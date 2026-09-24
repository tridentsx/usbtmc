// Copyright (c) 2015-2026 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

// Package wire implements the USBTMC and USBTMC-USB488 wire format:
// message ID values, bRequest values, status codes, and Bulk-IN/Bulk-OUT
// header encoding. It knows nothing about libusb, gousb, or any other USB
// transport -- everything here is pure encoding/decoding of bytes already
// agreed to have come from (or be going to) a USB bulk or control transfer.
//
// This package exists so that USBTMC device-side code -- something that
// answers these messages, such as the firmware transports at
// https://github.com/tridentsx/tmc-gateway -- can reuse the exact same wire
// vocabulary the host-side Device type in the parent package encodes and
// decodes, instead of duplicating it. Before this package existed, every
// name here was private to the parent package and unusable from outside it.
package wire

// ReservedField is USBTMC's reserved header byte, which every table in the
// spec requires to be 0x00.
const ReservedField = 0x00

// HeaderSize is the fixed size, in bytes, of every USBTMC Bulk-IN and
// Bulk-OUT header (USBTMC Table 1).
const HeaderSize = 12

// InterfaceClass is the first of the three bytes comprising the USB class
// code identifying a device's functionality (see
// http://www.usb.org/developers/defined_class). USBTMC's own bInterfaceClass
// value is ApplicationSpecificBaseClass.
type InterfaceClass byte

// ApplicationSpecificBaseClass is the only base class code USBTMC requires.
const ApplicationSpecificBaseClass InterfaceClass = 0xfe

// InterfaceSubClass is the second of the three bytes comprising the USB
// class code identifying a device's functionality.
type InterfaceSubClass byte

// USBTMCSubClass is USBTMC's own bInterfaceSubClass value.
const USBTMCSubClass InterfaceSubClass = 0x03

// InterfaceProtocol is the third of the three bytes comprising the USB class
// code identifying a device's functionality.
type InterfaceProtocol byte

// USBTMCProtocol and USB488Protocol are the two bInterfaceProtocol values
// USBTMC defines: plain USBTMC, and the USB488 subclass that adds GPIB-like
// control requests (READ_STATUS_BYTE, REN_CONTROL, GO_TO_LOCAL,
// LOCAL_LOCKOUT) and the TRIGGER message.
const (
	USBTMCProtocol InterfaceProtocol = 0x00
	USB488Protocol InterfaceProtocol = 0x01
)

// MessageID identifies the kind of USBTMC message a Bulk-IN or Bulk-OUT
// header describes.
type MessageID uint8

// The following MessageID values are found in Table 2 under the MACRO
// column of the USBTMC Specification 1.0, April 14, 2003. The end of line
// comment shows the MACRO names as given in the USBTMC specification.
// Trigger comes from Table 1 -- USB488 defined MessageID values of the
// USBTMC-USB488 Specification 1.0, April 14, 2003.
//
// RequestDevDepMsgIn and DevDepMsgIn share the value 2 deliberately: USBTMC
// disambiguates by direction, not by a distinct wire value. A host sends
// RequestDevDepMsgIn on Bulk-OUT to ask for data; a device replies with
// DevDepMsgIn on Bulk-IN. The same is true of the VendorSpecificIn pair.
const (
	DevDepMsgOut            MessageID = 1   // DEV_DEP_MSG_OUT
	RequestDevDepMsgIn      MessageID = 2   // REQUEST_DEV_DEP_MSG_IN
	DevDepMsgIn             MessageID = 2   // DEV_DEP_MSG_IN
	VendorSpecificOut       MessageID = 126 // VENDOR_SPECIFIC_OUT
	RequestVendorSpecificIn MessageID = 127 // REQUEST_VENDOR_SPECIFIC_IN
	VendorSpecificIn        MessageID = 127 // VENDOR_SPECIFIC_IN
	Trigger                 MessageID = 128 // TRIGGER
)

// Request is a USBTMC control-transfer bRequest value.
type Request uint8

// The USBTMC bRequest constants come from Table 15 -- USBTMC bRequest
// values in the USBTMC Specification 1.0, April 14, 2003.
const (
	InitiateAbortBulkOut    Request = 1  // INITIATE_ABORT_BULK_OUT
	CheckAbortBulkOutStatus Request = 2  // CHECK_ABORT_BULK_OUT_STATUS
	InitiateAbortBulkIn     Request = 3  // INITIATE_ABORT_BULK_IN
	CheckAbortBulkInStatus  Request = 4  // CHECK_ABORT_BULK_IN_STATUS
	InitiateClear           Request = 5  // INITIATE_CLEAR
	CheckClearStatus        Request = 6  // CHECK_CLEAR_STATUS
	GetCapabilities         Request = 7  // GET_CAPABILITIES
	IndicatorPulse          Request = 64 // INDICATOR_PULSE
)

// The USB488 bRequest constants come from Table 9 -- USB488 defined
// bRequest values in the USBTMC-USB488 Specification 1.0, April 14, 2003.
const (
	ReadStatusByte Request = 128 // READ_STATUS_BYTE
	RENControl     Request = 160 // REN_CONTROL
	GoToLocal      Request = 161 // GO_TO_LOCAL
	LocalLockout   Request = 162 // LOCAL_LOCKOUT
)

var requestDescription = map[Request]string{
	InitiateAbortBulkOut:    "Aborts a Bulk-OUT transfer.",
	CheckAbortBulkOutStatus: "Returns the status of the previously sent InitiateAbortBulkOut request.",
	InitiateAbortBulkIn:     "Aborts a Bulk-IN transfer.",
	CheckAbortBulkInStatus:  "Returns the status of the previously sent InitiateAbortBulkIn request.",
	InitiateClear:           "Clears all previously sent pending and unprocessed Bulk-OUT USBTMC message content and clears all pending Bulk-IN transfers from the USBTMC interface.",
	CheckClearStatus:        "Returns the status of the previously sent InitiateClear request.",
	GetCapabilities:         "Returns attributes and capabilities of the USBTMC interface.",
	IndicatorPulse:          "A mechanism to turn on an activity indicator for identification purposes. The device indicates whether or not it supports this request in the GetCapabilities response packet.",
	ReadStatusByte:          "Returns the IEEE 488 Status Byte.",
	RENControl:              "Mechanism to enable or disable local controls on a device.",
	GoToLocal:               "Mechanism to enable local controls on a device.",
	LocalLockout:            "Mechanism to disable local controls on a device.",
}

// String returns the request's description from the USBTMC/USB488 spec
// tables.
func (req Request) String() string {
	return requestDescription[req]
}

// Status is a USBTMC_status value, returned in a control-transfer response.
type Status byte

// The Status constant values come from Table 16 -- USBTMC_status values in
// the USBTMC Specification 1.0, April 14, 2003, and from Table 10 -- USB488
// defined USBTMC_status values in the USBTMC-USB488 Specification 1.0,
// April 14, 2003.
const (
	StatusSuccess               Status = 0x01 // STATUS_SUCCESS
	StatusPending               Status = 0x02 // STATUS_PENDING
	StatusInterruptInBusy       Status = 0x20 // STATUS_INTERRUPT_IN_BUSY
	StatusFailed                Status = 0x80 // STATUS_FAILED
	StatusTransferNotInProgress Status = 0x81 // STATUS_TRANSFER_NOT_IN_PROGRESS
	StatusSplitNotInProgress    Status = 0x82 // STATUS_SPLIT_NOT_IN_PROGRESS
	StatusSplitInProgress       Status = 0x83 // STATUS_SPLIT_IN_PROGRESS
)

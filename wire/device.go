// Copyright (c) 2015-2026 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package wire

import "encoding/binary"

// DecodeHeaderPrefix decodes the first four bytes of any USBTMC header
// (USBTMC Table 1): MsgID, bTag, and bTag's one's-complement check byte.
// ok is false if hdr is shorter than 4 bytes or the bTag/bTagInverse pair
// is inconsistent, which per the spec means the header must be discarded
// rather than acted on.
func DecodeHeaderPrefix(hdr []byte) (id MessageID, tag byte, ok bool) {
	if len(hdr) < 4 {
		return 0, 0, false
	}
	tag = hdr[1]
	if hdr[2] != InvertTag(tag) {
		return 0, 0, false
	}
	return MessageID(hdr[0]), tag, true
}

// DecodeBulkOutHeader decodes a DevDepMsgOut Bulk-OUT header (USBTMC
// Table 3), the message a host sends to write data to a device. It does
// not check that hdr's MsgID is actually DevDepMsgOut; callers dispatching
// on MessageID via DecodeHeaderPrefix first are expected to only reach
// here once they already know that.
//
// Verified against two independent real device-side implementations
// rather than inferred from EncodeBulkOutHeader's own OUT-direction code:
// OpenFFBoard's usbtmc_device.c and the Linux kernel's
// drivers/usb/class/usbtmc.c, which decode this exact same 12-byte shape
// from the opposite direction (Table 3's response counterpart, Table 2's
// DevDepMsgIn, is structurally identical -- see EncodeDevDepMsgInHeader).
func DecodeBulkOutHeader(hdr []byte) (tag byte, transferSize uint32, eom bool, ok bool) {
	if len(hdr) < HeaderSize {
		return 0, 0, false, false
	}
	id, tag, ok := DecodeHeaderPrefix(hdr)
	if !ok || id != DevDepMsgOut {
		return 0, 0, false, false
	}
	transferSize = binary.LittleEndian.Uint32(hdr[4:8])
	eom = hdr[8]&0x01 != 0
	return tag, transferSize, eom, true
}

// DecodeRequestDevDepMsgInHeader decodes a RequestDevDepMsgIn Bulk-OUT
// header (USBTMC Table 4), the message a host sends to ask a device for
// data. It does not check that hdr's MsgID is actually RequestDevDepMsgIn;
// see DecodeBulkOutHeader's doc comment for why.
func DecodeRequestDevDepMsgInHeader(hdr []byte) (tag byte, transferSize uint32, termCharEnabled bool, termChar byte, ok bool) {
	if len(hdr) < HeaderSize {
		return 0, 0, false, 0, false
	}
	id, tag, ok := DecodeHeaderPrefix(hdr)
	if !ok || id != RequestDevDepMsgIn {
		return 0, 0, false, 0, false
	}
	transferSize = binary.LittleEndian.Uint32(hdr[4:8])
	termCharEnabled = hdr[8]&0x02 != 0
	termChar = hdr[9]
	return tag, transferSize, termCharEnabled, termChar, true
}

// EncodeDevDepMsgInHeader creates the DevDepMsgIn Bulk-IN header a device
// sends back in response to a RequestDevDepMsgIn, with command-specific
// content as shown in USBTMC Table 2 (the response counterpart of
// EncodeBulkOutHeader's Table 3, and structurally identical to it: same
// four-byte prefix, the same TransferSize field, and bmTransferAttributes'
// D0 is EOM here exactly as it is in Table 3, just describing the
// device-to-host direction instead).
//
// usingTermChar reports whether TermChar (not carried in this header --
// only whether it was used) ended this transfer, USBTMC's bit D1 of
// bmTransferAttributes for this message. Set it only when the host
// enabled TermChar in the RequestDevDepMsgIn that prompted this response
// (DecodeRequestDevDepMsgInHeader's termCharEnabled) and this transfer
// really did stop because of it, not because eom was reached first.
//
// Verified against OpenFFBoard's usbtmc_device.c, a real device-side
// implementation building this exact header
// (usbtmc_msg_dev_dep_msg_in_header_t), and cross-checked against the
// Linux kernel's drivers/usb/class/usbtmc.c decoding the same shape from
// the host side.
func EncodeDevDepMsgInHeader(tag byte, transferSize uint32, eom bool, usingTermChar bool) [12]byte {
	prefix := EncodeHeaderPrefix(tag, DevDepMsgIn)
	packedTransferSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(packedTransferSize, transferSize)
	bmTransferAttributes := byte(0x00)
	if eom {
		bmTransferAttributes |= 0x01
	}
	if usingTermChar {
		bmTransferAttributes |= 0x02
	}
	return [12]byte{
		prefix[0],
		prefix[1],
		prefix[2],
		prefix[3],
		packedTransferSize[0],
		packedTransferSize[1],
		packedTransferSize[2],
		packedTransferSize[3],
		bmTransferAttributes,
		ReservedField,
		ReservedField,
		ReservedField,
	}
}

// Copyright (c) 2015-2026 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package wire

import "encoding/binary"

// NextTag returns the next bTag given the current one. Per the USBTMC
// standard, "the Host must set bTag such that 1<=bTag<=255."
func NextTag(tag byte) byte {
	return (tag % 255) + 1
}

// InvertTag returns the one's complement (inverse) of the given bTag.
func InvertTag(tag byte) byte {
	return tag ^ 0xff
}

// EncodeHeaderPrefix creates the first four bytes of a USBTMC Bulk-IN or
// Bulk-OUT header, as shown in USBTMC Table 1. id must match USBTMC Table 2.
func EncodeHeaderPrefix(tag byte, id MessageID) [4]byte {
	return [4]byte{
		byte(id),
		tag,
		InvertTag(tag),
		ReservedField,
	}
}

// EncodeBulkOutHeader creates the DevDepMsgOut Bulk-OUT header with
// command-specific content, as shown in USBTMC Table 3. This is what a host
// sends; a device implementation decodes the same 12 bytes with
// DecodeHeaderPrefix plus its own reading of offsets 4-11.
func EncodeBulkOutHeader(tag byte, transferSize uint32, eom bool) [12]byte {
	// Offset 0-3: see Table 1.
	prefix := EncodeHeaderPrefix(tag, DevDepMsgOut)
	// Offset 4-7: TransferSize. Per USBTMC Table 3, the "total number of
	// USBTMC message data bytes to be sent in this USB transfer. This does
	// not include the number of bytes in this Bulk-OUT Header or alignment
	// bytes. Sent least significant byte first, most significant byte
	// last. TransferSize must be > 0x00000000."
	packedTransferSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(packedTransferSize, transferSize)
	// Offset 8: bmTransferAttributes. Per USBTMC Table 3, D0 of
	// bmTransferAttributes:
	//   1 - The last USBTMC message data byte in the transfer is the last
	//       byte of the USBTMC message.
	//   0 - The last USBTMC message data byte in the transfer is not the
	//       last byte of the USBTMC message.
	// All other bits of bmTransferAttributes must be 0.
	bmTransferAttributes := byte(0x00)
	if eom {
		bmTransferAttributes = byte(0x01)
	}
	// Offset 9-11: ReservedField. Must be 0x000000.
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

// EncodeRequestDevDepMsgInHeader creates the RequestDevDepMsgIn Bulk-OUT
// header with command-specific content, as shown in USBTMC Table 4. This is
// what a host sends to ask a device for data; it is not the header shape a
// device uses to actually send that data back (that is the DevDepMsgIn
// Bulk-IN header, USBTMC Table 2's response side, which this package does
// not yet encode -- see the package doc comment).
func EncodeRequestDevDepMsgInHeader(
	tag byte,
	transferSize uint32,
	termCharEnabled bool,
	termChar byte,
) [12]byte {
	// Offset 0-3: see Table 1.
	prefix := EncodeHeaderPrefix(tag, RequestDevDepMsgIn)
	// Offset 4-7: TransferSize. Per USBTMC Table 4, the "maximum number of
	// USBTMC message data bytes to be sent in response to the command.
	// This does not include the number of bytes in this Bulk-IN Header or
	// alignment bytes. Sent least significant byte first, most significant
	// byte last. TransferSize must be > 0x00000000."
	packedTransferSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(packedTransferSize, transferSize)
	// Offset 8: bmTransferAttributes. Per USBTMC Table 4, D1 of
	// bmTransferAttributes:
	//   1 - "The Bulk-IN transfer must terminate on the specified TermChar.
	//       The Host may only set this bit if the USBTMC interface
	//       indicates it supports TermChar in the GetCapabilities response
	//       packet."
	//   0 - "The device must ignore TermChar."
	// All other bits of bmTransferAttributes must be 0.
	bmTransferAttributes := byte(0x00)
	if termCharEnabled {
		bmTransferAttributes = byte(0x02)
	}
	// Offset 9: TermChar. Offset 10-11: ReservedField. Must be 0x0000.
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
		termChar,
		ReservedField,
		ReservedField,
	}
}

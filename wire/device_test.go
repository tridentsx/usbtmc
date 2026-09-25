// Copyright (c) 2015-2026 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package wire

import "testing"

func TestDecodeHeaderPrefix(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		hdr := EncodeHeaderPrefix(5, DevDepMsgOut)
		id, tag, ok := DecodeHeaderPrefix(hdr[:])
		if !ok || id != DevDepMsgOut || tag != 5 {
			t.Errorf("DecodeHeaderPrefix() = (%v, %v, %v), want (DevDepMsgOut, 5, true)", id, tag, ok)
		}
	})

	t.Run("bTagInverse mismatch", func(t *testing.T) {
		hdr := EncodeHeaderPrefix(5, DevDepMsgOut)
		hdr[2] ^= 0xff // corrupt the check byte
		if _, _, ok := DecodeHeaderPrefix(hdr[:]); ok {
			t.Error("DecodeHeaderPrefix() = ok true for a corrupted bTagInverse, want false")
		}
	})

	t.Run("too short", func(t *testing.T) {
		if _, _, ok := DecodeHeaderPrefix([]byte{1, 2, 3}); ok {
			t.Error("DecodeHeaderPrefix() = ok true for a 3-byte buffer, want false")
		}
	})
}

func TestDecodeBulkOutHeaderRoundTrip(t *testing.T) {
	tests := []struct {
		name         string
		tag          byte
		transferSize uint32
		eom          bool
	}{
		{"eom_tag1", 1, 9, true},
		{"noEom_tag2", 2, 256, false},
		{"eom_tag255", 255, 512, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := EncodeBulkOutHeader(tt.tag, tt.transferSize, tt.eom)
			tag, transferSize, eom, ok := DecodeBulkOutHeader(encoded[:])
			if !ok {
				t.Fatal("DecodeBulkOutHeader() = ok false, want true")
			}
			if tag != tt.tag || transferSize != tt.transferSize || eom != tt.eom {
				t.Errorf("DecodeBulkOutHeader() = (%d, %d, %v), want (%d, %d, %v)",
					tag, transferSize, eom, tt.tag, tt.transferSize, tt.eom)
			}
		})
	}
}

// TestDecodeBulkOutHeaderRejectsWrongMessageID checks that a header
// carrying some other MessageID -- e.g. RequestDevDepMsgIn, sent to the
// wrong decoder -- is rejected rather than silently misinterpreted as a
// DevDepMsgOut.
func TestDecodeBulkOutHeaderRejectsWrongMessageID(t *testing.T) {
	encoded := EncodeRequestDevDepMsgInHeader(1, 64, true, '\n')
	if _, _, _, ok := DecodeBulkOutHeader(encoded[:]); ok {
		t.Error("DecodeBulkOutHeader() accepted a RequestDevDepMsgIn header, want rejected")
	}
}

func TestDecodeRequestDevDepMsgInHeaderRoundTrip(t *testing.T) {
	tests := []struct {
		name            string
		tag             byte
		transferSize    uint32
		termCharEnabled bool
		termChar        byte
	}{
		{"termChar_tag1", 1, 9, true, '\n'},
		{"noTermChar_tag2", 2, 512, false, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := EncodeRequestDevDepMsgInHeader(tt.tag, tt.transferSize, tt.termCharEnabled, tt.termChar)
			tag, transferSize, termCharEnabled, termChar, ok := DecodeRequestDevDepMsgInHeader(encoded[:])
			if !ok {
				t.Fatal("DecodeRequestDevDepMsgInHeader() = ok false, want true")
			}
			if tag != tt.tag || transferSize != tt.transferSize || termCharEnabled != tt.termCharEnabled {
				t.Errorf("DecodeRequestDevDepMsgInHeader() = (%d, %d, %v, %#x), want (%d, %d, %v, %#x)",
					tag, transferSize, termCharEnabled, termChar,
					tt.tag, tt.transferSize, tt.termCharEnabled, tt.termChar)
			}
			// termChar is only meaningful when termCharEnabled; still check it
			// round-trips exactly, since the byte is present either way.
			if termChar != tt.termChar {
				t.Errorf("termChar = %#x, want %#x", termChar, tt.termChar)
			}
		})
	}
}

func TestEncodeDevDepMsgInHeader(t *testing.T) {
	tests := []struct {
		name          string
		tag           byte
		transferSize  uint32
		eom           bool
		usingTermChar bool
		desired       [12]byte
	}{
		{
			"size9_eom_tag1",
			1, 9, true, false,
			[12]byte{0x02, 0x01, 0xfe, 0x00, 0x09, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00},
		},
		{
			"size256_noEom_noTermChar_tag2",
			2, 256, false, false,
			[12]byte{0x02, 0x02, 0xfd, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			"size9_eom_usingTermChar_tag1",
			1, 9, true, true,
			[12]byte{0x02, 0x01, 0xfe, 0x00, 0x09, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00},
		},
		{
			"size9_noEom_usingTermChar_tag1",
			1, 9, false, true,
			[12]byte{0x02, 0x01, 0xfe, 0x00, 0x09, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeDevDepMsgInHeader(tt.tag, tt.transferSize, tt.eom, tt.usingTermChar)
			if got != tt.desired {
				t.Errorf("EncodeDevDepMsgInHeader() = %x, want %x", got, tt.desired)
			}
		})
	}
}

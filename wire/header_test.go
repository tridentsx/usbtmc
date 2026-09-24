// Copyright (c) 2015-2026 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package wire

import (
	"fmt"
	"testing"
)

func TestNextTag(t *testing.T) {
	testCases := []struct {
		tag     byte
		nextTag byte
	}{
		{0x01, 0x02},
		{0xff, 0x01},
		{255, 1},
		{1, 2},
		{10, 11},
		{254, 255},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("tag_%d", tc.tag), func(t *testing.T) {
			got := NextTag(tc.tag)
			if got != tc.nextTag {
				t.Errorf(
					"NextTag == %x, want %x for given tag %x",
					got, tc.nextTag, tc.tag)
			}
		})
	}
}

func TestInvertTag(t *testing.T) {
	testCases := []struct {
		tag        byte
		tagInverse byte
	}{
		{0x00, 0xff},
		{0x0f, 0xf0},
		{0x55, 0xaa},
		{0xaa, 0x55},
		{0xf0, 0x0f},
		{0xff, 0x00},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("tag_%02x", tc.tag), func(t *testing.T) {
			got := InvertTag(tc.tag)
			if got != tc.tagInverse {
				t.Errorf(
					"tagInverse == %x, want %x for tag %x",
					got, tc.tagInverse, tc.tag)
			}
		})
	}
}

func TestEncodeHeaderPrefix(t *testing.T) {
	tests := []struct {
		name         string
		msgID        MessageID
		tag          byte
		headerPrefix [4]byte
	}{
		{"devDepMsgOut_tag2", DevDepMsgOut, 2, [4]byte{0x01, 0x02, 0xfd, 0x00}},
		{"devDepMsgOut_tag129", DevDepMsgOut, 129, [4]byte{0x01, 0x81, 0x7e, 0x00}},
		{"devDepMsgOut_tag255", DevDepMsgOut, 255, [4]byte{0x01, 0xff, 0x00, 0x00}},
		{"devDepMsgOut_tag1", DevDepMsgOut, 1, [4]byte{0x01, 0x01, 0xfe, 0x00}},
		{"requestDevDepMsgIn_tag4", RequestDevDepMsgIn, 4, [4]byte{0x02, 0x04, 0xfb, 0x00}},
		{"vendorSpecificOut_tag4", VendorSpecificOut, 4, [4]byte{0x7e, 0x04, 0xfb, 0x00}},
		{
			"requestVendorSpecificIn_tag4",
			RequestVendorSpecificIn, 4, [4]byte{0x7f, 0x04, 0xfb, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeHeaderPrefix(tt.tag, tt.msgID)
			if got != tt.headerPrefix {
				t.Errorf(
					"headerPrefix == %x, want %x",
					got, tt.headerPrefix)
			}
		})
	}
}

func TestEncodeBulkOutHeader(t *testing.T) {
	tests := []struct {
		name         string
		transferSize uint32
		eom          bool
		tag          byte
		desired      [12]byte
	}{
		{
			"size9_eom_tag1",
			9, true, 1,
			[12]byte{0x01, 0x01, 0xfe, 0x00, 0x09, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00},
		},
		{
			"size256_noEom_tag2",
			256, false, 2,
			[12]byte{0x01, 0x02, 0xfd, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			"size256_eom_tag2",
			256, true, 2,
			[12]byte{0x01, 0x02, 0xfd, 0x00, 0x00, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00},
		},
		{
			"size512_eom_tag2",
			512, true, 2,
			[12]byte{0x01, 0x02, 0xfd, 0x00, 0x00, 0x02, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeBulkOutHeader(tt.tag, tt.transferSize, tt.eom)
			if got != tt.desired {
				t.Errorf("BulkOutHeader == %x, want %x", got, tt.desired)
			}
		})
	}
}

func TestEncodeRequestDevDepMsgInHeader(t *testing.T) {
	tests := []struct {
		name            string
		tag             byte
		transferSize    uint32
		termCharEnabled bool
		termChar        byte
		desired         [12]byte
	}{
		{
			"size9_termChar_tag1",
			1, 9, true, '\n',
			[12]byte{0x02, 0x01, 0xfe, 0x00, 0x09, 0x00, 0x00, 0x00, 0x02, 0x0a, 0x00, 0x00},
		},
		{
			"size512_termChar_tag2",
			2, 512, true, '\n',
			[12]byte{0x02, 0x02, 0xfd, 0x00, 0x00, 0x02, 0x00, 0x00, 0x02, 0x0a, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeRequestDevDepMsgInHeader(
				tt.tag,
				tt.transferSize,
				tt.termCharEnabled,
				tt.termChar,
			)
			if got != tt.desired {
				t.Errorf("BulkOutHeader == %x, want %x", got, tt.desired)
			}
		})
	}
}

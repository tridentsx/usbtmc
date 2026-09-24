// Copyright (c) 2015-2026 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package wire

import (
	"fmt"
	"testing"
)

func TestRequestString(t *testing.T) {
	testCases := []struct {
		request     Request
		description string
	}{
		{InitiateAbortBulkOut, "Aborts a Bulk-OUT transfer."},
		{ReadStatusByte, "Returns the IEEE 488 Status Byte."},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Request_%d", tc.request), func(t *testing.T) {
			if tc.request.String() != tc.description {
				t.Errorf(
					"request.String() == %s, want %s",
					tc.request.String(),
					tc.description,
				)
			}
		})
	}
}

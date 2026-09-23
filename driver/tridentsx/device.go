// Copyright (c) 2015-2026 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package tridentsx

import (
	"context"
	"fmt"
	"time"

	usb "github.com/tridentsx/go-usb"
)

// defaultTimeout matches the other two drivers' default (2 seconds),
// used whenever a context has no deadline of its own.
const defaultTimeout = 2 * time.Second

// Device implements driver.USBDevice using github.com/tridentsx/go-usb.
type Device struct {
	handle  *usb.DeviceHandle
	bulkIn  uint8
	bulkOut uint8
	intrIn  uint8
	hasIntr bool
	timeout time.Duration
}

// Close closes the Device.
func (d *Device) Close() error {
	return d.handle.Close()
}

// String provides the Stringer interface method for Device.
func (d *Device) String() string {
	desc := d.handle.Descriptor()
	return fmt.Sprintf("VID: %#04x, PID: %#04x", desc.VendorID, desc.ProductID)
}

// Write writes to the USB device's bulk out endpoint.
func (d *Device) Write(p []byte) (n int, err error) {
	return d.handle.BulkTransfer(d.bulkOut, p, d.timeout)
}

// WriteString writes the given string to the Device and returns the number
// of bytes written along with an error code.
func (d *Device) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

// Read reads from the USB device's bulk in endpoint.
func (d *Device) Read(p []byte) (n int, err error) {
	return d.handle.BulkTransfer(d.bulkIn, p, d.timeout)
}

// ReadContext reads from the USB device's bulk in endpoint in a
// context-aware manner. If the context has a deadline, it is used as the
// transfer's timeout instead of the device's default.
func (d *Device) ReadContext(ctx context.Context, p []byte) (n int, err error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	return d.handle.BulkTransfer(d.bulkIn, p, d.contextTimeout(ctx))
}

// WriteContext writes to the USB device's bulk out endpoint in a
// context-aware manner. If the context has a deadline, it is used as the
// transfer's timeout instead of the device's default.
func (d *Device) WriteContext(ctx context.Context, p []byte) (n int, err error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	return d.handle.BulkTransfer(d.bulkOut, p, d.contextTimeout(ctx))
}

// contextTimeout derives a transfer timeout from the context's deadline,
// falling back to the device's default when there is none.
func (d *Device) contextTimeout(ctx context.Context) time.Duration {
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return time.Millisecond // minimum timeout to avoid blocking indefinitely
		}
		return remaining
	}
	return d.timeout
}

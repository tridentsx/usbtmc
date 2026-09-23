// Copyright (c) 2015-2026 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

// Package tridentsx implements the usbtmc driver interfaces using
// github.com/tridentsx/go-usb: a pure-Go USB backend with no cgo and no
// libusb C library dependency, unlike driver/google and driver/gotmc. Any
// program that only needs this driver builds and cross-compiles with
// nothing but the Go toolchain -- no libusb installed, no C compiler, no
// CGO_ENABLED=1.
package tridentsx

import (
	"fmt"

	usb "github.com/tridentsx/go-usb"

	"github.com/gotmc/usbtmc"
	"github.com/gotmc/usbtmc/driver"
)

// Driver implements the driver.Driver interface required by usbtmc using
// github.com/tridentsx/go-usb.
type Driver struct{}

func init() {
	usbtmc.Register(&Driver{})
}

// Context satisfies driver.Context. Unlike libusb, go-usb has no
// process-wide session object to create or tear down: DeviceList and
// OpenDevice call directly into the OS on each use, so this exists only to
// satisfy the interface.
type Context struct{}

// NewContext returns a new Context. It never fails, since there is nothing
// to actually initialize.
func (d Driver) NewContext() (driver.Context, error) {
	return &Context{}, nil
}

// SetDebugLevel is a no-op: go-usb has no equivalent verbose-logging knob
// to configure here.
func (c *Context) SetDebugLevel(level int) {}

// Close is a no-op: there is no session-level resource to release.
func (c *Context) Close() error {
	return nil
}

// NewDeviceByVIDPID opens the first USB device matching the given vendor ID
// and product ID, claims the first interface that has both a bulk IN and a
// bulk OUT endpoint (the USBTMC requirement per Section 3.2 of the USBTMC
// spec), and returns it ready for USBTMC framing.
func (c *Context) NewDeviceByVIDPID(VID, PID int) (driver.USBDevice, error) {
	handle, err := usb.OpenDevice(uint16(VID), uint16(PID)) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("opening device %04x:%04x: %w", VID, PID, err)
	}

	cfg, err := handle.GetActiveConfigDescriptor()
	if err != nil {
		_ = handle.Close()
		return nil, fmt.Errorf("getting active config descriptor: %w", err)
	}

	eps, found := findUSBTMCEndpoints(cfg)
	if !found {
		_ = handle.Close()
		return nil, fmt.Errorf("no interface with both a bulk IN and bulk OUT endpoint found on %04x:%04x", VID, PID)
	}

	if err := handle.ClaimInterface(eps.ifaceNum); err != nil {
		_ = handle.Close()
		return nil, fmt.Errorf("claiming interface %d: %w", eps.ifaceNum, err)
	}

	return &Device{
		handle:  handle,
		bulkIn:  eps.bulkIn,
		bulkOut: eps.bulkOut,
		intrIn:  eps.intrIn,
		hasIntr: eps.hasIntr,
		timeout: defaultTimeout,
	}, nil
}

// usbtmcEndpoints holds the endpoint addresses an interface needs to carry
// USBTMC traffic.
type usbtmcEndpoints struct {
	bulkIn, bulkOut, intrIn uint8
	hasIntr                 bool
	ifaceNum                uint8
}

// findUSBTMCEndpoints scans every interface's first alternate setting for
// one with both a bulk IN and a bulk OUT endpoint, returning that
// interface's endpoint addresses (plus an interrupt IN endpoint, if
// present, which USBTMC uses for status notification per Section 3.3).
func findUSBTMCEndpoints(cfg *usb.ConfigDescriptor) (usbtmcEndpoints, bool) {
	for _, iface := range cfg.Interfaces {
		if len(iface.AltSettings) == 0 {
			continue
		}
		alt := iface.AltSettings[0]

		var eps usbtmcEndpoints
		var haveIn, haveOut bool
		for i := range alt.Endpoints {
			ep := &alt.Endpoints[i]
			switch ep.TransferType() {
			case usb.TransferTypeBulk:
				if ep.IsInput() {
					eps.bulkIn, haveIn = ep.EndpointAddr, true
				} else {
					eps.bulkOut, haveOut = ep.EndpointAddr, true
				}
			case usb.TransferTypeInterrupt:
				if ep.IsInput() {
					eps.intrIn, eps.hasIntr = ep.EndpointAddr, true
				}
			}
		}
		if haveIn && haveOut {
			eps.ifaceNum = alt.InterfaceNumber
			return eps, true
		}
	}
	return usbtmcEndpoints{}, false
}

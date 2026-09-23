# usbtmc

Go library to communicate with a USB Test and Measurement Class (USBTMC)
enabled USB device.

[![GoDoc][godoc badge]][godoc link]
[![Go Report Card][report badge]][report card]
[![Build Status][travis badge]][travis link]
[![License Badge][license badge]][LICENSE.txt]

## Overview

[USBTMC][] is a USB device class specification for test equiment and
instrumentation devices, such as oscilloscopes, digital multimeters, and
function generators. USBTMC requires three endpoints:

- Control endpoint
- Bulk-OUT endpoint
- Bulk-IN endpoint

Additionally, the USBTMC subclass USB488 has an Interrupt-IN endpoint.

This packages enables controlling USBTMC compatible test equipment (e.g.,
oscilloscopes, function generators, multimeters, etc.) over USB. While this
package can be used by itself to send Standard Commands for Programmable
Instruments ([SCPI][]) commands to a piece of test equipment, it also serves to
provide an Instrument interface for both the [ivi][] and [visa][] packages. The
[ivi][] package provides standardized APIs for programming test instruments
following the [Interchangeable Virtual Instrument (IVI) standard][ivi-specs].

## USBTMC Descriptors

- Interface class = 0xFE (application-specific)
- Interface subclass = 0x03 (indicates USBTMC)

## Installation

```bash
$ go get github.com/gotmc/usbtmc
```

## Usage

To use the [usbtmc][gousbtmc] package, you must register which USB driver
backend should be used. Three backends are available:

- [github.com/google/gousb][gousb] — Go bindings for the [libusb][] C
  library
- [github.com/gotmc/libusb][golibusb] — another set of Go bindings for the
  [libusb][] C library
- [github.com/tridentsx/go-usb][gousb2] — a pure Go USB backend with no
  cgo and no [libusb][] C library dependency

The first two require [libusb][] ("a C library that provides generic
access to USB devices") installed on the system and a C compiler (cgo) to
build. The third needs neither: `driver/tridentsx` builds and
cross-compiles with nothing but the Go toolchain.

You'll need to install **_one_** of the above libraries using:

```bash
$ go get -v github.com/google/gousb
$ go get -v github.com/gotmc/gotmc
$ go get -v github.com/tridentsx/go-usb
```

To indicate which driver backend should be used, include **_one_** of the
following blank imports:

```go
import _ "github.com/gotmc/usbtmc/driver/google"
import _ "github.com/gotmc/usbtmc/driver/gotmc"
import _ "github.com/gotmc/usbtmc/driver/tridentsx"
```

## Documentation

Documentation can be found at either:

- <https://godoc.org/github.com/gotmc/usbtmc>
- <http://localhost:6060/pkg/github.com/gotmc/usbtmc/> after running `$
godoc -http=:6060`

## Contributing

Contributions are welcome! To contribute please:

1. Fork the repository
2. Create a feature branch
3. Code
4. Submit a [pull request][]

### Development Dependencies

- [just][] - task runner that replaces [GNU Make][make]

### Testing

Prior to submitting a [pull request][], please run:

```bash
$ just check
$ just lint
```

To update and view the test coverage report:

```bash
$ just cover
```

### Disclosure and Call for Help

While this package works, it does not fully implement the [USBTMC][]
specification. Please submit pull requests as needed to increase functionality,
maintainability, or reliability.

## License

[usbtmc][gousbtmc] is released under the MIT license. Please see the
[LICENSE.txt][] file for more information.

[godoc badge]: https://godoc.org/github.com/gotmc/usbtmc?status.svg
[godoc link]: https://godoc.org/github.com/gotmc/usbtmc
[golibusb]: https://github.com/gotmc/libusb
[gousb]: https://github.com/google/gousb
[gousb2]: https://github.com/tridentsx/go-usb
[gousbtmc]: https://github.com/gotmc/usbtmc
[ivi]: https://github.com/gotmc/ivi
[ivi-foundation]: http://www.ivifoundation.org/
[ivi-specs]: http://www.ivifoundation.org/specifications/
[just]: https://just.systems/man/en/
[libusb]: http://libusb.info
[LICENSE.txt]: https://github.com/gotmc/libusb/blob/master/LICENSE.txt
[license badge]: https://img.shields.io/badge/license-MIT-blue.svg
[lxi]: https://github.com/gotmc/lxi
[make]: https://www.gnu.org/software/make/
[pull request]: https://help.github.com/articles/using-pull-requests
[report badge]: https://goreportcard.com/badge/github.com/gotmc/usbtmc
[report card]: https://goreportcard.com/report/github.com/gotmc/usbtmc
[scpi]: https://www.ivifoundation.org/About-IVI/scpi.html
[travis badge]: http://img.shields.io/travis/gotmc/usbtmc/master.svg
[travis link]: https://travis-ci.org/gotmc/usbtmc
[usbtmc]: http://www.usb.org/developers/docs/devclass_docs/
[visa]: https://github.com/gotmc/visa

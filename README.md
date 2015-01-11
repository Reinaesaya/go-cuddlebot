# Cuddlebot Control Server

The Cuddlebot Control Server `cuddled` is implemented using the
[Go Programming Language][go] as a [RESTful][restful] API.


## Getting Started

To get started, install Go from the project [website][go]. Then, install
the Go package dependencies:

```sh
go get github.com/codegangsta/negroni
go get github.com/phyber/negroni-gzip/gzip
go get github.com/stretchr/graceful
go get github.com/mikepb/go-crc16
```

These packages include the [Negroni][negroni] HTTP Middleware for Go and
supporting packages.

To build binaries for Linux on ARM, you'll also need to install the
[GNU Tools for ARM Embedded Processors][gccarm]. Make sure that the tools
are available on your `PATH`.

A `Makefile` is available with the following targets:

- `build` compile `cuddled` and `cuddlespeak` for the current platform and
  for Linux/ARM
- `clean` remove the build directories

The binaries under `bin-arm-linux/` are used as part of the Yocto Embedded
Linux build process. More details are available as part of the Cuddlebot
system image project.


## Project File Organization

- `bin/` compiled binaries for the current platform
- `bin-arm-linux/` compiled binaries for the Linux/ARM
- `cuddle` implements the control server library
- `cuddled` implements the control server daemon
- `cuddlespeak` implements a command-line tool to control the motors


## License

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.


[go]: http://golang.org
[gccarm]: https://launchpad.net/gcc-arm-embedded
[restful]: http://www.restapitutorial.com
[negroni]: https://github.com/codegangsta/negroni
[yocto]: http://www.yoctoproject.org

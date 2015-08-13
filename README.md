# Cuddlebot Control Server

The Cuddlebot Control Server `cuddled` is implemented using the
[Go Programming Language][go] as a [RESTful][restful] API.

This version of the Control Server is an update from [mikepb's version][original]. It includes server control implementation for the smooth command, as well as detailed and improved build instructions for Ubuntu systems (up to date as of August 2015 for Ubuntu 14.04). This build is used in conjunction with instructions of the [Cuddlebot Yocto system image][cuddleyocto].


## Getting Started

### Go Setup

Open up terminal and check if Go exists
```bash
go version
```
If version is below 1.4, or gives an error, install Go with the following steps. If not, it may be possible to proceed directly to **Package Dependencies**.

If Go version exists, remove it:
```bash
# Remove base directories:
`sudo rm -rf /usr/lib/go`
# or
`sudo rm -rf /usr/local/go`
# or wherever go had been stored

# Remove existing golang directories:
`sudo apt-get remove golang-go`.
```

Install correct Go version:
```bash
# Install by source
cd /usr/lib/
git clone https://go.googlesource.com/go
cd go
git checkout go1.4.1
cd src
sudo ./all.bash

# Setup Go environment
mkdir $HOME/go
export PATH=$PATH:/usr/lib/go/bin
export GOPATH=$HOME/go
export GOROOT=/usr/lib/go
# Check
go env
go version

# Setup build environment
cd /usr/lib/go/src/
GOOS=linux GOARCH=arm GOARM=7 ./make.bash --no-clean
```

### Package Dependencies

Install the Go package dependencies:

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

The binaries under `bin-arm-linux/` are used as part of the Yocto Embedded Linux build process. More details are available as part of the [Cuddlebot system image project][cuddleyocto].


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

[original]: https://github.com/mikepb/go-cuddlebot
[cuddleyocto]: https://github.com/Reinaesaya/cuddlebot-yocto
[go]: http://golang.org
[gccarm]: https://launchpad.net/gcc-arm-embedded
[restful]: http://www.restapitutorial.com
[negroni]: https://github.com/codegangsta/negroni
[yocto]: http://www.yoctoproject.org

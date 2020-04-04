GOVM
----

[![GitHub release](https://img.shields.io/github/release/ljesparis/govm.svg)](https://github.com/ljesparis/govm/releases/latest)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/ljesparis/govm/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ljesparis/govm)](https://goreportcard.com/report/github.com/ljesparis/govm)

Download, configure and start using **go** never was so easy and faster.

**Note:** By now, govm is tested on linux.

<p align="center"><img src="./img/govm.gif?raw=true"/></p>

---
* [Installation](#installation)
* [Usage](#usage)
* [Get Involved](#get-involved)

---

## Installation

You can install pre-compiled binaries or compile from source.

#### Pre-compiled binaries
Download the latest released version from [releases page](https://github.com/ljesparis/govm/releases) and copy to the desired location.

#### Source (by now only supported on unix)
Follow the steps:
```sh
$ # Clone the repository outside the GOPATH
$ git clone https://github.com/ljesparis/govm.git
$ cd govm
$
$ # Build the source
$ make
$
$ # Install source whatever you want, by default will be installed at /usr/local/bin
$ make install 
```

#### Post-installation

##### Linux
On linux, govm will install go sources into **~/.local/bin** directory on unix.
Include this path into PATH environment variable if it does not exists already.

##### Windows 

On windows, govm will install go sources into **C:\Users\%USERPROFILE%\.govm\bin** directory,
make sure to add this folder into environment variables.

## Usage

```
govm is a go version manager

Usage:
  govm [flags] [command]
  govm [command]

Available Commands:
  delete      Delete golang source
  help        Help about any command
  list        List golang sources
  select      Select golang source

Flags:
  -h, --help      help for govm
      --version   version for govm

Use "govm [command] --help" for more information about a command.

```

## Get Involved

Pull requests are welcome!. Submit github issues for any feature enhancements, bugs or documentation problems.

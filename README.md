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

---

## Installation

You can install pre-compiled binaries or compile from source.

#### Pre-compiled binaries
Download the latest released version from [releases page](https://github.com/ljesparis/govm/releases) and copy to the desired location.

#### Source
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
By now **govm** does install go binaries at /usr/local/bin (this will change in the future). So if you want to install any source you want without any problem, follow the steps:
```sh
 $ # Lets change folder ownership
 $ sudo chown leonardoem:leonardoem /usr/local/bin
```

## Usage

```
govm is a go version manager

Usage:
  govm [flags] [command]
  govm [command]

Available Commands:
  help        Help about any command
  list        List golang sources
  select      Select golang source

Flags:
  -h, --help      help for govm
      --version   version for govm

Use "govm [command] --help" for more information about a command.

```

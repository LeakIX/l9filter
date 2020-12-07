# IP4scout

[![GitHub Release](https://img.shields.io/github/v/release/LeakIX/l9filter)](https://github.com/LeakIX/l9filter/releases)
[![Follow on Twitter](https://img.shields.io/twitter/follow/leak_ix.svg?logo=twitter)](https://twitter.com/leak_ix)

l9filter is a translation tool for stdin/stdout that understand the [l9format](https://github.com/LeakIX/l9format). Its main goal is to facilitate data communication
between different network recon software.  

## Features

- stdin/stdout oriented
- Conversion back and forward between formats
- Low memory/CPU footprint
- Easy to implement interface

## Usage

```
l9filter transform -h
```

Displays help for the random command (only implementation atm)

|Flag           |Description  |
|-----------------------|-------------------------------------------------------|
|--input-format       |Selects the input format to parse
|--output-format      |Selects the output format to use
|--source-file        |Use an input file instead of stdin|
|--target-file        |Use an output file instead of stdout|

## Installation Instructions

### From Binary

The installation is easy. You can download the pre-built binaries for your platform from the [Releases](https://github.com/LeakIX/l9filter/releases/) page.

```sh
▶ chmod +x l9filter-linux-64
▶ mv l9filter-linux-64 /usr/local/bin/l9filter
```

### From Source

```sh
▶ GO111MODULE=on go get -u -v github.com/LeakIX/l9filter/cmd/l9filter
▶ ${GOPATH}/bin/l9filter transform -h
```

## Running l9filter

l9filter requires an input to run. By default `stdin` will be used unles `input-file` is specified.

```sh
▶ l9filter -i l9 -o hostport
```

## Examples

[ip4scout](https://github.com/LeakIX/ip4scout) speaks [l9format](https://github.com/LeakIX/l9format) which is a JSON schema targeted at network recon.

Parsing its output would look like :

### Human output

```sh 
▶ ip4scout random --ports=3306,9200,6379|tee results.json|l9filter -i l9 -o human
```

Displays human-readable results on `stdout` while saving the scan results to `results.json` 


### Hostport output

```sh 
▶ ip4scout random -p 3304-3308,9200-9210|tee results.json|l9filter -i l9 -o hostport
```

Will display `host:port` results on `stdout` while saving the scan results to `results.json` 
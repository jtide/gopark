# GoPark
A demo REST API built with Golang for querying parking rates based on time-of-day.

# Installation

Use these instructions to install the `gopark` package into the go [workspace](https://golang.org/doc/code.html#Organization).

Prerequisites:
 - If not already installed, download Go from [here](https://golang.org/dl/)
 - Use terminal with GOPATH environment variable set

Install by performing:

    % go get github.com/jtide/gopark

The `gopark` package should now be installed at `$GOPATH/src/github.com/jtide/gopark` with the executable in
`$GOPATH/bin/gopark`.

# Usage

If $GOPATH/bin is in your path, then start the api with a default configuration by:

    % gopark

If desired, a custom configuration file can be used to specify the rates:

    % cd $GOPATH/src/github.com/jtide/gopark
    % ./gopark --config examples/sample-rates.json

# API Testing

For details on API testing and validation, see [API testing with curl](./doc/api-testing.md)

# Unit Tests
To run unit tests, install the following package, if not already installed.

    % go get github.com/stretchr/testify/assert

Then run unit tests by:

    % cd $GOPATH/src/github.com/jtide/gopark/api
    % go test -v

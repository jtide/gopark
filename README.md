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

    % cd $GOPATH/github.com/jtide/gopark
    % ./gopark --config examples/sample-rates.json

## Unit Tests
To run unit tests, install the following package, if not already installed.

    % go get github.com/stretchr/testify/assert

Then run unit tests by:

    % cd $GOPATH/github.com/jtide/gopark/api
    % go test -v

## API Testing with Curl

When run without configuration options, a set of default rates will be loaded from memory.

    % # Start the process without arguments to use default rates
    % ./gopark

    % curl -H "Accept: application/json"  "http://localhost:8080/api/rate?start=2015-07-01T07%3A00%3A00Z&end=2015-07-01T16%3A00%3A00Z"; echo
    {"start":"2015-07-01T07:00:00Z","end":"2015-07-01T16:00:00Z","price":1750}

    % curl -H "Accept: application/json"  "http://localhost:8080/api/rate?start=2015-07-01T07%3A00%3A00Z&end=2015-07-08T16%3A00%3A00Z"; echo
    {"start":"2015-07-01T07:00:00Z","end":"2015-07-08T16:00:00Z","price":"unavailable"}

Testing with `sample-rates.json` configuration file:

    % # Load the sample config file from app description
    % ./gopark --config examples/sample-rates.json

    % # Perform queries accepting JSON
    % curl -H "Accept: application/json"  "http://localhost:8080/api/rate?start=2015-07-01T07:00:00Z&end=2015-07-01T12:00:00Z"; echo
    {"start":"2015-07-01T07:00:00Z","end":"2015-07-01T12:00:00Z","price":1500}

    % curl -H "Accept: application/json"  "http://localhost:8080/api/rate?start=2015-07-04T07:00:00Z&end=2015-07-04T12:00:00Z"; echo
    {"start":"2015-07-04T07:00:00Z","end":"2015-07-04T12:00:00Z","price":2000}

    % curl -H "Accept: application/json"  "http://localhost:8080/api/rate?start=2015-07-04T07:00:00Z&end=2015-07-04T20:00:00Z"; echo
    {"start":"2015-07-04T07:00:00Z","end":"2015-07-04T20:00:00Z","price":"unavailable"}

    % # Perform queries accepting XML
    % curl -H "Accept: application/xml"  "http://localhost:8080/api/rate?start=2015-07-01T07:00:00Z&end=2015-07-01T12:00:00Z"; echo
    <Rate><Start>2015-07-01T07:00:00Z</Start><End>2015-07-01T12:00:00Z</End><Price>1500</Price></Rate>

    % curl -H "Accept: application/xml"  "http://localhost:8080/api/rate?start=2015-07-04T07:00:00Z&end=2015-07-04T12:00:00Z"; echo
    <Rate><Start>2015-07-04T07:00:00Z</Start><End>2015-07-04T12:00:00Z</End><Price>2000</Price></Rate>

    % curl -H "Accept: application/xml"  "http://localhost:8080/api/rate?start=2015-07-04T07:00:00Z&end=2015-07-04T20:00:00Z"; echo
    <UnknownRate><start>2015-07-04T07:00:00Z</start><end>2015-07-04T20:00:00Z</end><price>unavailable</price></UnknownRate>

In addition, a descriptive error message will be returned if the URL start or end parameters are not valid.

    % curl -H "Accept: application/json"  "http://localhost:8080/api/rate?start=2015-07-01T00:00:1234Z&end=2015-07-08T16:00:00Z"; echo
    {"error":"could not parse 'start' parameter [2015-07-01T00:00:1234Z]: parsing time \"2015-07-01T00:00:1234Z\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"34Z\" as \"Z07:00\""}

    % # View headers for above request with -I parameter to see 400 status code returned
    % curl -I -H "Accept: application/json"  "http://localhost:8080/api/rate?start=2015-07-01T00:00:1234Z&end=2015-07-08T16:00:00Z"; echo
    HTTP/1.1 400 Bad Request
    Content-Type: application/json; charset-utf-8
    Date: Wed, 02 May 2018 05:51:06 GMT
    Content-Length: 180

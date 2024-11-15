JSONx
========

[![GoDoc](https://godoc.org/github.com/jefferai/jsonx?status.svg)](https://godoc.org/github.com/jefferai/jsonx)

A Go (Golang) library to transform an object or existing JSON bytes into
[JSONx](https://www.ibm.com/support/knowledgecenter/SS9H2Y_7.5.0/com.ibm.dp.doc/json_jsonxconversionrules.html).
Because sometimes your luck runs out.

This follows the "standard" except for the handling of special and escaped
characters. Names and values are properly XML-escaped but there is no special
handling of values already escaped in JSON if they are valid in XML.

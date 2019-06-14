# cboring [![Build Status](https://travis-ci.org/dtn7/cboring.svg?branch=master)](https://travis-ci.org/dtn7/cboring) [![GoDoc](https://godoc.org/github.com/dtn7/cboring?status.svg)](https://godoc.org/github.com/dtn7/cboring)

Simple [CBOR][cbor] Go(lang) library for a selected subset of features,
developed to be used in [`dtn7-go`][dtn7-go], an implementation of the
[Bundle Protocol Version 7][bpbis]. The name is based on the fact that
`cboring` is both boring to use and bored about the amount of data to handle.


## Non-Features

- Supports the subset of [CBOR's][cbor] features necessary for [BPv7][bpbis]:
    - Unsigned Integer
    - Byte and Text String
    - Arrays, both of definite and indefinite length
- Small and clear codebase:
    - Only works on streams, Go's `io.Reader` or `io.Writer`
    - Does *not* use reflection or makes any strange assumptions
- Surprisingly fast


[bpbis]: https://tools.ietf.org/html/draft-ietf-dtn-bpbis-13
[cbor]: https://tools.ietf.org/html/rfc7049
[dtn7-go]: https://github.com/dtn7/dtn7-go

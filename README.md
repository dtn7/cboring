# cboring [![CI](https://github.com/dtn7/cboring/workflows/CI/badge.svg)](https://github.com/dtn7/cboring/actions) [![GoDoc](https://godoc.org/github.com/dtn7/cboring?status.svg)](https://godoc.org/github.com/dtn7/cboring)

A simple [CBOR][cbor] Go(lang) library for a selected subset of features,
developed to be used in [`dtn7-go`][dtn7-go], an implementation of the
[Bundle Protocol Version 7][bpbis]. The name is based on the fact that
`cboring` is both boring to use and bored about the amount of data to handle.


## Features

- Supports a selected subset of [CBOR's][cbor] features:
    - Unsigned Integer
    - Floating-point values
    - Byte and Text String
    - Arrays, both of definite and indefinite length
    - Maps of definite length
    - Booleans
- Small and clear codebase:
    - Only works on streams, Go's `io.Reader` or `io.Writer`
    - Does *not* use reflection or make any strange assumptions
- Surprisingly fast


[bpbis]: https://tools.ietf.org/html/draft-ietf-dtn-bpbis-29
[cbor]: https://tools.ietf.org/html/rfc7049
[dtn7-go]: https://github.com/dtn7/dtn7-go

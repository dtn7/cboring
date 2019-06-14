package main

import (
	"bytes"
	"fmt"

	_ "github.com/dtn7/cboring"
)

func main() {
	e := NewEndpointID("ipn:23.42")
	fmt.Println(e)

	buff := new(bytes.Buffer)
	e.MarshalCbor(buff)

	fmt.Printf("%x\n", buff.Bytes())
}

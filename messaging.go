/*
The MIT License (MIT)

Copyright © 2017 Peter Fern <golang@0xc0dedbad.com>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the "Software"),
to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/


package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"unsafe"
)

var nativeEndian binary.ByteOrder

func init() {
	var one int16 = 1
	b := (*byte)(unsafe.Pointer(&one))
	if *b == 0 {
		nativeEndian = binary.BigEndian
	} else {
		nativeEndian = binary.LittleEndian
	}
}

type encoder struct {
	w io.Writer
}

func newEncoder(w io.Writer) *encoder {
	return &encoder{w: w}
}

func (e *encoder) Encode(v interface{}) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	msgLen := uint32(len(buf))
	if err := binary.Write(e.w, nativeEndian, &msgLen); err != nil {
		return err
	}
	if _, err := e.w.Write(buf); err != nil {
		return err
	}
	return nil
}

type decoder struct {
	r io.Reader
}

func newDecoder(r io.Reader) *decoder {
	return &decoder{r: r}
}

func (d *decoder) Decode(v interface{}) error {
	var msgLen uint32
	if err := binary.Read(d.r, nativeEndian, &msgLen); err != nil {
		return err
	}
	buf := make([]byte, msgLen)
	if _, err := io.ReadFull(d.r, buf); err != nil {
		return err
	}
	return json.Unmarshal(buf, v)
}



/*
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <https://unlicense.org>
*/


package basefm

import (
	"io"
	"github.com/icza/bitio"
	"github.com/maxymania/weirdwire/hufftab"
)



type Encoder struct {
	wr *bitio.Writer
}
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{bitio.NewWriter(w)}
}

/*
Writes 7-bit clean data (or ASCII).
*/
func (e *Encoder) WriteAscii(str string) error {
	typs.TryWriteSymbol(e.wr,t_txt)
	for _,b := range []byte(str) {
		if b>=128 { continue }
		txt.TryWriteSymbol(e.wr,int(b))
	}
	txt.TryWriteSymbol(e.wr,128)
	return e.wr.TryError
}

/*
Writes 8-bit clean data (or UTF-8).
*/
func (e *Encoder) WriteUtf8(str string) error {
	typs.TryWriteSymbol(e.wr,t_utf)
	for _,b := range []byte(str) {
		utf.TryWriteSymbol(e.wr,int(b))
	}
	utf.TryWriteSymbol(e.wr,256)
	return e.wr.TryError
}

/*
Writes 8-bit clean, but NULL-terminated, data (or UTF-8).
*/
func (e *Encoder) WriteRaw(str string) error {
	typs.TryWriteSymbol(e.wr,t_bin)
	for _,b := range []byte(str) {
		if b==0 { continue }
		e.wr.TryWriteByte(b)
	}
	e.wr.TryWriteByte(0)
	utf.TryWriteSymbol(e.wr,256)
	return e.wr.TryError
}

/*
Writes text-data adaptively using eighter WriteAscii(),
WriteRaw() or WriteUtf8(), depending on the data.

WriteAscii() is prefered over WriteRaw(), and WriteRaw() is
prefered over WriteUtf8().
*/
func (e *Encoder) WriteEncoded(str string) error {
	has0,has8 := false,false
	for _,b := range []byte(str) {
		has0 = has0 || b==0
		has8 = has8 || b>=128
	}
	if !has8 { return e.WriteAscii(str) } // If data is not 8-bit-clean.
	if !has0 { return e.WriteRaw(str) } // If data contains a NULL-byte.
	return e.WriteUtf8(str)
}
func (e *Encoder) WriteSymbol(s uint32) error {
	typs.TryWriteSymbol(e.wr,t_symbol)
	if s<(1<<4) { // 0 xxxx
		e.wr.TryWriteBool(false)
		e.wr.TryWriteBits(uint64(s),4)
	} else if s<(1<<8) { // 10 xxxx xxxx
		e.wr.TryWriteBool(true)
		e.wr.TryWriteBool(false)
		e.wr.TryWriteBits(uint64(s),8)
	} else if s<(1<<16) { // 110 xxxx xxxx xxxx xxxx
		e.wr.TryWriteBool(true)
		e.wr.TryWriteBool(true)
		e.wr.TryWriteBool(false)
		e.wr.TryWriteBits(uint64(s),16)
	} else { // 111 xxxxxxxx*4
		e.wr.TryWriteBool(true)
		e.wr.TryWriteBool(true)
		e.wr.TryWriteBool(true)
		e.wr.TryWriteBits(uint64(s),32)
	}
	return e.wr.TryError
}
func (e *Encoder) Write(p []byte) (n int, err error) { return e.wr.Write(p) }
func (e *Encoder) WriteByte(b byte) (err error) { return e.wr.WriteByte(b) }
func (e *Encoder) Align() (skipped byte, err error) { return e.wr.Align() }

type Decoder struct {
	rd *bitio.Reader
}
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{bitio.NewReader(r)}
}

func (d *Decoder) readString(t *hufftab.Table, eos int) (string,error) {
	buf := make([]byte,0,16)
	for {
		s := t.TryReadSymbol(d.rd)
		if d.rd.TryError!=nil { return "",d.rd.TryError }
		if s>=eos { break }
		buf = append(buf,byte(s))
	}
	return string(buf),nil
}
func (d *Decoder) readRaw() (string,error) {
	buf := make([]byte,0,16)
	for {
		s := d.rd.TryReadByte()
		if d.rd.TryError!=nil { return "",d.rd.TryError }
		if s==0 { break }
		buf = append(buf,byte(s))
	}
	return string(buf),nil
}
func (d *Decoder) readSymbol() (r uint64,e error) {
	if !d.rd.TryReadBool() { // 0 xxxx
		r = d.rd.TryReadBits(4)
	} else if !d.rd.TryReadBool() { // 10 xxxx xxxx
		r = d.rd.TryReadBits(8)
	} else if !d.rd.TryReadBool() { // 110 xxxx xxxx xxxx xxxx
		r = d.rd.TryReadBits(16)
	} else { // 111 xxxxxxxx*4
		r = d.rd.TryReadBits(32)
	}
	e = d.rd.TryError
	return
}


func conv1(s string,e error) (uint32,string,error) { return 0,s,e }
func conv2(e error) (uint32,string,error) { return 0,"",e }
func conv3(u uint64,e error) (uint32,string,error) { return uint32(u),"",e }

func (d *Decoder) ReadSymbol() (uint32,string,error) {
	t := typs.TryReadSymbol(d.rd)
	if d.rd.TryError!=nil { return conv2(d.rd.TryError) }
	switch t {
	case t_txt: return conv1(d.readString(txt,128))
	case t_utf: return conv1(d.readString(utf,256))
	case t_bin: return conv1(d.readRaw())
	case t_symbol: return conv3(d.readSymbol())
	}
	panic("...")
}
func (d *Decoder) Read(p []byte) (n int, err error) { return d.rd.Read(p) }
func (d *Decoder) ReadByte() (b byte, err error) { return d.rd.ReadByte() }
func (d *Decoder) Align() (skipped byte) { return d.rd.Align() }


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
	"github.com/maxymania/weirdwire/hufftab"
)

var typs *hufftab.Table
var txt  *hufftab.Table
var utf  *hufftab.Table

const (
	t_txt = iota
	t_utf
	t_bin
	t_symbol
	num_types
)

func init() {
	typs = hufftab.NewTable(num_types)
	typs.IncrOne(t_txt,16)
	typs.IncrOne(t_utf,8)
	typs.IncrOne(t_bin,4)
	typs.Calculate()
	
	txt = hufftab.NewTable(129)
	txt.Incr('A','Z',20)
	txt.Incr('a','z',30)
	txt.Incr('0','9',10)
	txt.Calculate()
	
	utf = hufftab.NewTable(257)
	utf.Incr('A','Z',20)
	utf.Incr('a','z',30)
	utf.Incr('0','9',10)
	utf.Incr(128,255,5)
	utf.Calculate()
}


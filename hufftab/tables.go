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


package hufftab

import (
	"github.com/icza/bitio"
	"github.com/icza/huffman"
)

type TableNode struct{
	huffman.Node
	VBits uint64
	NBits byte
}
func (t *TableNode) Code() (r uint64,bits byte) { return t.VBits,t.NBits }
func (t *TableNode) cacheCode() {
	t.VBits,t.NBits = t.Node.Code()
}

type TableArr []TableNode
func NewTableArr(n int) TableArr {
	a := make(TableArr,n)
	for i := range a {
		a[i].Value = huffman.ValueType(i)
		a[i].Count = 1
	}
	return a
}
func (ta TableArr) IncrOne(i, incr int) {
	ta[i].Count += incr
}
func (ta TableArr) Incr(beg,end, incr int) {
	for i := beg; i<=end; i++ {
		ta[i].Count += incr
	}
}
func (ta TableArr) clone() []*huffman.Node {
	na := make([]*huffman.Node,len(ta))
	for i := range na { na[i] = &(ta[i].Node) }
	return na
}
func (ta TableArr) cacheCode() {
	for i := range ta {
		ta[i].cacheCode()
	}
}

type Table struct{
	TableArr
	*huffman.Node
}
func NewTable(n int) *Table { return &Table{TableArr:NewTableArr(n)} }

func (t *Table) Calculate() {
	t.Node = huffman.Build(t.clone())
	t.cacheCode()
}
func (t *Table) TryReadSymbol(r *bitio.Reader) int {
	cur := t.Node
	for cur.Left != nil {
		if r.TryReadBool() {
			cur = cur.Right
		} else {
			cur = cur.Left
		}
	}
	return int(cur.Value)
}
func (t *Table) TryWriteSymbol(w *bitio.Writer, i int) {
	w.TryWriteBits(t.TableArr[i].Code())
}
func (t *Table) Print() {
	huffman.Print(t.Node)
}


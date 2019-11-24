package cvdh

/*
#cgo CFLAGS: -I/usr/include/
#include <strings.h>
#include <stdio.h>
*/
import "C"
import (
	"errors"
	"github.com/dereklstinson/cutil"
	"github.com/dereklstinson/half"
	"unsafe"
)

//TensorP is a pinned memory tensor
type TensorP struct {
	a     MemManager
	p     *pinnedMem
	dims  []int
	stide []int
	d     DataType
}

//CreateTensorP will create a tensorP.  only datatypes supported are Float16 and Float32
//This is machine learning so...
func CreateTensorP(a MemManager, dims []int, d DataType) (*TensorP, error) {
	if d != Float16 && d != Float32 {
		return nil, errors.New("Unsupported Data Type Passed Float16, Float32 only")
	}
	t := new(TensorP)
	t.a = a
	t.dims = make([]int, len(dims))
	copy(t.dims, dims)
	t.stide = findstride(dims)
	t.p = new(pinnedMem)
	sib := (uint)(findvol(dims) * getbytes(d))
	err := t.a.Allocate(t.p, sib)
	if err != nil {
		return nil, err
	}
	return t, nil
}

//SetValue Assumes you know what you are doing
func (t *TensorP) SetValue(location []int, value float32) {
	var position int
	for i := range location {
		position += location[i] * t.stide[i]
	}
	if t.d == Float16 {
		x := half.NewFloat16(value)
		w, err := cutil.WrapGoMem(&x)
		if err != nil {
			panic(err)
		}

		t.a.Copy(t.p.offset((uint)(position)*2), w, 2)
	}
	w, err := cutil.WrapGoMem(&value)
	if err != nil {
		panic(err)
	}

	t.a.Copy(t.p.offset((uint)(position)*4), w, 4)

}

//pinned mem satisfies cutil.Mem
type pinnedMem struct {
	p        unsafe.Pointer
	sib      uint
	typesize uint
}

func (p *pinnedMem) offset(bybytes uint) cutil.Pointer {
	return cutil.Offset(p, bybytes)
}
func (p *pinnedMem) Ptr() unsafe.Pointer {
	return p.p
}
func (p *pinnedMem) DPtr() *unsafe.Pointer {
	return &p.p
}

//MemManager allocates pinned cpu memory the size of sib to p
type MemManager interface {
	Allocate(p cutil.Mem, sib uint) error
	Copy(dest cutil.Pointer, src cutil.Pointer, sib uint) error
}

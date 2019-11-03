package cvdh

import (
	"errors"
	"math/rand"
	"time"
)

//

//RandResizeHelper helps in getting a random offset point, and random
type RandResizeHelper struct {
	min, max *size
	rng      *rand.Rand
	rngs     rand.Source
}

//CreateResizeHelper creates will return
func CreateResizeHelper(min, max Size) *RandResizeHelper {
	rngs := rand.NewSource(time.Now().Unix())

	return &RandResizeHelper{
		rng:  rand.New(rngs),
		rngs: rngs,
		min:  newsize(min),
		max:  newsize(max),
	}
}

//RandomInputValues returns random input resize values.
func (r *RandResizeHelper) RandomInputValues(input Size) (offset Point, box Size, err error) {
	inputdimsize := input.Size()
	if len(inputdimsize) != len(r.min.s) {
		return nil, nil, errors.New("=input.Size() != r.min.Size() ")
	}
	for i := range inputdimsize {
		if inputdimsize[i] < r.min.s[i] {
			return nil, nil, errors.New("input.Size() < r.min.Size() ")
		}
	}
	mids := make([]int, len(inputdimsize))
	for i := range mids {
		if r.max.s[i] < inputdimsize[i] {
			mids[i] = r.max.s[i]
		} else {
			mids[i] = inputdimsize[i]
		}
	}
	offsetposition := make([]int, len(inputdimsize))
	for i, max := range mids {
		min := r.min.s[i]
		odimsize := inputdimsize[i]
		dimsize := r.rng.Intn(max-min) + min
		offsetposition[i] = r.rng.Intn(odimsize - dimsize)
		mids[i] = dimsize //reusing mids for the new box size
	}
	offset = &point{offsetposition}
	box = &size{mids}
	err = nil
	return offset, box, err
}

/*

func addp(a, b Point) Point {
	return Point{a.X + b.X, a.Y + b.Y}
}
func outofbounds(box Size, p Point) bool {
	if p.X > box.X || p.Y > box.Y {
		return true
	}
	return false
}
func maxpoint(min Point, box Size) (max Point) {
	max = Point{min.X + box.X, min.Y + box.Y}
	return max
}
*/

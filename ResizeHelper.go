package cvdh

import (
	"errors"
	"math/rand"
	"time"
)

//

//RandomHelper will create random values used in image manipulations
type RandomHelper struct {
	min, max *size
	rng      *rand.Rand
	rngs     rand.Source
}

//CreateRandomHelper creates will return
func CreateRandomHelper(min, max Size) *RandomHelper {
	rngs := rand.NewSource(time.Now().Unix())

	return &RandomHelper{
		rng:  rand.New(rngs),
		rngs: rngs,
		min:  newsize(min),
		max:  newsize(max),
	}
}

//ImageResizeOffset returns random input resize values.
func (r *RandomHelper) ImageResizeOffset(input Size) (offset Point, box Size, err error) {
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
		dimsize := 0
		if max-min == 0 {
			dimsize = min
		} else {
			dimsize = r.rng.Intn(max-min) + min
		}
		if odimsize-dimsize == 0 {
			offsetposition[i] = 0
		} else {
			offsetposition[i] = r.rng.Intn(odimsize - dimsize)
		}

		mids[i] = dimsize //reusing mids for the new box size
	}
	//angle =
	offset = &point{offsetposition}
	box = &size{mids}
	err = nil
	return offset, box, err
}

//Bool will return a random bool.
//This can also be used to choose if to mirror or not
func (r *RandomHelper) Bool() bool {
	if r.rng.Int()%2 == 0 {
		return false
	}
	return true
}

//AngleDegrees32 returns a random 360 degree angle with 1/4 resolution
func (r *RandomHelper) AngleDegrees32() float32 {
	return (float32)(r.rng.Intn(360)) + ((float32)(r.rng.Intn(5)))/((float32)(4))
}

//AngleDegrees64 returns a random 360 degree angle with 1/4 resolution
func (r *RandomHelper) AngleDegrees64() float64 {
	return (float64)(r.rng.Intn(360)) + ((float64)(r.rng.Intn(5)))/((float64)(4))
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

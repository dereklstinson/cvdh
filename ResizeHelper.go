package cvdh

import (
	"errors"
	"math/rand"
	"time"
)

//

//RandomHelper will create random values used in image manipulations
type RandomHelper struct {
	min, max []int
	rng      *rand.Rand
	rngs     rand.Source
}

//CreateRandomHelper creates will return
func CreateRandomHelper(min, max []int) *RandomHelper {
	rngs := rand.NewSource(time.Now().Unix())

	return &RandomHelper{
		rng:  rand.New(rngs),
		rngs: rngs,
		min:  min,
		max:  max,
	}
}

//ImageResizeOffset returns random input resize values.
func (r *RandomHelper) ImageResizeOffset(inputdims []int) (offset, box []int, err error) {

	if len(inputdims) != len(r.min) {
		return nil, nil, errors.New("=input.Size() != r.min.Size() ")
	}
	for i := range inputdims {
		if inputdims[i] < r.min[i] {
			return nil, nil, errors.New("input.Size() < r.min.Size() ")
		}
	}
	box = make([]int, len(inputdims))
	for i := range box {
		if r.max[i] < inputdims[i] {
			box[i] = r.max[i]
		} else {
			box[i] = inputdims[i]
		}
	}
	offset = make([]int, len(inputdims))
	for i, max := range box {
		min := r.min[i]
		odimsize := inputdims[i]
		dimsize := 0
		if max-min == 0 {
			dimsize = min
		} else {
			dimsize = r.rng.Intn(max-min) + min
		}
		if odimsize-dimsize == 0 {
			offset[i] = 0
		} else {
			offset[i] = r.rng.Intn(odimsize - dimsize)
		}

		box[i] = dimsize //reusing mids for the new box size
	}
	//angle =

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

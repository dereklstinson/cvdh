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

//ImageResizeOffset returns random crop point values in user's dim order
func (r *RandomHelper) ImageResizeOffset(inputdims []int) (pmin, pmax []int, err error) {

	if len(inputdims) != len(r.min) {
		return nil, nil, errors.New("=input.Size() != r.min.Size() ")
	}
	for i := range inputdims {
		if inputdims[i] < r.min[i] {
			return nil, nil, errors.New("input.Size() < r.min.Size() ")
		}
	}
	pmax = make([]int, len(inputdims))
	for i := range pmax {
		if r.max[i] < inputdims[i] {
			pmax[i] = r.max[i]
		} else {
			pmax[i] = inputdims[i]
		}
	}
	pmin = make([]int, len(inputdims))
	for i, max := range pmax {
		min := r.min[i]
		odimsize := inputdims[i]
		dimsize := min
		if max-min != 0 {
			dimsize = r.rng.Intn(max-min) + min
		}

		if odimsize-dimsize != 0 {
			pmin[i] = r.rng.Intn(odimsize - dimsize)
		}

		pmax[i] = dimsize + pmin[i] //reusing mids for the new box size
	}
	//angle =

	err = nil
	return pmin, pmax, err
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

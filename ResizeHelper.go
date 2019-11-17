package cvdh

import (
	//"errors"
	"math/rand"
	"time"
)

//

//RandomHelper will create random values used in image manipulations
type RandomHelper struct {
	offset              float32
	rng                 *rand.Rand
	rngs                rand.Source
	scalar              float32
	inputsmallerthanmin bool
}

//CreateRandomHelper creats a random helper
//scalar needs to be <1
//This will choose random values between  scalar <= x,y <= 1 to crop the image of x and y
//
//example
//
//image [1080,1920] * (1-(rand.float()*(1-scalar)))      (1-.9) =     1-(.1 * (random number between 0 and 1)) = .9x
//image [1080*.9x,1920*.9x]
// 1080(1-.9x) = max offset
func CreateRandomHelper(scalar float32) *RandomHelper {
	rngs := rand.NewSource(time.Now().Unix())

	return &RandomHelper{
		rng:    rand.New(rngs),
		rngs:   rngs,
		scalar: scalar,
	}
}

//Set sets the min and max
func (r *RandomHelper) Set(scalar float32) {
	r.scalar = scalar
}

//ImageResizeOffset returns random crop point values in user's dim order
func (r *RandomHelper) ImageResizeOffset(inputdims []int) (pmin, pmax []int, err error) {

	pmax = make([]int, len(inputdims))
	pmin = make([]int, len(inputdims))
	for i := range pmax {
		scalar := 1 - (r.rng.Float32() * (1 - r.scalar))
		maxsize := int((float32)(inputdims[i]) * scalar)
		offset := r.rng.Int() % (inputdims[i] - maxsize)
		pmin[i] = offset
		pmax[i] = offset + maxsize
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

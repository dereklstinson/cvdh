package main

import (
	"fmt"
	cvdh "github.com/dereklstinson/cvdhelper"
)

func main() {
	max := ([]int{512, 512, 3})
	min := ([]int{32, 32, 3})
	hlper := cvdh.CreateRandomHelper(min, max)
	iterations := 512
	mininputsize := 64
	for i := 0; i < iterations; i++ {
		size := ([]int{mininputsize + i, mininputsize + i, 3})
		cropmin, cropmax, err := hlper.ImageResizeOffset(size)
		if err != nil {
			panic(err)
		}

		for j := range cropmax {
			if cropmax[j] > mininputsize+i {
				panic("larger than mininputsize + i")
			}
		}
		fmt.Println("CropMin:", cropmin, "; CropMax:", cropmax)
		fmt.Println("Random Bool:", hlper.Bool())
		x := hlper.AngleDegrees32()

		if x > 360 {
			panic(x)
		}
		fmt.Println("Angle:", x)
	}

}

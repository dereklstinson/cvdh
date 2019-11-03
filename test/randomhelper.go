package main

import (
	"fmt"
	cvdh "github.com/dereklstinson/cvdhelper"
)

func main() {
	max := cvdh.CreateSize([]int{512, 512, 3})
	min := cvdh.CreateSize([]int{32, 32, 3})
	hlper := cvdh.CreateRandomHelper(min, max)
	iterations := 512
	mininputsize := 64
	for i := 0; i < iterations; i++ {
		size := cvdh.CreateSize([]int{mininputsize + i, mininputsize + i, 3})
		startpoint, box, err := hlper.ImageResizeOffset(size)
		if err != nil {
			panic(err)
		}
		pnt := startpoint.Point()
		bx := box.Size()
		for j := range pnt {
			if bx[j]+pnt[j] > mininputsize+i {
				panic("larger than mininputsize + i")
			}
		}
		fmt.Println("offset:", pnt, "; Box:", bx)
		fmt.Println("Random Bool:", hlper.Bool())
		x := hlper.AngleDegrees32()

		if x > 360 {
			panic(x)
		}
		fmt.Println("Angle:", x)
	}

}

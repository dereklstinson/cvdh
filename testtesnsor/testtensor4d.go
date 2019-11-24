package main

import (
	"fmt"
	"github.com/anthonynsimon/bild/effect"
	"image/jpeg"
	"image/png"
	"os"
	//	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	cvdh "github.com/dereklstinson/cvdhelper"
	"image"
	"strconv"
)

func main() {
	/*two, err := imgio.Open("hipgopher.png")
	if err != nil {
		panic(err)
	}
	*/
	file, err := os.Open("me.jpg")
	if err != nil {
		panic(err)
	}
	file1, err := os.Open("hipgopher.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	defer file1.Close()
	two, err := jpeg.Decode(file)
	if err != nil {
		panic(err)
	}
	five, err := png.Decode(file1)
	if err != nil {
		panic(err)
	}
	//five, err := imgio.Open("hipgopher.png")
	if err != nil {
		panic(err)
	}
	two2020 := transform.Resize(two, 400, 400, transform.NearestNeighbor)
	five2020 := transform.Resize(five, 400, 400, transform.NearestNeighbor)
	stwo2020 := effect.Sobel(two2020)
	sfive2020 := effect.Sobel(five2020)
	//isfive2020 := effect.Invert(sfive2020)
	//istwo2020 := effect.Invert(stwo2020)
	//err = imgio.Save("resizehip.png", two2020, imgio.PNGEncoder())
	//	err = imgio.Save("resizederek.png", five2020, imgio.PNGEncoder())
	//	err = imgio.Save("resizehipinvsobel.png", stwo2020, imgio.PNGEncoder())
	//	err = imgio.Save("resizederekinvsobel.png", sfive2020, imgio.PNGEncoder())

	//	newtensor := cvdh.Create4dTensor([]image.Image{two2020, five2020, sfive2020, stwo2020}, true)
	reg := []image.Image{two2020, five2020}
	edge := []image.Image{stwo2020, sfive2020}
	//	iedge := []image.Image{istwo2020, isfive2020}
	//	alltogether := [][]image.Image{reg, edge, iedge}
	newtensor := cvdh.CreateBatchTensorFromImageandGrayedEdgeKernel(reg, edge, true)

	fmt.Println("Max,Average:", newtensor.Max(), newtensor.Avg())
	newtensor.Add(-newtensor.Avg())
	newtensor.Divide(newtensor.Max())

	fmt.Println("AfterDivide Min", newtensor.Min())
	fmt.Println("AfterDivide Max", newtensor.Max())
	fmt.Println("AfterAverage", newtensor.Avg())
	//	fmt.Println(newtensor)
	imgs, err := newtensor.ToImages()
	if err != nil {
		panic(err)
	}
	for i := range imgs {
		name := "decode" + strconv.Itoa(i) + ".jpg"
		save, err := os.Create(name)
		if err != nil {
			panic(err)
		}
		defer save.Close()
		//name := "decode" + strconv.Itoa(i) + ".png"
		jpeg.Encode(save, imgs[i], nil)
	}

}

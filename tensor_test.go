package cvdh_test

import (
	"image"
	"image/png"
	"os"
	"strconv"
	"testing"

	cvdh "github.com/dereklstinson/cvdhelper"
)

func TestCreateBatchTensor(t *testing.T) {

	paths, err := cvdh.GetPaths("./testimgs/", []string{".jpg", ".png"})
	if err != nil {
		t.Error(err)
	}
	images := make([]image.Image, 4)
	for i, path := range paths {
		img, err := cvdh.GetImageHD(path)
		if err != nil {
			t.Error(err)
		}
		images[i] = img
	}

	tensor := cvdh.Create4dTensor(images, true)

	avg := tensor.Avg()
	cpy := tensor.Clone()
	avg1 := cpy.Avg()
	if avg != avg1 {
		t.Error("Averages should equal")
	}
	tensor2 := cvdh.Create4dTensor(images, false)
	avg2 := tensor2.Avg()
	if !(avg2 < avg+.02 && avg2 > avg-.02) {
		t.Error("Averages should have been within +- .02", avg, avg2)
	}
	imgs, err := cpy.ToImages()
	if err != nil {
		t.Error(err)
	}
	for i := range imgs {
		file, err := os.Create("./outputfromtesting/imgNCHW" + strconv.Itoa(i) + ".png")

		defer file.Close()
		if err != nil {
			t.Error(err)
		}
		png.Encode(file, imgs[i])
	}
	imgs2, err := tensor2.ToImages()
	if err != nil {
		t.Error(err)
	}
	for i := range imgs2 {
		file, err := os.Create("./outputfromtesting/imgNHWC" + strconv.Itoa(i) + ".png")

		defer file.Close()
		if err != nil {
			t.Error(err)
		}
		png.Encode(file, imgs2[i])
	}

}

func TestTensorOps(t *testing.T) {

	tensor := cvdh.MakeTensor4d([]int{3, 3, 3, 3}, true)

	avg := tensor.Avg()
	if avg != 0 {
		t.Error("average should be 0")
	}
	tensor.Add(1)

	avg = tensor.Avg()
	if avg != 1 {
		t.Error("average should be 1")
	}
	tensor.Add(0)

	avg = tensor.Avg()
	if avg != 1 {
		t.Error("average should be 1")
	}
	tensor.Multiply(4)
	avg = tensor.Avg()
	if avg != 4 {
		t.Error("average should be 4")
	}
	tensor.Multiply(1)
	avg = tensor.Avg()
	if avg != 4 {
		t.Error("average should be 4")
	}
	tensor.Divide(2)
	avg = tensor.Avg()
	if avg != 2 {
		t.Error("average should be 2")
	}
	tensor.Divide(1)
	avg = tensor.Avg()
	if avg != 2 {
		t.Error("average should be 2")
	}
	vol := tensor.Vol()
	if vol != 3*3*3*3 {
		t.Error("Vol should have been 81 ")
	}
	tensor.Place([]int{0, 0, 0, 0}, 200)
	tensor.Place([]int{0, 0, 1, 0}, -200)

	max := tensor.Max()
	min := tensor.Min()
	if max != 200 {
		t.Error("Max should have been 200")
	}
	if min != -200 {
		t.Error("Min should have been -200")
	}

	//Redundant part to test NHWC

	tensor1 := cvdh.MakeTensor4d([]int{4, 4, 4, 4}, false)

	tensor1.Add(1)

	avg1 := tensor1.Avg()
	if avg1 != 1 {
		t.Error("average should be 1")
	}

}

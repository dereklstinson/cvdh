package cvdhelper

import "image"

//BatchedTensor4d is a float32 representation of a 4d tensor
type BatchedTensor4d struct {
	Dims []int
	NCHW bool
	Data []float32
}

//FindMaxHW Will return the max HW
func FindMaxHW(imgs []image.Image) (h, w int) {

	for _, img := range imgs {
		y := img.Bounds().Max.Y
		if h < y {
			h = y
		}
		x := img.Bounds().Max.X
		if w < x {
			w = x
		}
	}
	return h, w
}

//FindMinHW Will return the min HW
func FindMinHW(imgs []image.Image) (h, w int) {

	for _, img := range imgs {
		y := img.Bounds().Max.Y
		if h > y {
			h = y
		}
		x := img.Bounds().Max.X
		if w > x {
			w = x
		}
	}
	return h, w
}

//CreateBatchedTensor creates a tensor from the largest dims found in the img batch it will create black bars on the sides of the positions that don't fit.
//channels is fixed to 3.
func CreateBatchedTensor(imgs []image.Image, NCHW bool) BatchedTensor4d {
	h, w := FindMaxHW(imgs)
	var dims []int
	if NCHW {
		dims = []int{len(imgs), 3, h, w}
	} else {
		dims = []int{len(imgs), h, w, 3}
	}
	hwcvol := findvol([]int{3, h, w})

	data := make([]float32, findvol(dims))
	for i, img := range imgs {
		y := img.Bounds().Max.Y
		x := img.Bounds().Max.X
		hoff := (h - y) / 2
		woff := (w - x) / 2
		if NCHW {
			imgdata := chw(img)
			batchvol := i * hwcvol
			for j := 0; j < 3; j++ {
				dcpos := h * w * j
				scpos := y * x * j
				for k := 0; k < y; k++ {
					dhpos := (hoff + k) * w
					shpos := (k * x)
					for l := 0; l < x; l++ {
						data[(batchvol + dcpos + dhpos + l + woff)] = imgdata[scpos+shpos+l]
					}
				}
			}

		} else {
			imgdata := hwc(img)
			batchvol := i * hwcvol
			for j := 0; j < y; j++ {
				dhpos := 3 * w * (j + hoff)
				shpos := 3 * x * j
				for k := 0; k < x; k++ {
					dwpos := (woff + k) * 3
					swpos := (k * 3)
					for l := 0; l < 3; l++ {
						data[(batchvol + dhpos + dwpos + l)] = imgdata[shpos+swpos+l]
					}
				}
			}
		}
	}
	return BatchedTensor4d{
		Dims: dims,
		Data: data,
		NCHW: NCHW,
	}
}
func chw(a image.Image) []float32 {
	c := 3
	ay := a.Bounds().Max.Y
	ax := a.Bounds().Max.X

	array := make([]float32, ay*ax*c)
	for i := 0; i < ay; i++ {
		for j := 0; j < ax; j++ {
			ra, ga, ba, _ := a.At(j, i).RGBA()

			array[(0*ax*ay)+(i*ax)+j] = float32(ra)
			array[(1*ax*ay)+(i*ax)+j] = float32(ga)
			array[(2*ax*ay)+(i*ax)+j] = float32(ba)

		}
	}
	//65535/x=255 ...x=257
	divideall(float32(257), array)
	return array
}
func hwc(a image.Image) []float32 {
	c := 3
	ay := a.Bounds().Max.Y
	ax := a.Bounds().Max.X

	array := make([]float32, ay*ax*c)
	for i := 0; i < ay; i++ {
		for j := 0; j < ax; j++ {
			ra, ga, ba, _ := a.At(j, i).RGBA()

			array[(i*ax*c)+(j*c)+0] = float32(ra)
			array[(i*ax*c)+(j*c)+1] = float32(ga)
			array[(i*ax*c)+(j*c)+2] = float32(ba)
		}
	}
	//65535/x=255 ...x=257
	divideall(float32(257), array)
	return array
}
func findvol(dims []int) (vol int) {
	vol = 1
	for _, dim := range dims {
		vol *= dim
	}
	return vol
}
func divideall(value float32, array []float32) {
	for i := 0; i < len(array); i++ {
		array[i] = array[i] / value
	}
}

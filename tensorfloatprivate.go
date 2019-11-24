package cvdh

import (
	"image"
	"image/color"
)

func hwgray(a image.Image) []float32 {

	ay := a.Bounds().Max.Y
	ax := a.Bounds().Max.X

	array := make([]float32, ay*ax)
	for i := 0; i < ay; i++ {
		for j := 0; j < ax; j++ {
			ra, ga, ba, _ := a.At(j, i).RGBA()
			avg := float32(ra+ga+ba) / float32(3)
			array[(i*ax)+(j)] = avg

		}
	}
	//65535/x=255 ...x=257
	divideall(float32(257), array)
	return array
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

func divideallint32(value int32, array []int32) {
	for i := 0; i < len(array); i++ {
		array[i] = array[i] / value
	}
}

func divideall(value float32, array []float32) {
	for i := 0; i < len(array); i++ {
		array[i] = array[i] / value
	}
}

func hwctoimage(data []float32, dims []int, nchw bool) image.Image {
	var (
		height  int
		width   int
		channel int
	)
	if nchw {
		channel = dims[0]
		height = dims[1]
		width = dims[2]

	} else {
		height = dims[0]
		width = dims[1]
		channel = dims[2]
	}
	stride := findstride(dims)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			rgba := make([]uint8, channel)
			if nchw {
				for c := 0; c < channel; c++ {
					rgba[c] = uint8(data[(stride[0]*c)+(stride[1]*h)+(stride[2]*w)])
				}
			} else {
				for c := 0; c < channel; c++ {
					rgba[c] = uint8(data[(stride[0]*h)+(stride[1]*w)+(stride[2]*c)])
				}
			}

			if channel == 1 {
				rgb := color.RGBA{R: rgba[0], G: rgba[0], B: rgba[0], A: uint8(255)}
				img.Set(w, h, rgb)
			} else if channel == 2 {
				rgb := color.RGBA{R: rgba[0], G: rgba[1], B: rgba[0]/2 + rgba[1]/2, A: uint8(255)}
				img.Set(w, h, rgb)
			} else if channel == 3 {
				rgb := color.RGBA{R: rgba[0], G: rgba[1], B: rgba[2], A: uint8(255)}
				img.Set(w, h, rgb)
			} else if channel == 4 {
				rgb := color.RGBA{R: rgba[0], G: rgba[1], B: rgba[2], A: rgba[3]}
				img.Set(w, h, rgb)
			}

		}
	}

	return img

}

/*
func chwtoimage(data []float32, dims []int) image.Image {
	channel := dims[0]
	height := dims[1]
	width := dims[2]
	stride := findstride(dims)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			rgba := make([]uint8, channel)
			for c := 0; c < channel; c++ {
				rgba[c] = uint8(data[(stride[0]*c)+(stride[1]*h)+(stride[2]*w)])
			}
			if channel == 1 {
				rgb := color.RGBA{R: rgba[0], G: rgba[0], B: rgba[0], A: uint8(255)}
				img.Set(w, h, rgb)
			} else if channel == 3 {
				rgb := color.RGBA{R: rgba[0], G: rgba[1], B: rgba[2], A: uint8(255)}
				img.Set(w, h, rgb)
			} else if channel == 4 {
				rgb := color.RGBA{R: rgba[0], G: rgba[1], B: rgba[2], A: rgba[3]}
				img.Set(w, h, rgb)
			}

		}
	}

	return img

}
*/

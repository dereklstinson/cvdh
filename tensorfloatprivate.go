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

func hwctoimage(data []float32, dims []int) image.Image {
	h := dims[0]
	w := dims[1]
	c := dims[2]

	img := image.NewRGBA(image.Rect(0, 0, h, w))
	if c == 1 {
		for i := 0; i < h; i++ {
			ht := i * w
			for j := 0; j < w; j++ {
				g := color.Gray{
					Y: uint8(data[ht+j]),
				}
				img.Set(j, i, g)
			}
		}
	} else if c == 3 {

		for i := 0; i < h; i++ {
			ht := i * w * c
			for j := 0; j < w; j++ {
				wh := j * c
				r := uint8(data[ht+wh])
				g := uint8(data[ht+wh+1])
				b := uint8(data[ht+wh+2])
				rgb := color.RGBA{R: r, G: g, B: b, A: uint8(255)}
				img.Set(j, i, rgb)
			}
		}
	} else if c == 4 {
		for i := 0; i < h; i++ {
			ht := i * w * c
			for j := 0; j < w; j++ {
				wh := j * c
				r := uint8(data[ht+wh])
				g := uint8(data[ht+wh+1])
				b := uint8(data[ht+wh+2])
				a := uint8(data[ht+wh+3])
				rgba := color.RGBA{R: r, G: g, B: b, A: a}
				img.Set(j, i, rgba)
			}
		}
	}
	return img

}
func chwtoimage(data []float32, dims []int) image.Image {
	h := dims[1]
	w := dims[2]
	c := dims[0]
	img := image.NewRGBA(image.Rect(0, 0, h, w))
	if c == 1 {
		for i := 0; i < h; i++ {
			ht := i * w
			for j := 0; j < w; j++ {
				g := color.Gray{
					Y: uint8(data[ht+j]),
				}
				img.Set(j, i, g)
			}
		}
	} else if c == 3 {
		chvol := h * w
		for i := 0; i < h; i++ {
			ht := i * w
			for j := 0; j < w; j++ {

				r := uint8(data[chvol*0+ht+j])
				g := uint8(data[chvol*1+ht+j])
				b := uint8(data[chvol*2+ht+j])
				rgb := color.RGBA{R: r, G: g, B: b, A: uint8(255)}
				img.Set(j, i, rgb)
			}
		}
	} else if c == 4 {
		chvol := h * w
		for i := 0; i < h; i++ {
			ht := i * w
			for j := 0; j < w; j++ {

				r := uint8(data[chvol*0+ht+j])
				g := uint8(data[chvol*1+ht+j])
				b := uint8(data[chvol*2+ht+j])
				a := uint8(data[chvol*3+ht+j])
				rgba := color.RGBA{R: r, G: g, B: b, A: a}
				img.Set(j, i, rgba)
			}
		}
	}
	return img

}

package cvdh

import (
	"errors"
	"image"
	"image/color"
)

//Tensor4d is a float32 representation of a 4d tensor
type Tensor4d struct {
	Dims []int
	NCHW bool
	Data []float32
}

//Min returns the minimum value of all data
func (b *Tensor4d) Min() (min float32) {
	min = 99999999
	for _, data := range b.Data {
		if data < min {
			min = data
		}
	}
	return min
}

//Vol returns the volume of the 4d tensor array
func (b *Tensor4d) Vol() int {
	return findvol(b.Dims)
}

//MakeTensor4d makes a zeroed tensor4d
func MakeTensor4d(dims []int, NCHW bool) Tensor4d {
	dims2 := make([]int, len(dims))
	copy(dims2, dims)
	return Tensor4d{
		Data: make([]float32, findvol(dims)),
		Dims: dims2,
		NCHW: NCHW,
	}

}

//Dims64 returns the dims in type int64
func (b *Tensor4d) Dims64() []int64 {
	dims := make([]int64, len(b.Dims))
	for i := range dims {
		dims[i] = int64(b.Dims[i])
	}
	return dims
}

//DimsU64 returns the dims in type uint64
func (b *Tensor4d) DimsU64() []uint64 {
	dims := make([]uint64, len(b.Dims))
	for i := range dims {
		dims[i] = uint64(b.Dims[i])
	}
	return dims
}

//DimsUInt returns the dims in type uint
func (b *Tensor4d) DimsUInt() []uint {
	dims := make([]uint, len(b.Dims))
	for i := range dims {
		dims[i] = uint(b.Dims[i])
	}
	return dims
}

//DimsU32 returns the dims in type uint32
func (b *Tensor4d) DimsU32() []uint32 {
	dims := make([]uint32, len(b.Dims))
	for i := range dims {
		dims[i] = uint32(b.Dims[i])
	}
	return dims
}

//Dims32 returns the dims in type int32
func (b *Tensor4d) Dims32() []int32 {
	dims := make([]int32, len(b.Dims))
	for i := range dims {
		dims[i] = int32(b.Dims[i])
	}
	return dims
}

//ZeroClone a a Tensor4d with the same specs but with zeros in the values.
func (b *Tensor4d) ZeroClone() Tensor4d {
	var a Tensor4d
	a.Data = make([]float32, len(b.Data))

	a.Dims = make([]int, len(b.Dims))
	copy(a.Dims, b.Dims)
	a.NCHW = b.NCHW
	return a

}

//Clone returns a copy the Tensor4d
func (b *Tensor4d) Clone() Tensor4d {
	var a Tensor4d
	a.Data = make([]float32, len(b.Data))
	copy(a.Data, b.Data)
	a.Dims = make([]int, len(b.Dims))
	copy(a.Dims, b.Dims)
	a.NCHW = b.NCHW
	return a

}

//Avg returns the average value of all data
func (b *Tensor4d) Avg() (avg float32) {
	for _, data := range b.Data {
		avg += data
	}
	return avg / float32(len(b.Data))
}

//Max returns the maximum value of all data
func (b *Tensor4d) Max() (max float32) {
	max = -99999999
	for _, data := range b.Data {
		if data > max {
			max = data
		}
	}
	return max
}

//Divide divides all the elements by the value passed
func (b *Tensor4d) Divide(value float32) {
	if value == 1 {
		return
	}
	for i := range b.Data {
		b.Data[i] = b.Data[i] / value
	}
}

//Multiply multiplies all elemenents by value passed
func (b *Tensor4d) Multiply(value float32) {
	if value == 1 {
		return
	}
	for i := range b.Data {
		b.Data[i] = b.Data[i] * value
	}
}

//Add addes all elements by value passed
func (b *Tensor4d) Add(value float32) {
	if value == 0 {
		return
	}
	for i := range b.Data {
		b.Data[i] = b.Data[i] + value
	}
}

//ToImages will convert the batched tensor back into an image.Image. Channel size should be either 1 (gray) or 3 (RGB).
//It will rescale all values to fit between 0 and 255.  This will scale all the values on all of the NCHW or NHWC data as a whole.
//not on individual items ei per HWC or CHW
func (b *Tensor4d) ToImages() ([]image.Image, error) {
	if b.NCHW {
		if !(b.Dims[1] == 1 || b.Dims[1] == 3) {
			return nil, errors.New("Channel Needs to be 1 or 3")
		}
	}
	if !b.NCHW {
		if !(b.Dims[3] == 1 || b.Dims[3] == 3) {
			return nil, errors.New("Channel Needs to be 1 or 3")
		}
	}
	min := b.Min()
	max := b.Max()

	a := b.Clone()
	a.Add(-min)
	a.Multiply(255 / max)

	imgs := make([]image.Image, 0)
	hwcvol := findvol(a.Dims[1:])
	for i := 0; i < a.Dims[0]; i++ {
		data := a.Data[i*hwcvol : (i+1)*hwcvol]
		dims := a.Dims[1:]
		if a.NCHW {
			imgs = append(imgs, chwtoimage(data, dims))
		} else {
			imgs = append(imgs, hwctoimage(data, dims))
		}

	}

	return imgs, nil
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

//Create4dTensorGray creates a tensor from the largest dims found in the img batch it will create black bars on the sides of the positions that don't fit.
//channels is fixed to 1. This also scales the values to 0 to 255.
func Create4dTensorGray(imgs []image.Image, NCHW bool) Tensor4d {
	h, w := FindMaxHW(imgs)
	var dims []int
	if NCHW {
		dims = []int{len(imgs), 1, h, w}
	} else {
		dims = []int{len(imgs), h, w, 1}
	}
	hwcvol := findvol([]int{1, h, w})

	data := make([]float32, findvol(dims))
	for i, img := range imgs {
		y := img.Bounds().Max.Y
		x := img.Bounds().Max.X
		hoff := (h - y) / 2
		woff := (w - x) / 2
		if NCHW {
			imgdata := hwgray(img)
			batchvol := i * hwcvol
			for j := 0; j < 1; j++ {
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
			imgdata := hwgray(img)
			batchvol := i * hwcvol
			for j := 0; j < y; j++ {
				dhpos := 3 * w * (j + hoff)
				shpos := 3 * x * j
				for k := 0; k < x; k++ {
					dwpos := (woff + k) * 3
					swpos := (k * 3)
					for l := 0; l < 1; l++ {
						data[(batchvol + dhpos + dwpos + l)] = imgdata[shpos+swpos+l]
					}
				}
			}
		}
	}
	return Tensor4d{
		Dims: dims,
		Data: data,
		NCHW: NCHW,
	}
}

//Create4dTensor creates a tensor from the largest dims found in the img batch it will create black bars on the sides of the positions that don't fit.
//channels is fixed to 3. This also scales the values to 0 to 255.
func Create4dTensor(imgs []image.Image, NCHW bool) Tensor4d {
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
	return Tensor4d{
		Dims: dims,
		Data: data,
		NCHW: NCHW,
	}
}

//MirroredDim will return the mirrored position given the dim size and position in the dim given.
func MirroredDim(dimsize int, position int) int {
	return (dimsize - 1) - position
}

//MirrorCopy changes each of the batch images into a mirrored reflection
func (b *Tensor4d) MirrorCopy() Tensor4d {
	cpy := b.ZeroClone()
	var (
		n int
		c int
		h int
		w int
	)
	var flipped int
	switch b.NCHW {

	case true:

		n = b.Dims[0]
		c = b.Dims[1]
		h = b.Dims[2]
		w = b.Dims[3]
		batchvol := c * h * w
		chanvol := h * w
		hvol := w
		for i := 0; i < n; i++ {
			for j := 0; j < c; j++ {
				for k := 0; k < h; k++ {
					for l := 0; l < w; l++ {
						flipped = (w - 1) - l
						cpy.Data[(i*batchvol)+(j*chanvol)+(k*hvol)+flipped] = b.Data[(i*batchvol)+(j*chanvol)+(k*hvol)+l]
					}
				}
			}
		}

	case false:

		n = b.Dims[0]
		h = b.Dims[1]
		w = b.Dims[2]
		c = b.Dims[3]
		batchvol := c * h * w
		hvol := c * w
		wvol := c
		for i := 0; i < n; i++ {
			for j := 0; j < h; j++ {
				for k := 0; k < w; k++ {
					flipped = (w - 1) - k
					for l := 0; l < c; l++ {

						cpy.Data[(i*batchvol)+(j*hvol)+(flipped*wvol)+l] = b.Data[(i*batchvol)+(j*hvol)+(k*wvol)+l]
					}
				}
			}
		}
	}
	return cpy
}

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
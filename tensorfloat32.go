package cvdh

import (
	"errors"
	//"fmt"
	"github.com/dereklstinson/half"
	"image"
	"sync"
)

//Tensor4d is a float32 representation of a 4d tensor
type Tensor4d struct {
	dims   []int
	nchw   bool
	stride []int
	data   []float32
}

//Data returns a copy of the data
func (b *Tensor4d) Data() (datacopy []float32) {
	datacopy = make([]float32, len(b.data))
	copy(datacopy, b.data)
	return datacopy
}

//SetTensorValues sets the values of the tensor if values passed are nil. Then the original values will remain.
func (b *Tensor4d) SetTensorValues(d []float32, stride, dims []int) {
	if d != nil {
		b.data = d
	}
	if stride != nil {
		b.stride = stride
	}
	if dims != nil {
		b.dims = dims
	}
	return
}

//DataFP16 returns a copy of the Tensor in FP16
func (b *Tensor4d) DataFP16() (copy []half.Float16) {
	return half.NewFloat16Array(b.data)

}

//Place places the value in the position passed
func (b *Tensor4d) Place(position []int, value float32) {
	if b.stride == nil {
		b.stride = b.Stride()
	}
	pos := 0
	for i := range position {

		pos += b.stride[i] * position[i]
	}

	b.data[pos] = value
}

//Get gets the value at position
func (b *Tensor4d) Get(position []int) (value float32) {
	if b.stride == nil {
		b.stride = b.Stride()
	}
	pos := 0
	for i := range position {

		pos += b.stride[i] * position[i]
	}

	value = b.data[pos]
	return value
}

//TransformTensorCopy returns a copy of the tensor in another format.
//If tensor was in NHWC the it will return one in NCHW.
func (b *Tensor4d) TransformTensorCopy() (cpy *Tensor4d) {
	tdata := make([]float32, b.Vol())
	bstride := b.Stride()
	if b.nchw {

		cpy = &Tensor4d{
			dims: []int{b.dims[0], b.dims[2], b.dims[3], b.dims[1]},
			nchw: false,
			data: tdata,
		}
		tstride := cpy.Stride()

		for i := 0; i < (b.dims[0]); i++ { //n
			for j := 0; j < b.dims[1]; j++ { //c
				for k := 0; k < b.dims[2]; k++ { //h
					for l := 0; l < b.dims[3]; l++ { //w
						cpy.data[i*tstride[0]+(j*tstride[3])+(k*tstride[1])+(l*tstride[2])] = b.data[i*bstride[0]+j*bstride[1]+k*bstride[2]+l*bstride[3]]
					}
				}
			}
		}
		return cpy
	}

	cpy = &Tensor4d{
		dims: []int{b.dims[0], b.dims[3], b.dims[1], b.dims[2]},
		nchw: true,
		data: tdata,
	}
	tstride := cpy.Stride()

	for i := 0; i < (b.dims[0]); i++ { //n
		for j := 0; j < b.dims[1]; j++ { //h
			for k := 0; k < b.dims[2]; k++ { //w
				for l := 0; l < b.dims[3]; l++ { //c
					cpy.data[i*tstride[0]+(l*tstride[1])+(j*tstride[2])+(k*tstride[3])] = b.data[i*bstride[0]+j*bstride[1]+k*bstride[2]+l*bstride[3]]
				}
			}
		}
	}
	return cpy

}

//NCHW if true tensor in NCHW format
func (b *Tensor4d) NCHW() bool {
	return b.nchw
}

//Min returns the minimum value of all data
func (b *Tensor4d) Min() (min float32) {
	min = 99999999
	for _, data := range b.data {
		if data < min {
			min = data
		}
	}
	return min
}

//Stride returns the stride of the tensor offset
func (b *Tensor4d) Stride() (strides []int) {
	strides = make([]int, len(b.dims))
	stride := 1
	for i := range strides {
		strides[i] = stride
		stride *= b.dims[i]
	}
	return strides
}

//Vol returns the volume of the 4d tensor array
func (b *Tensor4d) Vol() int {
	return findvol(b.dims)
}

//MakeTensor4d makes a zeroed tensor4d
func MakeTensor4d(dims []int, NCHW bool) *Tensor4d {
	dims2 := make([]int, len(dims))
	copy(dims2, dims)
	return &Tensor4d{
		data: make([]float32, findvol(dims)),
		dims: dims2,
		nchw: NCHW,
	}

}

/*
func MakeStride4dTensor(dims,strides []int, NCHW bool)*Tensor4d{
	dims2:=make()
}
*/

//Dims64 returns the dims in type int64
func (b *Tensor4d) Dims64() []int64 {
	dims := make([]int64, len(b.dims))
	for i := range dims {
		dims[i] = int64(b.dims[i])
	}
	return dims
}

//DimsU64 returns the dims in type uint64
func (b *Tensor4d) DimsU64() []uint64 {
	dims := make([]uint64, len(b.dims))
	for i := range dims {
		dims[i] = uint64(b.dims[i])
	}
	return dims
}

//DimsUInt returns the dims in type uint
func (b *Tensor4d) DimsUInt() []uint {
	dims := make([]uint, len(b.dims))
	for i := range dims {
		dims[i] = uint(b.dims[i])
	}
	return dims
}

//DimsU32 returns the dims in type uint32
func (b *Tensor4d) DimsU32() []uint32 {
	dims := make([]uint32, len(b.dims))
	for i := range dims {
		dims[i] = uint32(b.dims[i])
	}
	return dims
}

//Dims32 returns the dims in type int32
func (b *Tensor4d) Dims32() []int32 {
	dims := make([]int32, len(b.dims))
	for i := range dims {
		dims[i] = int32(b.dims[i])
	}
	return dims
}

//Dims returns the dims in type int
func (b *Tensor4d) Dims() []int {
	dims := make([]int, len(b.dims))
	copy(dims, b.dims)
	return dims
}

//ZeroClone a a Tensor4d with the same specs but with zeros in the values.
func (b *Tensor4d) ZeroClone() *Tensor4d {
	a := new(Tensor4d)
	a.data = make([]float32, len(b.data))
	a.dims = make([]int, len(b.dims))
	copy(a.dims, b.dims)
	a.nchw = b.nchw
	return a
}

//Clone returns a copy the Tensor4d
func (b *Tensor4d) Clone() *Tensor4d {
	a := new(Tensor4d)
	a.data = make([]float32, len(b.data))
	copy(a.data, b.data)
	a.dims = make([]int, len(b.dims))
	copy(a.dims, b.dims)
	a.nchw = b.nchw
	return a

}

//Avg returns the average value of all data
func (b *Tensor4d) Avg() (avg float32) {
	for _, data := range b.data {
		avg += data
	}
	return avg / float32(len(b.data))
}

//Max returns the maximum value of all data
func (b *Tensor4d) Max() (max float32) {
	max = -99999999
	for _, data := range b.data {
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
	for i := range b.data {
		b.data[i] = b.data[i] / value
	}
}

//Multiply multiplies all elements by value passed
func (b *Tensor4d) Multiply(value float32) {
	if value == 1 {
		return
	}
	for i := range b.data {
		b.data[i] = b.data[i] * value
	}
}

//Add adds all elements by value passed
func (b *Tensor4d) Add(value float32) {
	if value == 0 {
		return
	}
	for i := range b.data {
		b.data[i] = b.data[i] + value
	}
}

//ToImages will convert the batched tensor back into an image.Image. Channel size should be either 1 (gray) or 3 (RGB).
//It will rescale all values to fit between 0 and 255.  This will scale all the values on all of the NCHW or NHWC data as a whole.
//not on individual items ei per HWC or CHW
func (b *Tensor4d) ToImages() ([]image.Image, error) {
	if b.nchw {
		if !(b.dims[1] == 1 || b.dims[1] == 3) {
			return nil, errors.New("Channel Needs to be 1 or 3")
		}
	}
	if !b.nchw {
		if !(b.dims[3] == 1 || b.dims[3] == 3) {
			return nil, errors.New("Channel Needs to be 1 or 3")
		}
	}
	min := b.Min()
	max := b.Max()

	a := b.Clone()
	a.Add(-min)
	a.Multiply(255 / max)

	imgs := make([]image.Image, 0)
	hwcvol := findvol(a.dims[1:])
	for i := 0; i < a.dims[0]; i++ {
		data := a.data[i*hwcvol : (i+1)*hwcvol]
		dims := a.dims[1:]
		if a.nchw {
			imgs = append(imgs, chwtoimage(data, dims))
		} else {
			imgs = append(imgs, hwctoimage(data, dims))
		}

	}

	return imgs, nil
}

//Create4dTensorGray creates a tensor from the largest dims found in the img batch it will create black bars on the sides of the positions that don't fit.
//channels is fixed to 1. This also scales the values to 0 to 255.
func Create4dTensorGray(imgs []image.Image, NCHW bool) *Tensor4d {
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
	return &Tensor4d{
		dims: dims,
		data: data,
		nchw: NCHW,
	}
}

//CreateTensorGrayImageGrayEdgeKernel will take a color image Turn it gray and take an EDK and gray it
//makeing an 2 channel tensor
func CreateTensorGrayImageGrayEdgeKernel(original, edgedetection []image.Image, NCHW bool) *Tensor4d {
	w, h := original[0].Bounds().Max.X, original[0].Bounds().Max.Y
	var dims []int
	batchsize := len(original)
	if NCHW {
		dims = []int{batchsize, 2, h, w}
	} else {
		dims = []int{batchsize, h, w, 2}
	}
	hwcvol := findvol(dims)
	data := make([]float32, hwcvol)
	if NCHW {
		hwc := h * w * 2
		for n := 0; n < batchsize; n++ {
			go func(n int) {
				boffset := hwc * n
				coffset := h * w
				for i := 0; i < h; i++ {
					hoffset := i * w
					for k := 0; k < w; k++ {

						r, g, b, _ := original[n].At(k, i).RGBA()
						re, ge, be, _ := edgedetection[n].At(k, i).RGBA()

						data[boffset+0*coffset+hoffset+k] = (float32)(((r + g + b) / 3) / 257)
						data[boffset+1*coffset+hoffset+k] = (float32)(((re + ge + be) / 3) / 257)
					}
				}

			}(n)

		}

	} else {
		hwc := h * w * 2
		for n := 0; n < batchsize; n++ {
			go func(n int) {
				boffset := hwc * n
				for i := 0; i < h; i++ {
					ioff := i * w * 2
					for k := 0; k < w; k++ {
						koff := k * 2
						r, g, b, _ := original[n].At(k, i).RGBA()
						re, ge, be, _ := edgedetection[n].At(k, i).RGBA()

						data[boffset+ioff+koff+0] = (float32)(((r + g + b) / 3) / 257)
						data[boffset+ioff+koff+1] = (float32)(((re + ge + be) / 3) / 257)

					}
				}
			}(n)
		}
	}
	return &Tensor4d{
		dims: dims,
		data: data,
		nchw: NCHW,
	}
}

//CreateBatchTensorFromImageandGrayedEdgeKernel this would be used for something like edge detection.
//This will make a 4 channel tensor
func CreateBatchTensorFromImageandGrayedEdgeKernel(original, edgedetection []image.Image, NCHW bool) *Tensor4d {
	w, h := original[0].Bounds().Max.X, original[0].Bounds().Max.Y
	var dims []int
	batchsize := len(original)
	if NCHW {
		dims = []int{batchsize, 4, h, w}
	} else {
		dims = []int{batchsize, h, w, 4}
	}
	hwcvol := findvol(dims)
	data := make([]float32, hwcvol)
	if NCHW {
		hwc := h * w * 4
		for n := 0; n < batchsize; n++ {
			go func(n int) {
				boffset := hwc * n
				coffset := h * w
				for i := 0; i < h; i++ {
					hoffset := i * w
					for k := 0; k < w; k++ {

						r, g, b, _ := original[n].At(k, i).RGBA()
						re, ge, be, _ := edgedetection[n].At(k, i).RGBA()
						r /= 257
						g /= 257
						b /= 257
						gray := ((re + ge + be) / 3) / 257

						data[boffset+0*coffset+hoffset+k] = (float32)(r)
						data[boffset+1*coffset+hoffset+k] = (float32)(g)
						data[boffset+2*coffset+hoffset+k] = (float32)(b)
						data[boffset+3*coffset+hoffset+k] = (float32)(gray)
					}
				}

			}(n)

		}

	} else {
		hwc := h * w * 4
		for n := 0; n < batchsize; n++ {
			go func(n int) {
				boffset := hwc * n
				for i := 0; i < h; i++ {
					ioff := i * w * 4
					for k := 0; k < w; k++ {
						koff := k * 4
						r, g, b, _ := original[n].At(k, i).RGBA()
						re, ge, be, ra := edgedetection[n].At(k, i).RGBA()
						var gray uint32
						if ra <= 65535/2 {
							r /= 257
							g /= 257
							b /= 257
							gray = ra / 257
						} else {
							r /= 257
							g /= 257
							b /= 257
							gray = ((re + ge + be) / 3) / 257
							gray = ((re + ge + be) / 3) / 257
						}

						data[boffset+ioff+koff+0] = (float32)(r)
						data[boffset+ioff+koff+1] = (float32)(g)
						data[boffset+ioff+koff+2] = (float32)(b)
						data[boffset+ioff+koff+3] = (float32)(gray)
					}
				}
			}(n)
		}
	}
	return &Tensor4d{
		dims: dims,
		data: data,
		nchw: NCHW,
	}
}

//CreateTensorFromImageandGrayedEdgeKernel this would be used for something like edge detection.
func CreateTensorFromImageandGrayedEdgeKernel(original, edgedetection image.Image, NCHW bool) *Tensor4d {
	w, h := original.Bounds().Max.X, original.Bounds().Max.Y
	var dims []int
	if NCHW {
		dims = []int{1, 3, h, w}
	} else {
		dims = []int{1, h, w, 3}
	}
	hwcvol := findvol(dims)
	data := make([]float32, hwcvol)
	if NCHW {
		coffset := h * w
		for i := 0; i < h; i++ {
			for k := 0; k < w; k++ {
				r, g, b, _ := original.At(k, i).RGBA()
				re, ge, be, _ := edgedetection.At(k, i).RGBA()
				r /= 257
				g /= 257
				b /= 257
				gray := ((re + ge + be) / 3) / 257

				data[0*coffset+i*w+k] = (float32)(r)
				data[1*coffset+i*w+k] = (float32)(g)
				data[2*coffset+i*w+k] = (float32)(b)
				data[3*coffset+i*w+k] = (float32)(gray)
			}
		}

	} else {
		for i := 0; i < h; i++ {
			ioff := i * w * 4
			for k := 0; k < w; k++ {
				koff := k * 4
				r, g, b, _ := original.At(k, i).RGBA()
				re, ge, be, ra := edgedetection.At(k, i).RGBA()
				var gray uint32
				if ra <= 65535/2 {
					r /= 257
					g /= 257
					b /= 257
					gray = ra / 257
				} else {
					r /= 257
					g /= 257
					b /= 257
					gray = ((re + ge + be) / 3) / 257
					gray = ((re + ge + be) / 3) / 257
				}

				data[ioff+koff+0] = (float32)(r)
				data[ioff+koff+1] = (float32)(g)
				data[ioff+koff+2] = (float32)(b)
				data[ioff+koff+3] = (float32)(gray)
			}
		}
	}
	return &Tensor4d{
		dims: dims,
		data: data,
		nchw: NCHW,
	}
}

//Batch1dTensors requires all tensors to have the same dims
func Batch1dTensors(t []*Tensor4d, threads int) (b *Tensor4d) {
	b = new(Tensor4d)

	b.nchw = t[0].nchw
	n := len(t)
	dims := t[0].Dims()
	if dims[0] != 1 {
		panic(dims[0])
	}
	dims[0] = n
	b.dims = dims
	boffset := t[0].Vol()

	b.data = make([]float32, b.Vol())
	var wg sync.WaitGroup
	for i := range t {
		wg.Add(1)
		offset := boffset * i
		go func(offset int) {
			for j := range t[0].data {
				b.data[offset+j] = t[i].data[j]
			}
			wg.Done()
		}(offset)
		if i%threads == threads-1 {
			wg.Wait()
		}
	}
	wg.Wait()
	return b
}

//Create4dTensor creates a tensor from the largest dims found in the img batch it will create black bars on the sides of the positions that don't fit.
//channels is fixed to 3. This also scales the values to 0 to 255.
func Create4dTensor(imgs []image.Image, NCHW bool) *Tensor4d {
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
	return &Tensor4d{
		dims: dims,
		data: data,
		nchw: NCHW,
	}
}

//MirroredDim will return the mirrored position given the dim size and position in the dim given.
func MirroredDim(dimsize int, position int) int {
	return (dimsize - 1) - position
}

//MirrorCopy changes each of the batch images into a mirrored reflection
func (b *Tensor4d) MirrorCopy() *Tensor4d {
	cpy := b.ZeroClone()
	var (
		n int
		c int
		h int
		w int
	)
	var flipped int
	switch b.nchw {

	case true:

		n = b.dims[0]
		c = b.dims[1]
		h = b.dims[2]
		w = b.dims[3]
		batchvol := c * h * w
		chanvol := h * w
		hvol := w
		for i := 0; i < n; i++ {
			for j := 0; j < c; j++ {
				for k := 0; k < h; k++ {
					for l := 0; l < w; l++ {
						flipped = (w - 1) - l
						cpy.data[(i*batchvol)+(j*chanvol)+(k*hvol)+flipped] = b.data[(i*batchvol)+(j*chanvol)+(k*hvol)+l]
					}
				}
			}
		}

	case false:

		n = b.dims[0]
		h = b.dims[1]
		w = b.dims[2]
		c = b.dims[3]
		batchvol := c * h * w
		hvol := c * w
		wvol := c
		for i := 0; i < n; i++ {
			for j := 0; j < h; j++ {
				for k := 0; k < w; k++ {
					flipped = (w - 1) - k
					for l := 0; l < c; l++ {

						cpy.data[(i*batchvol)+(j*hvol)+(flipped*wvol)+l] = b.data[(i*batchvol)+(j*hvol)+(k*wvol)+l]
					}
				}
			}
		}
	}
	return cpy
}

//ConcatTensors concats tensors into a new 4d tensor.  if dest is nil. Function will will allocate new memory and return a pointer to it.
//
//Example :
//var dest *TensorD
//dest=ConcatTensors(srcs,dest)
//
//
// Tensors must all have the same batch, height, and width, Channel size can be different.
//
// Only NCHW for now
func ConcatTensors(tensors []*Tensor4d, dest *Tensor4d) *Tensor4d {

	var (
		pb = tensors[0].dims[0]
		c  int
		ph = tensors[0].dims[2]
		pw = tensors[0].dims[3]
	)
	for i := range tensors {
		if pb != tensors[i].dims[0] {
			return nil
		}
		c += tensors[i].dims[1]
		if ph != tensors[i].dims[2] {
			return nil
		}
		if pw != tensors[i].dims[3] {
			return nil
		}

		if dest == nil {
			dest = MakeTensor4d([]int{pb, c, ph, pw}, true)
		}
		if dest.dims[1] != c {
			return nil
		}
	}
	koffset := 0
	for i := range tensors {
		for j := 0; j < pb; j++ {
			for k := 0; k < tensors[i].dims[1]; k++ {
				for l := 0; l < ph; l++ {
					for m := 0; m < pw; m++ {
						value := tensors[i].Get([]int{j, k, l, m})
						dest.Place([]int{j, k + koffset, l, m}, value)
					}
				}
			}

		}
		koffset += tensors[i].dims[1]
	}
	return dest
}

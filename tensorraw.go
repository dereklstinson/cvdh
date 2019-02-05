package cvdh

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"runtime"

	"github.com/dereklstinson/half"
)

//TensorRaw is a more compressed and accurate Tensor
type TensorRaw struct {
	Dims      []int
	NCHW      bool
	Dtype     DataType
	BigEndian bool
	Data      []byte
}

//DataType is the flags for
type DataType byte

//These are the flags for DataType used in TensorBytes
const (
	Int8    DataType = 1
	UInt8   DataType = 2
	Int16   DataType = 3
	UInt16  DataType = 4
	Int32   DataType = 5
	UInt32  DataType = 6
	Int64   DataType = 7
	UInt64  DataType = 8
	Float16 DataType = 11
	Float32 DataType = 12
	Float64 DataType = 13
)

func getbytes(input DataType) int {
	switch input {

	case Int8:
		return 1
	case UInt8:
		return 1
	case Int16:
		return 2
	case UInt16:
		return 2
	case Int32:
		return 4
	case UInt32:
		return 4
	case Int64:
		return 8
	case UInt64:
		return 8
	case Float16:
		return 2
	case Float32:
		return 4
	case Float64:
		return 8
	}
	return 0
}

//Dims32 returns the dims in Dims32 format
func (t *TensorRaw) Dims32() []int32 {
	dims := make([]int32, len(t.Dims))
	for i := range dims {
		dims[i] = int32(t.Dims[i])
	}
	return dims
}

//MakeTensorRawZeros makes a zero initialzed tensor
func MakeTensorRawZeros(dims []int, dtype DataType, NCHW bool) *TensorRaw {
	runtime.CPUProfile()
	dims1 := make([]int, len(dims))
	copy(dims1, dims)
	vol := findvol(dims)

	return &TensorRaw{
		NCHW:  NCHW,
		Dtype: dtype,
		Dims:  dims1,
		Data:  make([]byte, getbytes(dtype)*vol),
	}
}

//GetSliceSize returls the size you should make the slice
func (t *TensorRaw) GetSliceSize() int {
	x := getbytes(t.Dtype)
	if x == 0 {
		return 0
	}
	return len(t.Data) / x
}

//FillSlice will fill the slice up to the length.  Slice is empty it will allocate memory to it and fill it up.
func (t *TensorRaw) FillSlice(input interface{}) error {

	if t.Dtype != finddatatype(input) {
		return errors.New("Slice Doesn't Match tensor datatype")
	}
	switch x := input.(type) {
	case []int8:
		buf := bytes.NewReader(t.Data)
		if len(x) == 0 {
			size := getbytes(t.Dtype) * len(t.Data)
			x = make([]int8, size)
			if t.BigEndian {
				return binary.Read(buf, binary.BigEndian, &x)
			}
			return binary.Read(buf, binary.LittleEndian, x)

		}
		if t.BigEndian {
			return binary.Read(buf, binary.BigEndian, &x)
		}
		return binary.Read(buf, binary.LittleEndian, x)

	case []uint8:
		buf := bytes.NewReader(t.Data)
		if len(x) == 0 {
			size := getbytes(t.Dtype) * len(t.Data)
			x = make([]uint8, size)
			if t.BigEndian {
				return binary.Read(buf, binary.BigEndian, &x)
			}
			return binary.Read(buf, binary.LittleEndian, x)

		}
		if t.BigEndian {
			return binary.Read(buf, binary.BigEndian, &x)
		}
		return binary.Read(buf, binary.LittleEndian, x)

	case []int16:
		buf := bytes.NewReader(t.Data)
		if len(x) == 0 {
			size := getbytes(t.Dtype) * len(t.Data)
			x = make([]int16, size)
			if t.BigEndian {
				return binary.Read(buf, binary.BigEndian, &x)
			}
			return binary.Read(buf, binary.LittleEndian, x)

		}
		if t.BigEndian {
			return binary.Read(buf, binary.BigEndian, &x)
		}
		return binary.Read(buf, binary.LittleEndian, x)
	case []uint16:
		buf := bytes.NewReader(t.Data)
		if len(x) == 0 {
			size := getbytes(t.Dtype) * len(t.Data)
			x = make([]uint16, size)
			if t.BigEndian {
				return binary.Read(buf, binary.BigEndian, &x)
			}
			return binary.Read(buf, binary.LittleEndian, x)

		}
		if t.BigEndian {
			return binary.Read(buf, binary.BigEndian, &x)
		}
		return binary.Read(buf, binary.LittleEndian, x)
	case []int32:
		buf := bytes.NewReader(t.Data)
		if len(x) == 0 {
			size := getbytes(t.Dtype) * len(t.Data)
			x = make([]int32, size)
			if t.BigEndian {
				return binary.Read(buf, binary.BigEndian, &x)
			}
			return binary.Read(buf, binary.LittleEndian, x)

		}
		if t.BigEndian {
			return binary.Read(buf, binary.BigEndian, &x)
		}
		return binary.Read(buf, binary.LittleEndian, x)
	case []uint32:
		buf := bytes.NewReader(t.Data)
		if len(x) == 0 {
			size := getbytes(t.Dtype) * len(t.Data)
			x = make([]uint32, size)
			if t.BigEndian {
				return binary.Read(buf, binary.BigEndian, &x)
			}
			return binary.Read(buf, binary.LittleEndian, x)

		}
		if t.BigEndian {
			return binary.Read(buf, binary.BigEndian, &x)
		}
		return binary.Read(buf, binary.LittleEndian, x)
	case []int64:
		buf := bytes.NewReader(t.Data)
		if len(x) == 0 {
			size := getbytes(t.Dtype) * len(t.Data)
			x = make([]int64, size)
			if t.BigEndian {
				return binary.Read(buf, binary.BigEndian, &x)
			}
			return binary.Read(buf, binary.LittleEndian, x)

		}
		if t.BigEndian {
			return binary.Read(buf, binary.BigEndian, &x)
		}
		return binary.Read(buf, binary.LittleEndian, x)
	case []uint64:
		buf := bytes.NewReader(t.Data)
		if len(x) == 0 {
			size := getbytes(t.Dtype) * len(t.Data)
			x = make([]uint64, size)
			if t.BigEndian {
				return binary.Read(buf, binary.BigEndian, &x)
			}
			return binary.Read(buf, binary.LittleEndian, x)

		}
		if t.BigEndian {
			return binary.Read(buf, binary.BigEndian, &x)
		}
		return binary.Read(buf, binary.LittleEndian, x)
	case []half.Float16:
		buf := bytes.NewReader(t.Data)
		if len(x) == 0 {
			size := getbytes(t.Dtype) * len(t.Data)
			x = make([]half.Float16, size)
			if t.BigEndian {
				return binary.Read(buf, binary.BigEndian, &x)
			}
			return binary.Read(buf, binary.LittleEndian, x)

		}
		if t.BigEndian {
			return binary.Read(buf, binary.BigEndian, &x)
		}
		return binary.Read(buf, binary.LittleEndian, x)
	case []float32:
		buf := bytes.NewReader(t.Data)
		if len(x) == 0 {
			size := getbytes(t.Dtype) * len(t.Data)
			x = make([]float32, size)
			if t.BigEndian {
				return binary.Read(buf, binary.BigEndian, &x)
			}
			return binary.Read(buf, binary.LittleEndian, x)

		}
		if t.BigEndian {
			return binary.Read(buf, binary.BigEndian, &x)
		}
		return binary.Read(buf, binary.LittleEndian, x)
	case []float64:
		buf := bytes.NewReader(t.Data)
		if len(x) == 0 {
			size := getbytes(t.Dtype) * len(t.Data)
			x = make([]float64, size)
			if t.BigEndian {
				return binary.Read(buf, binary.BigEndian, &x)
			}
			return binary.Read(buf, binary.LittleEndian, x)

		}
		if t.BigEndian {
			return binary.Read(buf, binary.BigEndian, &x)
		}
		return binary.Read(buf, binary.LittleEndian, x)

	}
	return errors.New("Unsupported slice")
}

//FindAverage finds the average byte
func FindAverage(data []TensorRaw) byte {
	var adder uint64
	var counter uint64
	for i := range data {
		for j := range data[i].Data {
			adder += uint64(data[i].Data[j])

			counter++
		}
	}
	return byte(adder / counter)
}

//MakeTensorRaw makes a raw tensor from the datatype passed
func MakeTensorRaw(dims []int, NCHW bool, input interface{}) (*TensorRaw, error) {

	dtype := finddatatype(input)
	if dtype == 0 {
		return nil, errors.New("Not Supported type for input")
	}
	bts, err := makebytes(input, false)
	if err != nil {
		return nil, err
	}
	return &TensorRaw{
		Dtype: dtype,
		Data:  bts,
	}, nil
}

func makebytes(input interface{}, Bigendian bool) ([]byte, error) {

	if Bigendian {
		endian := binary.BigEndian
		buf := new(bytes.Buffer)

		err := binary.Write(buf, endian, input)
		if err != nil {
			return nil, err
		}
		return (buf.Bytes()), nil
	}
	endian := binary.LittleEndian
	buf := new(bytes.Buffer)

	err := binary.Write(buf, endian, input)
	if err != nil {
		return nil, err
	}
	return (buf.Bytes()), nil

}

func finddatatype(input interface{}) DataType {
	switch input.(type) {
	case []int8:
		return Int8
	case []uint8:
		return UInt8
	case []int16:
		return Int16
	case []uint16:
		return UInt16
	case []int32:
		return Int32
	case []uint32:
		return UInt32
	case []int64:
		return Int64
	case []uint64:
		return UInt64
	case []half.Float16:
		return Float16
	case []float32:
		return Float32
	case []float64:
		return Float64

	}

	return 0

}

//CreateRawTensorGrayInt8 creates a tensor from the largest dims found in the img batch it will create black bars on the sides of the positions that don't fit.
//channels is fixed to 1. This also scales the values to 0 to 255.
func CreateRawTensorGrayInt8(imgs []image.Image, NCHW bool) TensorRaw {
	h, w := FindMaxHW(imgs)
	var dims []int
	if NCHW {
		dims = []int{len(imgs), 1, h, w}
	} else {
		dims = []int{len(imgs), h, w, 1}
	}
	hwcvol := findvol([]int{1, h, w})

	data := make([]byte, findvol(dims))
	for i, img := range imgs {
		y := img.Bounds().Max.Y
		x := img.Bounds().Max.X
		hoff := (h - y) / 2
		woff := (w - x) / 2
		if NCHW {
			imgdata := hwgraybyteint8(img)
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
			imgdata := hwgraybyteint8(img)
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
	return TensorRaw{
		Dims:  dims,
		Dtype: UInt8,
		Data:  data,
		NCHW:  NCHW,
	}
}

//FindAverageByte will find the average byte out of all of the tensorRaws
func FindAverageByte(tensors []TensorRaw) (uint8, error) {
	var adder uint32
	var counter uint32
	for i := range tensors {
		if tensors[i].Dtype != UInt8 {
			return 0, errors.New("All Data Type Needs to be uint8")
		}
		for j := range tensors[i].Data {

			adder += uint32(tensors[i].Data[j])
			counter++
		}

	}
	return uint8(adder / counter), nil
}

//Uint8toInt8TensorRawFromAverage changes the data type to Int8 and takes the average away.
func Uint8toInt8TensorRawFromAverage(tensors []TensorRaw, avg uint8) error {
	for i := range tensors {
		if tensors[i].Dtype != UInt8 {
			return errors.New("All Data Type Needs to be uint8")
		}
		tensors[i].Dtype = Int8
		for j := range tensors[i].Data {

			tensors[i].Data[j] -= avg

		}
	}
	return nil
}
func hwgraybyteint8(a image.Image) []byte {
	ay := a.Bounds().Max.Y
	ax := a.Bounds().Max.X

	array := make([]byte, ay*ax)
	for i := 0; i < ay; i++ {
		for j := 0; j < ax; j++ {
			ra, ga, ba, _ := a.At(j, i).RGBA()
			avg := (ra + ga + ba) / (3 * 257)
			avg2 := int8(int32(avg) - int32(128))
			array[(i*ax)+(j)] = byte(avg2)

		}
	}

	return array
}

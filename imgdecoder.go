package cvdh

import (
	"image"
	"image/color"
)

//Encoder helps take images and maps the colors to vectors into a tensor.
type Encoder struct {
	m    map[color.Color][]float32
	dims int
}

//CreateEncoder an array of float32 vectors that will be mapped to an array of color.RGBA.
//These encoders are best used with one hot state.
func CreateEncoder(vectors [][]float32, colors []color.Color) *Encoder {

	m := make(map[color.Color][]float32)
	if len(vectors) != len(colors) {
		return nil
	}
	dims := len(vectors[0])
	for i := range colors {
		if len(vectors[i]) != dims {
			return nil
		}
		m[colors[i]] = vectors[i]
	}
	return &Encoder{
		m:    m,
		dims: dims,
	}
}

//AddToMap will add to the vector.color.RGBA mapping.
func (d *Encoder) AddToMap(vector []float32, c color.RGBA) {

	d.m[c] = vector
}

//Encode encodes the image to what is mapped in the encoder.
//Values that are not mapped will get a vector of all zeros
//
//Try to have a well planned out color scheme to avoid zeros.
//
func (d *Encoder) Encode(img image.Image) *Tensor4d {
	xmax := img.Bounds().Max.X
	ymax := img.Bounds().Max.Y
	stridey := xmax * d.dims
	stridex := d.dims
	outputvol := make([]float32, xmax*ymax*d.dims)
	for i := 0; i < ymax; i++ {
		ypos := i * stridey
		for j := 0; j < xmax; j++ {
			xpos := j * stridex
			pixel := img.At(j, i)
			//r, g, b, a := pixel.RGBA()
			//	r = (r * 255) / a
			//g = (g * 255) / a
			//b = (b * 255) / a
			//	vect := d.m[color.RGBA{(uint8)(r), (uint8)(g), (uint8)(b), (uint8)(255)}]
			vect := d.m[pixel]
			for k := 0; k < len(vect); k++ {
				outputvol[ypos+xpos+k] = vect[k]
			}
		}
	}
	return &Tensor4d{
		data: outputvol,
		dims: []int{1, ymax, xmax, d.dims},
	}
}

//Decoder will store a vector/color mappings
//Mostly works the one hot states.
type Decoder struct {
	maps []dmap
}
type dmap struct {
	vector []float32
	c      color.Color
}

//CreateDecoder creates maps to help decode vectors to colors
func CreateDecoder(vectors [][]float32, colors []color.Color) *Decoder {
	dmaps := make([]dmap, len(colors))
	for i := range colors {
		dmaps[i] = makedmap(vectors[i], colors[i])
	}
	return &Decoder{
		maps: dmaps,
	}
}
func makedmap(v []float32, c color.Color) dmap {
	vect := make([]float32, len(v))
	copy(vect, v)
	return dmap{
		vector: v,
		c:      c,
	}
}
func (d *dmap) checkhotstate(v []float32) float32 {

	for i := range d.vector {
		if d.vector[i] != 0 {
			return d.vector[i] - v[i]
		}
	}
	return 0
}

//Decode - looks at the one hot state high one.  takes the difference with that in the vector.
// The one with the highest value gets that color in the map.
//works for only NHWC for now.
func (d *Decoder) Decode(tensor *Tensor4d) []image.Image {
	if tensor.nchw {
		return nil
	}
	b := tensor.dims[0]
	h := tensor.dims[1]
	w := tensor.dims[2]
	c := tensor.dims[3]
	strides := tensor.Stride()
	multipleimages := make([]image.Image, b)

	for n := 0; n < b; n++ {
		im := image.NewRGBA(image.Rect(0, 0, w, h))
		for i := 0; i < h; i++ {
			for j := 0; j < w; j++ {
				vector := make([]float32, c)
				for k := 0; k < c; k++ {
					vector[k] = tensor.data[strides[0]*n+strides[1]*h+strides[2]*w+strides[3]*c]
				}
				var best = float32(-9999999)
				var position int
				for k := range d.maps {
					current := d.maps[k].checkhotstate(vector)
					if best < current {
						best = current
						position = k
					}
				}
				im.Set(j, i, d.maps[position].c)
			}
		}
		multipleimages[n] = im
	}
	return multipleimages

}

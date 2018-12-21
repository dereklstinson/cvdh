package cvdh

import "image"

//FindMaxHW Will return the max h and w
func FindMaxHW(imgs []image.Image) (h, w int) {
	h, w = -99999999, -999999999
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

//FindMinHW Will return the min h and w
func FindMinHW(imgs []image.Image) (h, w int) {
	h, w = 9999999999, 9999999999
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

//FindAvgHW returns the average h and w.
func FindAvgHW(imgs []image.Image) (h, w int) {
	for _, img := range imgs {
		h += img.Bounds().Max.Y

		w += img.Bounds().Max.X

	}
	return h / len(imgs), w / len(imgs)
}

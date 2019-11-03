package cvdh

import "image"

//FindImageStats will return some statistics of the images Size will return []int{h,w}
func FindImageStats(imgs []image.Image) (max, min, mean, mode, median []int) {
	maxh, maxw := -99999999, -999999999
	minh, minw := 9999999999, 9999999999
	avgh, avgw := 0, 0
	numofsizesarrayX := make([]int, 0)
	numofsizesarrayY := make([]int, 0)
	for _, img := range imgs {
		y := img.Bounds().Max.Y
		x := img.Bounds().Max.X
		avgh += y
		avgw += x
		if maxh < y {
			maxh = y
		}
		if maxw < x {
			maxw = x
		}
		if minh > y {
			minh = y
		}
		if minw < x {
			minw = x
		}
		if len(numofsizesarrayX) <= x {
			appendx := make([]int, x-len(numofsizesarrayX)+1)
			numofsizesarrayX = append(numofsizesarrayX, appendx...)
		}
		if len(numofsizesarrayY) <= y {
			appendy := make([]int, y-len(numofsizesarrayY)+1)
			numofsizesarrayY = append(numofsizesarrayY, appendy...)
		}
		numofsizesarrayX[x]++
		numofsizesarrayY[y]++
	}
	maxinarray := -9999
	modey := -1
	offset := 0
	medianx := 0
	for i, sizei := range numofsizesarrayY {
		offset += sizei
		if sizei > maxinarray {
			maxinarray = sizei
			modey = i
		}
		if offset >= len(imgs)/2 {
			medianx = i
		}
	}
	maxinarray = -9999
	offset = 0
	mediany := 0
	modex := -1
	for i, sizei := range numofsizesarrayX {
		offset += sizei
		if sizei > maxinarray {
			maxinarray = sizei
			modex = i
		}
		if offset >= len(imgs)/2 {
			mediany = i
		}
	}
	max = ([]int{maxh, maxw})
	min = ([]int{minh, minw})
	mean = ([]int{avgh / len(imgs), avgw / len(imgs)})
	median = ([]int{mediany, medianx})
	mode = ([]int{modey, modex})
	return max, min, mean, mode, median

}

//FindMaxSize returns Size in row major [h,w]
func FindMaxSize(imgs []image.Image) []int {
	h, w := FindMaxHW(imgs)
	return ([]int{h, w})
}

//FindMinSize returns Size in row major [h,w]
func FindMinSize(imgs []image.Image) []int {
	h, w := FindMinHW(imgs)
	return ([]int{h, w})
}

//FindAvgSize returns Size in row major [h,w]
func FindAvgSize(imgs []image.Image) []int {
	h, w := FindAvgHW(imgs)
	return ([]int{h, w})
}

//FindMaxHW Will return the max h and w
//This is old but it is staying here to not break backwardscompatability
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
//This is old but it is staying here to not break backwardscompatability
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
//This is old but it is staying here to not break backwardscompatability
func FindAvgHW(imgs []image.Image) (h, w int) {
	for _, img := range imgs {
		h += img.Bounds().Max.Y

		w += img.Bounds().Max.X

	}
	return h / len(imgs), w / len(imgs)
}

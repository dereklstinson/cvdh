package cvdh

import (
	"fmt"
	"image"
	"sync"
)

//FolderImageInfo is the folder info
type FolderImageInfo struct {
	Max        []int       `json:"max,omitempty"`
	Min        []int       `json:"min,omitempty"`
	Mean       []int       `json:"mean,omitempty"`
	Mode       []int       `json:"mode,omitempty"`
	Median     []int       `json:"median,omitempty"`
	Folder     string      `json:"folder,omitempty"`
	AvgPixel   uint        `json:"avg_pixel,omitempty"`
	TotalPixel uint        `json:"total_pixel,omitempty"`
	List       []ImageInfo `json:"list,omitempty"`
}

//ImageInfo is the image info
type ImageInfo struct {
	Size []int  `json:"size,omitempty"`
	Name string `json:"name,omitempty"`
}

//GetSizesandAveragePixelandTotalPixel gets all image info probably wanted.
func GetSizesandAveragePixelandTotalPixel(paths []string, threads int) (sizes [][]int, avgpixel uint, totalpixel uint) {
	sizes = make([][]int, len(paths))
	var wg sync.WaitGroup
	var totalpixelvalue uint

	for i := 0; i < len(paths); i += threads {
		totalthreads := threads
		if i+threads >= len(paths) {
			totalthreads = len(paths) - i

		}
		totalpixelarray := make([]uint, totalthreads)
		averagepixelarray := make([]uint, totalthreads)

		for j, ti := i, 0; j < totalthreads+i; j, ti = j+1, ti+1 {
			path := paths[j]
			sizes[j] = make([]int, 2)
			wg.Add(1)

			go func(path string, j int, ti int) {
				var avg uint
				img, err := GetImageHD(path)
				if err != nil {
					fmt.Println("Skipping BadFile:", path)
					sizes[j][0], sizes[j][1] = -1, -1
				} else {
					y, x := img.Bounds().Max.Y, img.Bounds().Max.X
					sizes[j][0], sizes[j][1] = y, x
					for k := 0; k < y; k++ {
						for l := 0; l < x; l++ {
							r, g, b, _ := img.At(l, k).RGBA()
							avg += uint(r+g+b) / (257)
						}
					}
					totalpixelarray[ti] = uint(x * y)
					averagepixelarray[ti] = avg
				}

				wg.Done()
			}(path, j, ti)
		}
		wg.Wait()
		for j := range averagepixelarray {
			totalpixelvalue += averagepixelarray[j]
			totalpixel += totalpixelarray[j]
		}
	}
	avgpixel = totalpixelvalue / totalpixel
	avgpixel /= 3
	return sizes, avgpixel, totalpixel
}

//GetStats gets image size stats
func GetStats(sizes [][]int) (max, min, mean, mode, median []int) {
	maxh, maxw := -99999999, -999999999
	minh, minw := 9999999999, 9999999999
	avgh, avgw := 0, 0
	numofsizesarrayX := make([]int, 0)
	numofsizesarrayY := make([]int, 0)
	var baddcounter int
	for i := range sizes {
		y := sizes[i][0]
		x := sizes[i][1]
		if y < 1 || x < 1 {
			baddcounter--
		} else {
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
			if minw > x {
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
	}
	maxinarray := -9999
	modey := -1
	offset := 0
	medianx := 0
	totalcounted := len(sizes) - baddcounter
	for i, sizei := range numofsizesarrayY {
		offset += sizei
		if sizei > maxinarray {
			maxinarray = sizei
			modey = i
		}
		if offset >= totalcounted/2 && medianx == 0 {
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
		if offset >= totalcounted/2 && mediany == 0 {
			mediany = i
		}
	}

	max = ([]int{maxh, maxw})
	min = ([]int{minh, minw})
	mean = ([]int{avgh / totalcounted, avgw / totalcounted})
	median = ([]int{mediany, medianx})
	mode = ([]int{modey, modex})
	return max, min, mean, mode, median
}

//FindImageStatsFromPaths will return image stats from all the paths sent.
func FindImageStatsFromPaths(paths []string, threads int) (max, min, mean, mode, median []int) {
	maxh, maxw := -99999999, -999999999
	minh, minw := 9999999999, 9999999999
	avgh, avgw := 0, 0
	numofsizesarrayX := make([]int, 0)
	numofsizesarrayY := make([]int, 0)
	var baddcounter int
	var mtx sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < len(paths); i += threads {
		totalthreads := threads
		if i+threads >= len(paths) {
			totalthreads = len(paths) - i

		}
		for j := i; j < totalthreads+i; j++ {

			path := paths[j]
			wg.Add(1)
			go func(path string) {
				img, err := GetImageHD(path)
				if err != nil {
					fmt.Println("Skipping BadFile:", path)
					mtx.Lock()
					baddcounter++
					mtx.Unlock()
				} else {

					y := img.Bounds().Max.Y
					x := img.Bounds().Max.X
					mtx.Lock()
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
					if minw > x {
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
					mtx.Unlock()
				}
				wg.Done()
			}(path)

		}
		wg.Wait()

	}

	maxinarray := -9999
	modey := -1
	offset := 0
	medianx := 0
	totalcounted := len(paths) - baddcounter
	for i, sizei := range numofsizesarrayY {
		offset += sizei
		if sizei > maxinarray {
			maxinarray = sizei
			modey = i
		}
		if offset >= totalcounted/2 {
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
		if offset >= totalcounted/2 {
			mediany = i
		}
	}

	max = ([]int{maxh, maxw})
	min = ([]int{minh, minw})
	mean = ([]int{avgh / totalcounted, avgw / totalcounted})
	median = ([]int{mediany, medianx})
	mode = ([]int{modey, modex})
	return max, min, mean, mode, median
}

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

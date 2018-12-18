package cvdhelper

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

//GetPaths will return the paths of all the files in the directory and sub directorys. If suffixes is nil then it will grab everything.
func GetPaths(folder string, withsuffixes []string) (paths []string, err error) {
	paths = make([]string, 0)
	err = filepath.Walk(folder, func(path string, info os.FileInfo, errw error) error {
		if errw != nil {
			return errw
		}
		if info.IsDir() {
			return nil
		}
		if withsuffixes == nil {
			paths = append(paths, path)
		} else {
			for i := range withsuffixes {
				if strings.HasSuffix(path, withsuffixes[i]) {
					paths = append(paths, path)
				}

			}
		}

		return nil
	})
	return paths, err
}

//GetImageHD will return an image.Image based on the prefix.  Supported files are .jpeg, jpg, png.
func GetImageHD(path string) (image.Image, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(path, ".jpeg") || strings.HasSuffix(path, ".jpg") {
		return jpeg.Decode(file)
	} else if strings.HasSuffix(path, ".png") {
		return png.Decode(file)
	}
	return nil, errors.New("Unsupported Format")
}

/*

package cvdhelper

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//GetPaths will return the paths of all the files in the directory and sub directorys. If suffixes is nil then it will grab everything.
func GetPaths(folder string, withsuffixes []string) (paths []string, err error) {
	paths = make([]string, 0)
	err = filepath.Walk(folder, func(path string, info os.FileInfo, errw error) error {
		if errw != nil {
			return errw
		}
		if info.IsDir() {
			return nil
		}
		if withsuffixes == nil {
			paths = append(paths, path)
		} else {
			for i := range withsuffixes {
				if strings.HasSuffix(path, withsuffixes[i]) {
					paths = append(paths, path)
				}

			}
		}

		return nil
	})
	return paths, err
}

//GetImageHD will return an image.Image based on the prefix.  Supported files are .jpeg, jpg, png.
func GetImageHD(path string) (image.Image, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(path, ".jpeg") || strings.HasSuffix(path, ".jpg") {
		return jpeg.Decode(file)
	} else if strings.HasSuffix(path, ".png") {
		return png.Decode(file)
	}
	return nil, errors.New("Unsupported Format")
}

//GetImagesHD will return a slice of image.Image from hard drive.  This can be parallelized with goroutines >1.
func GetImagesHD(paths []string, goroutines int) ([]image.Image, error) {
	if goroutines < 1 {
		return nil, errors.New("Must have  goroutines> 0")
	}
	if len(paths) < goroutines {
		return nil, errors.New("number of paths can't be less than number of goroutines")
	}

	imgs := make([]image.Image, len(paths))
	if goroutines == 1 {
		var err error
		for i, path := range paths {
			imgs[i], err = GetImageHD(path)
			if err != nil {
				return nil, err
			}
		}
	}
	imgchans := make([]chan image.Image, goroutines)

	err := goranges(imgs, paths, imgchans)
	return imgs, err
}
func goranges(imgs []image.Image, paths []string, imgchans []chan image.Image) error {
	offset := len(imgs) / len(imgchans)
	errs := make([]chan int, len(imgchans))
	var wg sync.WaitGroup
	for i, imgchan := range imgchans {

		errs[i] = make(chan int, 1)
		imgchan = make(chan image.Image, 10)

		wg.Add(1)
		startat := i * offset
		endat := startat + offset - 1
		if endat < len(imgs) {
			go channelsitter(imgs[startat:endat], imgchan, &wg)
			go imagefinder(imgchan, errs[i], paths[startat:endat])

		} else {
			go channelsitter(imgs[startat:], imgchan, &wg)
			go imagefinder(imgchan, errs[i], paths[startat:])

		}

	}
	for _, err := range errs {
		x := <-err
		if x != 0 {
			return errors.New("1")
		}
	}

	wg.Wait() //this is probably redundant
	return nil
}
func imagefinder(imgchan chan<- image.Image, errchan chan<- int, paths []string) {
	for _, path := range paths {
		img, err := GetImageHD(path)

		if err != nil {
			close(imgchan)
			errchan <- 1 //panic(err)
		}

		imgchan <- img

	}
	fmt.Println("Didn't Find an Error closing image chan")
	close(imgchan)
	fmt.Println("Sending a zero")
	errchan <- 0
}
func channelsitter(imgs []image.Image, imgchan <-chan image.Image, wg *sync.WaitGroup) {
	i := 0
	for img := range imgchan {
		imgs[i] = img
		i++
	}
	fmt.Println("Sending Done for waitgroup")
	wg.Done()
}


//GetImage takes an io.Reader and first checks to see if it is a jpeg. If not it will check to see if it is an png.
func GetImage(r io.Reader) (image.Image, error) {
	img, err1 := jpeg.Decode(r)
	if err1 != nil {
		strerr := err1.Error()
		img1, err2 := png.Decode(r)
		if err2 != nil {
			strerr2 := err2.Error()
			return nil, errors.New("JPEG: " + strerr + ", PNG: " + strerr2 + " ---Unsupported Format")
		}
		return img1, err2
	}
	return img, err1
}
*/

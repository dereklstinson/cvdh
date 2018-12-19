package cvdh

import (
	"errors"
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

	if goroutines == 1 {
		imgs := make([]image.Image, len(paths))
		var err error
		for i, path := range paths {
			imgs[i], err = GetImageHD(path)
			if err != nil {
				return nil, err
			}
		}
		return imgs, nil
	}
	imgchans := make([]chan imgmessage, goroutines)

	return goranges(paths, imgchans)

}
func goranges(paths []string, imgchans []chan imgmessage) ([]image.Image, error) {
	offset := len(paths) / len(imgchans)
	imgs := make([][]image.Image, len(imgchans))
	errs := make([]error, len(imgchans))
	var wg sync.WaitGroup
	for i, imgchan := range imgchans {
		wg.Add(1)
		imgchan = make(chan imgmessage, 2)
		startat := i * offset
		endat := startat + offset //- 1

		go func(imgchan chan imgmessage, i int) {
			msg := <-imgchan

			if msg.err != nil {
				errs[i] = msg.err
			}
			imgs[msg.worker] = msg.imgs
			wg.Done()
		}(imgchan, i)

		if endat < len(imgs) {

			go imagefinder(imgchan, i, paths[startat:endat])

		} else {

			go imagefinder(imgchan, i, paths[startat:])

		}

	}

	wg.Wait()
	for _, e := range errs {
		if e != nil {
			return nil, e
		}
	}
	imgs2 := make([]image.Image, 0)
	for i := range imgs {
		imgs2 = append(imgs2, imgs[i]...)
	}

	return imgs2, nil
}
func imagefinder(imgchan chan<- imgmessage, index int, paths []string) {
	imgs := make([]image.Image, len(paths))
	for i, path := range paths {
		img, err := GetImageHD(path)

		if err != nil {
			imgchan <- imgmessage{
				worker: index,
				err:    err,
			}
			close(imgchan)

		}

		imgs[i] = img

	}
	imgchan <- imgmessage{
		worker: index,
		imgs:   imgs,
	}
	close(imgchan)

}

type imgmessage struct {
	worker int
	imgs   []image.Image
	err    error
}

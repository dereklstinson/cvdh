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

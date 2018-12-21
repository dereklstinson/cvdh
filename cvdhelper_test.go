package cvdh_test

import (
	"testing"

	cvdh "github.com/dereklstinson/cvdhelper"
)

//asdf
func TestGetImages(t *testing.T) {
	paths, err := cvdh.GetPaths("./testimgs/", []string{".jpg", ".png"})
	if err != nil {
		t.Error(err)
	}
	threads := len(paths) /*don't do this. I am only doing this because I have 4 images I am testing.
	If you have 100,000 images the program will crash.  Use this threads=min(number of cputhreads,len(images))
	*/

	imgs, err := cvdh.GetImagesHD(paths, threads)
	if err != nil {
		t.Error(err)
	}

	if len(imgs) != len(paths) {
		t.Error("Length of imgs and paths needs to be the same", len(imgs), len(paths))
	}
	minh, minw := cvdh.FindMinHW(imgs)
	if minh == 0 && minw == 0 {
		t.Error(minh, minw)
	}
	//Redundacy
	imgs, err = cvdh.GetImagesHD(paths, 1)
	if err != nil {
		t.Error(err)
	}

	if len(imgs) != len(paths) {
		t.Error("Length of imgs and paths needs to be the same", len(imgs), len(paths))
	}
	minh, minw = cvdh.FindMinHW(imgs)
	if minh == 0 && minw == 0 {
		t.Error(minh, minw)
	}
	avgh, avgw := cvdh.FindMinHW(imgs)
	if avgh == 0 && avgw == 0 {
		t.Error(avgh, avgw)
	}
}

func TestGetPaths(t *testing.T) {
	paths, err := cvdh.GetPaths("./testimgs/", []string{".png"})
	if err != nil {
		t.Error(err)
	}
	if len(paths) != 3 {
		t.Error("Should Be three")
	}
	paths, err = cvdh.GetPaths("./testimgs/", nil)
	if err != nil {
		t.Error(err)
	}
	if len(paths) != 5 {
		t.Error("Should Be five")
	}
	paths, err = cvdh.GetPaths("./testimgs/", []string{".jpg"})
	if err != nil {
		t.Error(err)
	}
	if len(paths) != 1 {
		t.Error("Should Be one")
	}
	paths, err = cvdh.GetPaths("./testimgs/", []string{".jpg", ".png"})
	if err != nil {
		t.Error(err)
	}
	if len(paths) != 4 {
		t.Error("Should Be four")
	}
}

func TestGetImageHD(t *testing.T) {
	paths, err := cvdh.GetPaths("./testimgs/", []string{".jpg", ".png"})
	if err != nil {
		t.Error(err)
	}
	for _, path := range paths {
		img, err := cvdh.GetImageHD(path)
		if err != nil {
			t.Error(err)
		}
		if img == nil {
			t.Error("Shouldn't Be nil")
		}
	}
	paths, err = cvdh.GetPaths("./testimgs/", []string{".mat"})
	if err != nil {
		t.Error(err)
	}
	for _, path := range paths {
		img, err := cvdh.GetImageHD(path)
		if err == nil {
			t.Error("This should have been nil")
		}
		if img != nil {
			t.Error("This should be nil")
		}
	}
}

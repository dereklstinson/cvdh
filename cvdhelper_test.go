package cvdhelper_test

import (
	"testing"

	"github.com/dereklstinson/cvdhelper"
)

const mpiiimgaelocation = "/home/derek/Desktop/mpii-images/"

func TestGetImages(t *testing.T) {
	paths, err := cvdhelper.GetPaths(mpiiimgaelocation, []string{"2.jpg"})
	if err != nil {
		t.Error(err)
	}
	_, err = cvdhelper.GetImages(paths, 24)
	if err != nil {
		t.Error(err)
	}

}

func TestGetPaths(t *testing.T) {
	paths, err := cvdhelper.GetPaths("./testimgs/", []string{".png"})
	if err != nil {
		t.Error(err)
	}
	if len(paths) != 3 {
		t.Error("Should Be three")
	}
	paths, err = cvdhelper.GetPaths("./testimgs/", nil)
	if err != nil {
		t.Error(err)
	}
	if len(paths) != 5 {
		t.Error("Should Be five")
	}
	paths, err = cvdhelper.GetPaths("./testimgs/", []string{".jpg"})
	if err != nil {
		t.Error(err)
	}
	if len(paths) != 1 {
		t.Error("Should Be one")
	}
	paths, err = cvdhelper.GetPaths("./testimgs/", []string{".jpg", ".png"})
	if err != nil {
		t.Error(err)
	}
	if len(paths) != 4 {
		t.Error("Should Be four")
	}
}

func TestGetImageHD(t *testing.T) {
	paths, err := cvdhelper.GetPaths("./testimgs/", []string{".jpg", ".png"})
	if err != nil {
		t.Error(err)
	}
	for _, path := range paths {
		img, err := cvdhelper.GetImageHD(path)
		if err != nil {
			t.Error(err)
		}
		if img == nil {
			t.Error("Shouldn't Be nil")
		}
	}
	paths, err = cvdhelper.GetPaths("./testimgs/", []string{".mat"})
	if err != nil {
		t.Error(err)
	}
	for _, path := range paths {
		img, err := cvdhelper.GetImageHD(path)
		if err == nil {
			t.Error("This should have been nil")
		}
		if img != nil {
			t.Error("This should be nil")
		}
	}
}

/*
func TestGetImage(t *testing.T) {
	paths, err := cvdhelper.GetPaths("./testimgs/", []string{".jpg", ".png"})
	if err != nil {
		t.Error(err)
	}
	for _, path := range paths {
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			t.Error(err)
		}
		img, err := cvdhelper.GetImage(file)
		if err != nil {
			t.Error(err)
		}
		if img == nil {
			t.Error("Shouldn't Be nil")
		}
	}
}
*/

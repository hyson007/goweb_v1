package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//Image is NOT stored in DB
type Image struct {
	GalleryID uint
	Filename  string
}

// func (i *Image) String() string {
// 	// this makes sure the old code can still works
// 	// but seems no luck...
// 	return i.Path()
// }

func (i *Image) Path() string {
	return "/" + i.RelativePath()
}

func (i *Image) RelativePath() string {
	return fmt.Sprintf("images/galleries/%v/%v", i.GalleryID, i.Filename)
}

type ImageService interface {
	CreateImage(galleryID uint, r io.ReadCloser, filename string) error
	ByGalleryID(gallery uint) ([]Image, error)
	Delete(i *Image) error
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	//grab every files under that gallery directory3
	galleryPath := is.imagePath(galleryID)
	imgStrings, err := filepath.Glob(galleryPath + "*")
	ret := make([]Image, len(imgStrings))
	if err != nil {
		return nil, err
	}
	for i := range imgStrings {
		imgStrings[i] = strings.Replace(imgStrings[i], galleryPath, "", 1)
		ret[i] = Image{
			GalleryID: galleryID,
			Filename:  imgStrings[i],
		}
	}
	return ret, nil
}

func (is *imageService) imagePath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

func (is *imageService) CreateImage(galleryID uint, r io.ReadCloser, filename string) error {
	defer r.Close()
	path, err := is.mkImagePath(galleryID)
	if err != nil {
		return err
	}

	//f.Filename comes from form
	//create dst file
	dst, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	//copy reader data over to dst file
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}

	return nil
}

func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}

func (is *imageService) Delete(i *Image) error {
	// delete from file system
	return os.Remove(i.RelativePath())
}

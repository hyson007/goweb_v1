package models

import (
	"github.com/jinzhu/gorm"
)

type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null; index"`
	Title  string `gorm:"not_null"`
}

type GalleryService interface {
	GalleryDB
}

type galleryService struct {
	GalleryDB
}

type GalleryDB interface {
	Create(gallery *Gallery) error
}

func NewGalleryservice(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{db: db},
		},
	}
}

type galleryValidator struct {
	GalleryDB
}

// if we just create like userDB, we will end up with two DB connections
type galleryGorm struct {
	db *gorm.DB
}

//make sure galleryGorm implement GalleryDB
var _ GalleryDB = &galleryGorm{}

type GalleryValFunc func(*Gallery) error

func runGalleryValFuncs(g *Gallery, fns ...GalleryValFunc) error {
	for _, fn := range fns {
		if err := fn(g); err != nil {
			return err
		}
	}
	return nil
}

func (gv *galleryValidator) UserIDRequired(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) TitleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (gv *galleryValidator) Create(g *Gallery) error {
	// thi varadic will work even without passing any function
	if err := runGalleryValFuncs(g,
		gv.TitleRequired,
		gv.UserIDRequired,
	); err != nil {
		return err
	}
	return gv.GalleryDB.Create(g)
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

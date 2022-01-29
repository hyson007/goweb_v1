package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Gallery struct {
	gorm.Model
	UserID uint    `gorm:"not_null; index"`
	Title  string  `gorm:"not_null"`
	Images []Image `gorm:"-"`
}

func (g Gallery) ImageSplit(n int) [][]Image {
	ret := make([][]Image, n)
	for i := 0; i < n; i++ {
		ret[i] = make([]Image, 0)
	}
	for i, img := range g.Images {
		bucket := i % n
		ret[bucket] = append(ret[bucket], img)
	}
	return ret
}

type GalleryService interface {
	GalleryDB
}

type galleryService struct {
	GalleryDB
}

type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	ByUserID(userID uint) ([]Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
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

func (gv *galleryValidator) Update(g *Gallery) error {
	// thi varadic will work even without passing any function
	if err := runGalleryValFuncs(g,
		gv.TitleRequired,
		gv.UserIDRequired,
	); err != nil {
		return err
	}
	return gv.GalleryDB.Update(g)
}

// delete user
func (gv *galleryValidator) Delete(id uint) error {
	if id < 0 {
		return ErrInvalidID
	}
	return gv.GalleryDB.Delete(id)
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) Update(gallery *Gallery) error {
	fmt.Println("from gg update", gallery)
	return gg.db.Save(gallery).Error
}

func (gg *galleryGorm) Delete(id uint) error {
	g := Gallery{Model: gorm.Model{ID: id}}
	return gg.db.Delete(&g).Error
}

func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var g Gallery
	err := gg.db.Where("id = ?", id).First(&g).Error
	// fmt.Println("gallerygorm by ID", g, err)
	switch err {
	case nil:
		return &g, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (gg *galleryGorm) ByUserID(UserID uint) ([]Gallery, error) {
	var galleries []Gallery
	gg.db.Where("user_id = ?", UserID).Find(&galleries)
	return galleries, nil
}

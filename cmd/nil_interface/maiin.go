package main

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

// type Animal interface {
// 	speak()
// }

// func jiao(a Animal) {
// 	a.speak()
// }

// type Cat struct {
// }

// func (c Cat) speak() {
// 	fmt.Println("cat")
// }

// type Dog struct {
// }

// func (d Dog) speak() {
// 	fmt.Println("Dog")
// }

// func main() {
// 	c := Cat{}
// 	d := Dog{}
// 	jiao(c)
// 	jiao(d)
// }

type Animal struct {
	speaker
}

type speaker interface {
	speak()
}

type Cat struct {
}

type Fenail struct {
	speaker
}

func (c Cat) speak() {
	fmt.Println("cat")
}

type userGorm struct {
	db *gorm.DB
	// we don't need hmac here anymore once we move it to validator
	//hmac hash.HMAC
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"unique_index; non null"`
	Password     string `gorm:"-"` // ignore this in DB
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"unique_index; non null"`
}

// write a method create userGorm
func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	db.LogMode(true)
	if err != nil {
		return nil, err
	}
	// defer db.Close()
	// it's not good to leave this close in this function, but rather it
	// should be in a seperate func to close it

	//hmac not required anymore in userGorm
	//hmac := hash.NewHMAC(hmacSecretKey)
	return &userGorm{
		db: db,
	}, nil
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var u User
	err := ug.db.Where("id = ?", id).First(&u).Error
	switch err {
	case nil:
		return &u, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

var ErrNotFound = errors.New("models: resource not found")

func bar() {
	fmt.Println("foo")
}

func main() {
	d := Animal{&Fenail{Cat{}}}
	d.speak()
	bar()

}

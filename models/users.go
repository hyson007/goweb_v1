package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	// general error when resource is not found in db
	ErrNotFound  = errors.New("models: resource not found")
	ErrInvalidID = errors.New("models: Invalid User ID")
)

const userPWPepper = "secret-random-string"

func NewUserService(connectionInfo string) (*UserService, error) {

	db, err := gorm.Open("postgres", connectionInfo)
	db.LogMode(true)
	if err != nil {
		return nil, err
	}
	// defer db.Close()
	// it's not good to leave this close in this function, but rather it should be in a seperate func to close it
	return &UserService{db: db}, nil
}

type UserService struct {
	db *gorm.DB
}

//ByID will lookup user by the ID provided
//case 1 found user, return:  users, nil
//case 2 not found user, return: nil, ErrNotFound
//case 3 something wrong with db, nil, otherError, need to return 500

func (us *UserService) ByID(id uint) (*User, error) {
	var u User
	err := us.db.Where("id = ?", id).First(&u).Error
	switch err {
	case nil:
		return &u, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (us *UserService) ByEmail(email string) (*User, error) {
	var u User
	err := us.db.Where("email = ?", email).First(&u).Error
	switch err {
	case nil:
		return &u, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// create user
func (us *UserService) Create(u *User) error {
	//adding pepper into user pwd
	pwByte := []byte(u.Password + userPWPepper)

	//bcrypto takes a byte of string
	hashByte, err := bcrypt.GenerateFromPassword(pwByte, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashByte)
	// not required but good exercies
	u.Password = ""
	return us.db.Create(u).Error
}

// update user
func (us *UserService) Update(u *User) error {
	return us.db.Save(u).Error
}

// delete user
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// close user Service DB connection, note how this is a function for UserService struct, rather than a new function.

func (us *UserService) Close() error {
	return us.db.Close()
}

// Drop table and then auto migrate
func (us *UserService) ResetDB() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"unique_index; non null"`
	Password     string `gorm:"-"` // ignore this in DB
	PasswordHash string `gorm:"not null"`
}

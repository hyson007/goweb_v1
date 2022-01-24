package models

import (
	"errors"
	"goweb_v1/hash"
	"goweb_v1/rand"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	// general error when resource is not found in db
	ErrNotFound  = errors.New("models: resource not found")
	ErrInvalidID = errors.New("models: Invalid User ID")
	//ErrInvalidEmail = errors.New("models: Invalid User Email provided")
	ErrInvalidPwd = errors.New("models: Invalid User Password provided")
)

const userPWPepper = "secret-random-string"
const hmacSecretKey = "whatever"

//UserDB is used to interact with the user DB
//case 1 found user, return:  users, nil
//case 2 not found user, return: nil, ErrNotFound
//case 3 something wrong with db, nil, otherError, need to return 500
type UserDB interface {
	// Method for query single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	//Method for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	//Used to close DB connection
	Close() error

	//Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

func NewUserService(connectionInfo string) (*UserService, error) {
	ug, err := NewUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	// we switch to usergorm instead
	// db, err := gorm.Open("postgres", connectionInfo)
	// db.LogMode(true)
	// if err != nil {
	// 	return nil, err
	// }
	// // defer db.Close()
	// // it's not good to leave this close in this function, but rather it should be in a seperate func to close it
	// hmac := hash.NewHMAC(hmacSecretKey)
	return &UserService{
		UserDB: &userValidator{
			UserDB: ug,
		},
	}, nil
}

type UserService struct {
	UserDB
}

type userValidator struct {
	UserDB
}

// this line is a checker whether userGorm type matches with UserDB
// or if userGorm implements the correct UserDB all interfaces or not
var _ UserDB = &userGorm{}

// write a method create userGorm
func NewUserGorm(connectionInfo string) (*userGorm, error) {

	db, err := gorm.Open("postgres", connectionInfo)
	db.LogMode(true)
	if err != nil {
		return nil, err
	}
	// defer db.Close()
	// it's not good to leave this close in this function, but rather it should be in a seperate func to close it
	hmac := hash.NewHMAC(hmacSecretKey)
	return &userGorm{db: db, hmac: hmac}, nil
}

type userGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
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

//ByID will lookup user by the ID provided
//case 1 found user, return:  users, nil
//case 2 not found user, return: nil, ErrNotFound
//case 3 something wrong with db, nil, otherError, need to return 500

// func (us *UserService) ByID(id uint) (*User, error) {
// 	var u User
// 	err := us.db.Where("id = ?", id).First(&u).Error
// 	switch err {
// 	case nil:
// 		return &u, nil
// 	case gorm.ErrRecordNotFound:
// 		return nil, ErrNotFound
// 	default:
// 		return nil, err
// 	}
// }

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

// func (us *UserService) ByEmail(email string) (*User, error) {
// 	var u User
// 	err := us.db.Where("email = ?", email).First(&u).Error
// 	switch err {
// 	case nil:
// 		return &u, nil
// 	case gorm.ErrRecordNotFound:
// 		return nil, ErrNotFound
// 	default:
// 		return nil, err
// 	}
// }

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var u User
	err := ug.db.Where("email = ?", email).First(&u).Error
	switch err {
	case nil:
		return &u, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// By Remember looks up a user with the given token and returns that user,
// the method will handle the hashing the token for us
// func (us *UserService) ByRemember(token string) (*User, error) {
// 	hashedToken := us.hmac.Hash(token)
// 	var u User

// gorm by default use snake case for column name
// https://gorm.io/docs/conventions.html

// 	Column db name uses the field’s name’s snake_case by convention.

// type User struct {
//   ID        uint      // column name is `id`
//   Name      string    // column name is `name`
//   Birthday  time.Time // column name is `birthday`
//   CreatedAt time.Time // column name is `created_at`
// }
// You can override the column name with tag column or use NamingStrategy

// type Animal struct {
//   AnimalID int64     `gorm:"column:beast_id"`         // set name to `beast_id`
//   Birthday time.Time `gorm:"column:day_of_the_beast"` // set name to `day_of_the_beast`
//   Age      int64     `gorm:"column:age_of_the_beast"` // set name to `age_of_the_beast`
// }

// 	err := us.db.Where("remember_hash = ?", hashedToken).First(&u).Error
// 	switch err {
// 	case nil:
// 		return &u, nil
// 	case gorm.ErrRecordNotFound:
// 		return nil, ErrNotFound
// 	default:
// 		return nil, err
// 	}
// }

func (ug *userGorm) ByRemember(token string) (*User, error) {
	hashedToken := ug.hmac.Hash(token)
	var u User
	err := ug.db.Where("remember_hash = ?", hashedToken).First(&u).Error
	switch err {
	case nil:
		return &u, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Authenticate can be used to authenticate a user with the provided email
// address and password
// we will leave this func under User Service not User Gorm as it's not DB write
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	// fmt.Println(password)
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPWPepper))

	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPwd
		default:
			// below err is from bcrypt
			return nil, err
		}
	}

	return foundUser, nil
}

// create user
// func (us *UserService) Create(u *User) error {
// 	//adding pepper into user pwd
// 	pwByte := []byte(u.Password + userPWPepper)

// 	//bcrypto takes a byte of string
// 	hashByte, err := bcrypt.GenerateFromPassword(pwByte, bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}
// 	u.PasswordHash = string(hashByte)
// 	// not required but good exercies
// 	u.Password = ""

// 	// have a check for user, if no remember we set one.
// 	// and we should have a cookie whenever user create new account

// 	if u.Remember == "" {
// 		token, err := rand.RememberToken()
// 		if err != nil {
// 			return err
// 		}
// 		u.Remember = token
// 	}

// 	//now all users should have remember cookie then we create a hash for it
// 	u.RememberHash = us.hmac.Hash(u.Remember)
// 	return us.db.Create(u).Error
// }

func (ug *userGorm) Create(u *User) error {
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

	// have a check for user, if no remember we set one.
	// and we should have a cookie whenever user create new account

	if u.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		u.Remember = token
	}

	//now all users should have remember cookie then we create a hash for it
	u.RememberHash = ug.hmac.Hash(u.Remember)
	return ug.db.Create(u).Error
}

// update user
// func (us *UserService) Update(u *User) error {
// 	if u.Remember != "" {
// 		u.RememberHash = us.hmac.Hash(u.Remember)
// 	}

// 	return us.db.Save(u).Error
// }

func (ug *userGorm) Update(u *User) error {
	if u.Remember != "" {
		u.RememberHash = ug.hmac.Hash(u.Remember)
	}

	return ug.db.Save(u).Error
}

// delete user
// func (us *UserService) Delete(id uint) error {
// 	if id == 0 {
// 		return ErrInvalidID
// 	}
// 	user := User{Model: gorm.Model{ID: id}}
// 	return us.db.Delete(&user).Error
// }
func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// close user Service DB connection, note how this is a function for UserService struct, rather than a new function.

// func (us *UserService) Close() error {
// 	return us.db.Close()
// }
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// Drop table and then auto migrate
// func (us *UserService) ResetDB() {
// 	us.db.DropTableIfExists(&User{})
// 	us.db.AutoMigrate(&User{})
// }
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

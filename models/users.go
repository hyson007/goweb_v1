package models

import (
	"errors"
	"goweb_v1/hash"
	"goweb_v1/rand"
	"regexp"
	"strings"

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

//jon		@	calhoun		.	com
// we can move this to uservalidator
// var (
// 	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9._%+\-]+\.[a-z]{2,16}$`)
// )

var (
	ErrEmailRequired = errors.New("email address is required")
	ErrEmailInvalid  = errors.New("email address is not valid")
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

// these lines are checker whether a particular struct type matches with interface
var _ UserDB = &userValidator{}
var _ UserDB = &userGorm{}
var _ UserService = &userService{}

// UserService is a set of methods used to manipulate and work
// with the user model
type UserService interface {
	// Authenticate will verify the proovided email/pwd,
	// if they are correct, they user corresponding to that
	// email will be returned, otherwise you will receive either
	// ErrNotFound, ErrInvalidPassword or another error if thing
	// goes wrong
	Authenticate(email, password string) (*User, error)
	UserDB
}

// User represent the user model we stored in DB
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"unique_index; non null"`
	Password     string `gorm:"-"` // ignore this in DB
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"unique_index; non null"`
}

// Authenticate can be used to authenticate a user with the provided email
// address and password
// we will leave this func under User Service not User Gorm as it's not DB write
func (us *userService) Authenticate(email, password string) (*User, error) {
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

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)

	uv := newUserValidator(ug, hmac)
	// uv := &userValidator{
	// 	UserDB: ug,
	// 	hmac:   hmac,
	// }
	// we switch to usergorm instead
	// db, err := gorm.Open("postgres", connectionInfo)
	// db.LogMode(true)
	// if err != nil {
	// 	return nil, err
	// }
	// // defer db.Close()
	// // it's not good to leave this close in this function, but rather it should be in a seperate func to close it
	// hmac := hash.NewHMAC(hmacSecretKey)

	//when we return interface, we are essentially return a struct
	// which has implemented what method that interface demand!!!
	return &userService{
		UserDB: uv,
	}, nil
}

// this name can be anything
type userService struct {
	UserDB
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}

	}
	return nil
}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9._%+\-]+\.[a-z]{2,16}$`),
	}
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

//ByRemember will hash the remember token and then call the ByRemember on the
//subsequent UserDB layer
func (uv *userValidator) ByRemember(token string) (*User, error) {

	u := User{
		Remember: token,
	}

	if err := runUserValFuncs(&u,
		uv.hmacRemember); err != nil {
		return nil, err
	}

	//rememberHash := uv.hmac.Hash(token)

	// way 1
	// user, err := uv.UserDB.ByRemember(rememberHash)
	// if err != nil {
	// 	return nil, err
	// }
	// return user, nil

	// way 2
	return uv.UserDB.ByRemember(u.RememberHash)

}

// ByEmail will normalise address before calling By Email on the UserDB field
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user,
		uv.normalizeEmail,
	); err != nil {
		return nil, err
	}
	// this is not email but user.Email !!!
	return uv.UserDB.ByEmail(user.Email)
}

// Create will create the provided user and backfill data like ID, createdAt
// and UpdatedAt fields
// moved major to validator section
func (uv *userValidator) Create(u *User) error {
	// //adding pepper into user pwd
	// pwByte := []byte(u.Password + userPWPepper)

	// //bcrypto takes a byte of string
	// hashByte, err := bcrypt.GenerateFromPassword(pwByte, bcrypt.DefaultCost)
	// if err != nil {
	// 	return err
	// }
	// u.PasswordHash = string(hashByte)
	// // not required but good exercies
	// u.Password = ""

	// have a check for user, if no remember we set one.
	// and we should have a cookie whenever user create new account
	// if u.Remember == "" {
	// 	token, err := rand.RememberToken()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	u.Remember = token
	// }

	if err := runUserValFuncs(u,
		uv.pwdMinlen,
		uv.pwdRequired,
		uv.bcryptPassword,
		uv.pwdHashRequired,
		// set Rem should run before hmac
		uv.setRemIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail); err != nil {
		return err
	}

	//now all users should have remember cookie then we create a hash for it
	//u.RememberHash = uv.hmac.Hash(u.Remember)
	return uv.UserDB.Create(u)
}

// update wiill hash a remember token if it's provided.
func (uv *userValidator) Update(u *User) error {
	if err := runUserValFuncs(u,
		uv.pwdMinlen,
		uv.bcryptPassword,
		uv.pwdHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail); err != nil {
		return err
	}

	// if u.Remember != "" {
	// 	u.RememberHash = uv.hmac.Hash(u.Remember)
	// }

	return uv.UserDB.Update(u)
}

// delete user
func (uv *userValidator) Delete(id uint) error {
	user := User{
		Model: gorm.Model{
			ID: id,
		},
	}

	// alternate way
	//var user User
	//user.ID = id

	if err := runUserValFuncs(&user,
		uv.idGreaterThanzero); err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// if the user password is set, then the func will generate pwd hash with a
// predefined pepper (userPwPepper) and bcrypt
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwByte := []byte(user.Password + userPWPepper)
	hashByte, err := bcrypt.GenerateFromPassword(pwByte, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashByte)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) setRemIfUnset(user *User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	return nil
}

func (uv *userValidator) idGreaterThanzero(user *User) error {
	if user.ID <= 0 {
		return ErrInvalidID
	}
	return nil
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	// fmt.Println("from validator", user.Email)
	// // moving the emailregex under uv to avoid using global variable
	// if uv.emailRegex.MatchString(user.Email) {
	// 	return ErrEmailInvalid
	// }
	// return nil
	//fmt.Println("from validator", user.Email)
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	// we found a user with email address
	// if the found user has the same ID as this user, it;s an update
	// and this is the same user
	if user.ID != existing.ID {
		return errors.New("models: email address is already taken")
	}
	return nil
}

func (uv *userValidator) pwdMinlen(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return errors.New("password too short")
	}
	return nil
}

func (uv *userValidator) pwdRequired(user *User) error {
	if user.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (uv *userValidator) pwdHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return errors.New("password Hash is required")
	}
	return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return errors.New("models: remember token must be at least 32 bytes")
	}
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return errors.New("rememberhash is required")
	}
	return nil
}

type userGorm struct {
	db *gorm.DB
	// we don't need hmac here anymore once we move it to validator
	//hmac hash.HMAC
}

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
// the method expect the remember token already been hashed
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

//update, moved the hash part to validator
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var u User
	err := ug.db.Where("remember_hash = ?", rememberHash).First(&u).Error
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

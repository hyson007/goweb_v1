package main

import (
	"bufio"
	"fmt"
	"goweb_v1/models"
	"os"
	"strings"

	"github.com/jinzhu/gorm"

	//gorm's specific way to initiate
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "baloo"
	password = "junglebook"
	dbname   = "lenslocked"
)

type User struct {
	// this is embedded, User will have a Model key and some value under that
	// fmt.Println(user.Model.CreatedAt)
	// fmt.Println(user.CreatedAt)
	// these two print the same thing.

	// type Model struct {
	// 	ID        uint `gorm:"primary_key"`
	// 	CreatedAt time.Time
	// 	UpdatedAt time.Time
	// 	DeletedAt *time.Time `sql:"index"`
	gorm.Model
	Name   string
	Email  string `gorm:"unique_index; not null"`
	Orders []Order
}

type Order struct {
	gorm.Model
	UserID      uint
	Amount      int
	Description string
}

func main() {
	psqlinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := models.NewUserService(psqlinfo)
	if err != nil {
		panic(err)
	}

	// db.ResetDB()
	// usr, err := db.ByID(2)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(usr)

	// user := &models.User{
	// 	Name:  "jijo",
	// 	Email: "jijo@test.com",
	// }

	// if err := db.Create(user); err != nil {
	// 	panic(err)
	// }

	// user.Email = "jijo2@test.com"

	// if err := db.Update(user); err != nil {
	// 	panic(err)
	// }

	// emailUsr, err := db.ByEmail("jijo2@test.com")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(emailUsr)

	db.Delete(1)

	// db, err := gorm.Open("postgres", psqlinfo)
	// defer db.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// if err = db.DB().Ping(); err != nil {
	// 	panic(err)
	// }

	// showing SQL command runs in the backgroud, have to enable first...
	// db.LogMode(true)

	// dropping table
	//db.DropTableIfExists(&User{})

	// passing an empty object and it knows how to migrate
	// if struct has some new field, it will kick in auto migration.
	// db.AutoMigrate(&User{}, &Order{})

	//fmt.Println(db)

	// result
	// lenslocked=# select * from USERS;
	// id | created_at | updated_at | deleted_at | name | email
	// ----+------------+------------+------------+------+-------
	// (0 rows)

	// name, email := getInfoFromKeyboard()
	// u := User{
	// 	Name:  name,
	// 	Email: email,
	// }
	// db.Create(&u)
	// if err = db.Create(&u).Error; err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("+%v\n", u)

	// different ways to query
	// var u User
	// db.First(&u)
	// fmt.Println(u)

	// // different ways to query
	// db.Last(&u)
	// fmt.Println(u)

	// // different ways to query, it will assume 1 is the id
	// db.First(&u, 1)
	// fmt.Println(u)
	// // or use question mark
	// db.First(&u, "name = ?", "jackxie")
	// fmt.Println(u)

	// // different ways to query, it will assume 1 is the id
	// // gorm uses question mark
	// db.Where("name = ? AND email = ?", "jackxie", "hyson007@gmail.com").First(&u)
	// fmt.Println(u)

	// // multiple results
	// users := []User{}
	// db.Find(&users)
	// fmt.Println(len(users))
	// fmt.Println(users)

	// // error handling
	// if err := db.Where("email = ?", "b@b.com").First(&u).Error; err != nil {
	// 	switch err {
	// 	case gorm.ErrRecordNotFound:
	// 		fmt.Println("users not found")
	// 	default:
	// 		panic(err)
	// 	}
	// }

	//create orders for a few users
	// if err := db.First(&u).Error; err != nil {
	// 	panic(err)
	// }
	// // CreateOrders(db, u, 100, "description for orders1")
	// // CreateOrders(db, u, 888, "description for orders2")
	// // CreateOrders(db, u, 444, "description for orders3")

	// // relational query
	// if err := db.Preload("Orders").First(&u).Error; err != nil {
	// 	panic(err)
	// }
	// fmt.Println(u)
}

func getInfoFromKeyboard() (name, email string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("what's your name")

	// ReadString takes delim byte, so it has to be single quote
	name, _ = reader.ReadString('\n')

	fmt.Println("what's your Email address")
	email, _ = reader.ReadString('\n')

	return strings.TrimSpace(name), strings.TrimSpace(email)
}

func CreateOrders(db *gorm.DB, user User, amount int, desc string) {
	if err := db.Create(&Order{
		UserID:      user.ID,
		Amount:      amount,
		Description: desc,
	}).Error; err != nil {
		panic(err)
	}

}

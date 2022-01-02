package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "baloo"
	password = "junglebook"
	dbname   = "lenslocked"
)

func main() {
	psqlinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlinfo)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	// DB execute, Exec will not return anything (nil)
	// _, err = db.Exec(`
	// INSERT INTO USERS(name, email)
	// VALUES($1, $2)
	// `, "user1", "test@abc.com")
	// if err != nil {
	// 	panic(err)
	// }

	// Execute SQL insert and same time asking return id
	// var id int
	// err = db.QueryRow(`
	// INSERT INTO USERS(name, email)
	// VALUES($1, $2)
	// RETURNING id`,
	// 	"user2", "test2@abc.com").Scan(&id)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("id is ...", id)

	// Query a single result
	// var id int
	// var name, email string
	// row := db.QueryRow(`
	// SELECT id, name, email
	// FROM USERS
	// WHERE id=$1`, 4)

	// err = row.Scan(&id, &name, &email)

	// if err != nil {
	// 	// this will catch if no rows in the result set
	// 	if err == sql.ErrNoRows {
	// 		fmt.Println("No rows")
	// 	} else {
	// 		panic(err)
	// 	}
	// }

	// fmt.Println("query result, id is ...", id, name, email)

	// Query multiple result/rows
	// type Users struct {
	// 	ID    int
	// 	Name  string
	// 	Email string
	// }
	// var users []Users

	// rows, err := db.Query(`
	// SELECT id, name, email
	// FROM USERS`)

	// defer rows.Close()

	// for rows.Next() {
	// 	var user Users
	// 	if err = rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
	// 		panic(err)
	// 	}
	// 	users = append(users, user)
	// }

	// // should check error when row.next() finishes.
	// if rows.Err() != nil {
	// 	fmt.Println("some error during rows loop")
	// }
	// fmt.Println(users)

	// Query relational DB

	// CREATE TABLE ORDERS (
	// 	id serial PRIMARY KEY,
	// 	user_id int,
	// 	amount int,
	// 	description TEXT);

	// for i := 1; i <= 6; i++ {
	// 	userid := 1
	// 	amount := 100 * i
	// 	description := fmt.Sprintf("USB type C X %d", i)
	// 	if i > 3 {
	// 		userid = 4
	// 	}
	// 	_, err = db.Exec(`
	// 	INSERT INTO orders(user_id, amount, description)
	// 	VALUES($1, $2, $3)`, userid, amount, description)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// lenslocked=# select * from orders
	// lenslocked-# ;
	// id | user_id | amount |  description
	// ----+---------+--------+----------------
	// 1 |       1 |    100 | USB type C X 1
	// 2 |       1 |    200 | USB type C X 2
	// 3 |       1 |    300 | USB type C X 3
	// 4 |       4 |    400 | USB type C X 4
	// 5 |       4 |    500 | USB type C X 5
	// 6 |       4 |    600 | USB type C X 6
	// (6 rows)

	// lenslocked=# select * from users
	// lenslocked-# ;
	// id | name  |     email
	// ----+-------+---------------
	// 1 | user1 | test@abc.com
	// 3 | user2 | test2@abc.com
	// 4 | b     | bbb@bbb.com
	// (3 rows)

	// lenslocked=# SELECT * from USERS
	// INNER JOIN orders ON USERS.ID=ORDERS.user_id;
	// id | name  |    email     | id | user_id | amount |  description
	// ----+-------+--------------+----+---------+--------+----------------
	// 1 | user1 | test@abc.com |  1 |       1 |    100 | USB type C X 1
	// 1 | user1 | test@abc.com |  2 |       1 |    200 | USB type C X 2
	// 1 | user1 | test@abc.com |  3 |       1 |    300 | USB type C X 3
	// 4 | b     | bbb@bbb.com  |  4 |       4 |    400 | USB type C X 4
	// 4 | b     | bbb@bbb.com  |  5 |       4 |    500 | USB type C X 5
	// 4 | b     | bbb@bbb.com  |  6 |       4 |    600 | USB type C X 6
	// (6 rows)

	rows, err := db.Query(`
	SELECT * FROM USERS
	INNER JOIN orders ON USERS.ID=ORDERS.user_id`)

	defer rows.Close()

	for rows.Next() {
		var userID, orderID, amount, id int
		var userName, email, description string
		if err = rows.Scan(&userID, &userName, &email, &orderID, &id, &amount, &description); err != nil {
			panic(err)
		}
		fmt.Println("userID", userID, "userName", userName, "email", email, "orderID", orderID, "id", id, "amount", amount, "description", description)
	}

	// should check error when row.next() finishes.
	if rows.Err() != nil {
		panic(rows.Err())
	}

	// note above there are TWO ID columns, we can use AS to solve this.
	// lenslocked=# SELECT users.id, users.name, users.email, orders.id AS order_id, orders.amount, orders.description from USERS
	// INNER JOIN orders ON USERS.ID=ORDERS.user_id;
	// id | name  |    email     | order_id | amount |  description
	// ----+-------+--------------+----------+--------+----------------
	// 1 | user1 | test@abc.com |        1 |    100 | USB type C X 1
	// 1 | user1 | test@abc.com |        2 |    200 | USB type C X 2
	// 1 | user1 | test@abc.com |        3 |    300 | USB type C X 3
	// 4 | b     | bbb@bbb.com  |        4 |    400 | USB type C X 4
	// 4 | b     | bbb@bbb.com  |        5 |    500 | USB type C X 5
	// 4 | b     | bbb@bbb.com  |        6 |    600 | USB type C X 6

}

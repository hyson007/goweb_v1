package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Dog struct {
	ID     int
	Name   string
	Breed  string
	BornAt time.Time
}

type JSONDog struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Breed  string `json:"breed"`
	BornAt int64  `json:"born_at"`
}

func (d Dog) ConvertJsonDog() JSONDog {
	return JSONDog{
		d.ID,
		d.Name,
		d.Breed,
		d.BornAt.Unix(),
	}
}

func ToDog(j JSONDog) Dog {
	return Dog{
		j.ID,
		j.Name,
		j.Breed,
		time.Unix(j.BornAt, 0),
	}
}

func (d Dog) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.ConvertJsonDog())
}

// both methods are bind to Dog type.
func (d *Dog) UnmarshalJSON(data []byte) error {
	var jd JSONDog
	err := json.Unmarshal(data, &jd)
	if err != nil {
		panic(err)
	}
	*d = jd.ConvertDog()
	return nil
}

func (j JSONDog) ConvertDog() Dog {
	return Dog{
		j.ID,
		j.Name,
		j.Breed,
		time.Unix(j.BornAt, 0),
	}
}

func main() {
	var dog = Dog{1, "abo", "whatever", time.Now()}
	b, err := json.Marshal(dog)
	if err != nil {
		panic(err)
	}
	// fmt.Println("Original Dog Format")
	fmt.Println("using interface method , new Dog Format")
	fmt.Println(string(b))

	// fmt.Println("New Dog Format")
	// b, err = json.Marshal(ToJsonDog(dog))
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(b))

	by := []byte(`{
		"id":1,
		"name":"bowser",
		"breed":"husky",
		"born_at":1480979203}`)

	// var jd JSONDog
	var d Dog
	err = json.Unmarshal(by, &d)
	if err != nil {
		panic(err)
	}
	fmt.Println("unmarshal from api using interface method")
	fmt.Println(d)
	// fmt.Println("convert to our format")
	// fmt.Println(jd.ConvertDog())

}

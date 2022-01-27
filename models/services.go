package models

import "github.com/jinzhu/gorm"

func NewService(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)

	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		// this is the critical part, we have to initiated the User service here
		// otherwise we will get nil pointer error
		Gallery: NewGalleryservice(db),
		User:    NewUserService(db),
		db:      db,
	}, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

// close user Service DB connection, note how this is a function for UserService struct, rather than a new function.
// func (us *UserService) Close() error {
// 	return us.db.Close()
// }
func (s *Services) Close() error {
	return s.db.Close()
}

// Drop table and then auto migrate
// func (us *UserService) ResetDB() {
// 	us.db.DropTableIfExists(&User{})
// 	us.db.AutoMigrate(&User{})
// }
func (s *Services) DestructiveReset() error {
	if err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}

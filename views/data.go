package views

import (
	"goweb_v1/models"
	"log"
)

const (
	AlertLvError    = "danger"
	AlertLvWarning  = "warning"
	ALertLvInfo     = "info"
	AlertLvSuccess  = "success"
	AlertMsgGeneric = "something went wrong, pls try again and contact us if persists"
)

// ALert is used to render Boostrap Alert message in templates
type Alert struct {
	Level   string
	Message string
}

// Data is the top level struct that views expect data to come in
type Data struct {
	// setting to pointer can make it to nil
	// which can benefit the condition checking in bootstrap template
	Alert *Alert

	//by passing the user into Data struct will help to identify
	//if a given user is logged in or not.
	User  *models.User
	Yield interface{}
}

// this func takes any error and check if it implements Public or not
func (d *Data) SetAlert(err error) {
	// this is type assertion, if this err implement Public interface
	// if it does, then pErr will be the err casted into PublicError Type
	if pErr, ok := err.(PublicError); ok {
		//inside the if, we can call
		//pErr.Public()
		// log.Panicln("hit", pErr)
		d.Alert = &Alert{
			Level:   AlertLvError,
			Message: pErr.Public(),
		}
	} else {
		// means this is not public error
		log.Println(err)
		d.Alert = &Alert{
			Level:   AlertLvError,
			Message: AlertMsgGeneric,
		}
	}
}

// this helper func set alert message
func (d *Data) SetAlertText(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvError,
		Message: msg,
	}
}

// this interface will be used, if errors are allow public ones, then we can
// show to users, otherwise, we can only show a generic error
type PublicError interface {
	error
	Public() string
}

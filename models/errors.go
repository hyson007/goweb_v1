package models

import "strings"

const (
	// change the type of this to a new type which implements both error
	// and public interfaces

	ErrEmailRequired modelError = "models: email address is required"
	ErrEmailInvalid  modelError = "models: email address is not valid"
	// general error when resource is not found in db
	ErrNotFound        modelError = "models: resource not found"
	ErrInvalidID       modelError = "models: Invalid User ID"
	ErrInvalidPwd      modelError = "models: Invalid User Password provided"
	ErrEmailTaken      modelError = "models: email address is already taken"
	ErrPwdShort        modelError = "models: password too short"
	ErrPwdRequired     modelError = "models: password is required"
	ErrPwdHashRequired modelError = "models: password Hash is required"
	ErrRemToken        modelError = "models: remember token must be at least 32 bytes"
	ErrRemHashToken    modelError = "models: rememberhash is required"
	ErrUserIDRequired  modelError = "models: User ID is required"
	ErrTitleRequired   modelError = "models: Title is required"
)

type modelError string

// this implement the error interface
func (e modelError) Error() string {
	return string(e)
}

// this implement the public interface
// we just change a bit the prefix to show the difference
func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	return strings.Title(s)
}

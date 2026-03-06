package storage

import (
	"errors"
)

var (
	ErrorUserNotExist = errors.New("user is not exist")
)

type IStorage interface {
	// Create
	// Read
	// Update
	// Delete
}

type ICreator interface {
	// Create
}

//type IGeter interface {
// etc ...

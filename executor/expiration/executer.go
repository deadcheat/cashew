package expiration

import (
	"log"

	"github.com/deadcheat/cashew"
)

// Executor is a struct as receiver implements Executor interface
type Executor struct {
	eu cashew.ExpirationUseCase
}

func New(eu cashew.ExpirationUseCase) cashew.Executor {
	return &Executor{eu}
}

// Execute is implement for Executor interface and do batch process for expiration
func (e *Executor) Execute() {
	if err := e.eu.RevokeAll(); err != nil {
		log.Println(err)
	}
}

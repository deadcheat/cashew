package expiration

import (
	"log"

	"github.com/deadcheat/cashew"
)

// Executor is a struct as receiver implements Executor interface
type Executor struct{}

func New() cashew.Executor {
	return &Executor{}
}

// Execute is implement for Executor interface and do batch process for expiration
func (e *Executor) Execute() {
	log.Println("hoge")
}

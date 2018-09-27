package expiration

import (
	"log"

	"github.com/deadcheat/cashew"
)

// Executer is a struct as receiver implements Executer interface
type Executer struct{}

func New() cashew.Executer {
	return &Executer{}
}

// Execute is implement for Executer interface and do batch process for expiration
func (e *Executer) Execute() {
	log.Println("hoge")
}

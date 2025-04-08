package suser

import (
	"fmt"
	"math"
)

type SUser struct {
	ID int
	Firstname string
	Lastname string
	Age int
	Weight int
	Height int
	BodyFat float64
	IMC float64
	TargetWeight int
}
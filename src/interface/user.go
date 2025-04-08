package iuser

import (
	"fmt"
	"math"
)

type IUser interface {
	GetID() int
	GetFirstname() string
	GetLastname() string
	GetAge() int
	GetWeight() int
	GetHeight() int
	GetBodyFat() float64
	GetIMC() float64
	GetTargetWeight() int

	Build(id int, firstname string, lastname string, age int, weight int, height int) IUser

	SetTargetWeight(targetWeight int) IUser
	SetBodyFat(bodyFat float64) IUser
	SetIMC(imc float64) IUser
	SetWeight(weight int) IUser
	SetHeight(height int) IUser
	SetAge(age int) IUser
}

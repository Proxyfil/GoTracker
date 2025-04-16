package iuser

import (
	
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

	Create(id int, firstname string, lastname string, age int, weight int, height int) IUser

	SetTargetWeight(targetWeight int) IUser
	SetWeight(weight int) IUser
	SetHeight(height int) IUser
	SetAge(age int) IUser
}

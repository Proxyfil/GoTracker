package suser

import (
	iuser "gotracker/interface"
	"math"
)

type SUser struct {
	ID           int
	Firstname    string
	Lastname     string
	Age          int
	Weight       int
	Height       int
	TargetWeight int
}

// GetBodyFat returns the BodyFat value of the user
func (u *SUser) GetBodyFat() float64 {
	// We assume a simple formula for Body Fat calculation
	return float64(u.Weight) * 0.2
}

// GetIMC calculates and returns the IMC (Body Mass Index) of the user
func (u *SUser) GetIMC() float64 {
	if u.Height == 0 {
		return 0 // Avoid division by zero
	}
	return float64(u.Weight) / math.Pow(float64(u.Height)/100, 2)
}

// SetTargetWeight sets the target weight of the user
func (u *SUser) SetTargetWeight(targetWeight int) iuser.IUser {
	u.TargetWeight = targetWeight
	return u
}

// SetWeight sets the weight of the user
func (u *SUser) SetWeight(weight int) iuser.IUser {
	u.Weight = weight
	return u
}

// SetHeight sets the height of the user
func (u *SUser) SetHeight(height int) iuser.IUser {
	u.Height = height
	return u
}

// SetAge sets the age of the user
func (u *SUser) SetAge(age int) iuser.IUser {
	u.Age = age
	return u
}

func (u *SUser) GetID() int {
	return u.ID
}

func (u *SUser) GetFirstname() string {
	return u.Firstname
}

func (u *SUser) GetLastname() string {
	return u.Lastname
}

func (u *SUser) GetAge() int {
	return u.Age
}

func (u *SUser) GetWeight() int {
	return u.Weight
}

func (u *SUser) GetHeight() int {
	return u.Height
}

func (u *SUser) GetTargetWeight() int {
	return u.TargetWeight
}

func (u *SUser) Create(id int, firstname string, lastname string, age int, weight int, height int) iuser.IUser {
	u.ID = id
	u.Firstname = firstname
	u.Lastname = lastname
	u.Age = age
	u.Weight = weight
	u.Height = height
	return u
}

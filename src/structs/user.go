package suser

import (
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
	TargetWeight int
}

// GetBodyFat returns the BodyFat value of the user
func (u *SUser) GetBodyFat() float64 {
    return u.BodyFat
}

// GetIMC calculates and returns the IMC (Body Mass Index) of the user
func (u *SUser) GetIMC() float64 {
    if u.Height == 0 {
        return 0 // Avoid division by zero
    }
    return float64(u.Weight) / math.Pow(float64(u.Height)/100, 2)
}

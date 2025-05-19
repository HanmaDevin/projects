package types

import (
	"errors"

	"gorm.io/gorm"
)

type Object struct {
	gorm.Model

	Number1 *float64 `gorm:"not null" json:"number1"`
	Number2 *float64 `gorm:"not null" json:"number2"`
}

type Result struct {
	gorm.Model

	Res  float64 `gorm:"not null" json:"result"`
	Desc string  `gorm:"not null" json:"description"`
}

func (o *Object) Add() float64 {
	return *o.Number1 + *o.Number2
}

func (o *Object) Subtract() float64 {
	return *o.Number1 - *o.Number2
}

func (o *Object) Divide() (float64, error) {
	var divisor *float64 = o.Number2
	if *divisor == 0 {
		return -1, errors.New("division by 0 not possible")
	}

	return *o.Number1 / *divisor, nil
}

package types

import (
	"errors"
)

type Object struct {
	Number1 *float64 `json:"number1"`
	Number2 *float64 `json:"number2"`
}

type Result struct {
	Res  float64 `json:"result"`
	Desc string  `json:"description"`
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

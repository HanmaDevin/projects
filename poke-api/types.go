package main

type Pokemon struct {
	Name            string  `json:"name"`
	BaseExpierience int     `json:"base_experience"`
	Height          int     `json:"height"`
	Weight          int     `json:"weight"`
	Sprites         Sprites `json:"sprites"`
	Types           []Type  `json:"types"`
}

type Sprites struct {
	Other        OtherSprites `json:"other"`
	FrontDefault string       `json:"front_default"`
}

type OtherSprites struct {
	Showdown Showdown `json:"showdown"`
}

type Showdown struct {
	FrontDefault string `json:"front_default"`
}

type Type struct {
	Slot       int        `json:"slot"`
	TypeDetail TypeDetail `json:"type"`
}

type TypeDetail struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

package model

type Composition struct {
	User            User `json:"user"`
	Address         Address `json:"address"`
}
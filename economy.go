package main

import "gorm.io/gorm"

type EconomyUser struct {
	gorm.Model
	ID    string
	Coins int
}

func (e *EconomyUser) AddCoins(amount int) {
	e.Coins += amount
}
func (e *EconomyUser) RemoveCoins(amount int) {
	e.Coins -= amount
}

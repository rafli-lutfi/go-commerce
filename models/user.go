package models

import "time"

type User struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	Username      string        `json:"username" gorm:"size:100; not null; unique"`
	Email         string        `json:"email" gorm:"size:100; not null; unique"`
	Password      string        `json:"-" gorm:"size:256; not null"`
	FirstName     string        `json:"first_name" gorm:"size:256; not null"`
	LastName      string        `json:"last_name" gorm:"size:256; not null"`
	MobileNumber  string        `json:"mobile_number" gorm:"size:256; not null"`
	UserAddresses []UserAddress `json:"user_addresses"`
}

type UserAddress struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	UserID     uint   `json:"user_id" gorm:"not null"`
	Address    string `json:"address" gorm:"not null"`
	City       string `json:"city" gorm:"not null"`
	PostalCode string `json:"postal_code" gorm:"not null"`
	Country    string `json:"country" gorm:"not null"`
}

type NewUser struct {
	Username     string     `json:"username" binding:"required"`
	Password     string     `json:"password" binding:"required"`
	Email        string     `json:"email" binding:"required"`
	FirstName    string     `json:"first_name" binding:"required"`
	LastName     string     `json:"last_name" binding:"required"`
	MobileNumber string     `json:"mobile_number" binding:"required"`
	NewAddress   NewAddress `json:"address"`
}

type NewAddress struct {
	UserID     uint   `json:"user_id"`
	Address    string `json:"address" binding:"required"`
	City       string `json:"city" binding:"required"`
	PostalCode string `json:"postal_code" binding:"required"`
	Country    string `json:"country" binding:"required"`
}

type UserInfo struct {
	ID           uint          `json:"id"`
	Username     string        `json:"username"`
	Email        string        `json:"email"`
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	MobileNumber string        `json:"mobile_number"`
	Addresses    []UserAddress `json:"addresses"`
}

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

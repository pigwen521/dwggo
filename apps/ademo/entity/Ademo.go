package entity

import "gorm.io/gorm"

type Ademo struct {
	gorm.Model
	Name string
}

/*
默认ademos
func (Ademo) TableName() string {
	return "ademo"
} */

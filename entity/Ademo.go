package entity

import "gorm.io/gorm"

//XX实体
type Ademo struct {
	gorm.Model
	Name string
}

/*
重写表名，默认表名ademos
func (Ademo) TableName() string {
	return "ademo"
} */

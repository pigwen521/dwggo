package service

import (
	"dsjk.com/dwggo/core"
	"dsjk.com/dwggo/entity"
	"dsjk.com/dwggo/lib/helper/str"
	"dsjk.com/dwggo/model"
	"gorm.io/gorm"
)

type AdemoService struct{}

func (service *AdemoService) FindOne(name string) entity.Ademo {
	var obj entity.Ademo
	orm := core.Db
	orm.Where("name=?", name).Find(&obj)
	return obj
}
func (ser *AdemoService) Query(arg model.ArgAdemoQueryInModel) (*model.ArgAdemoQueryOutModel, error) {
	var ademo entity.Ademo
	t := ser.BuildCond(arg)
	err := t.Find(&ademo).Error
	if err != nil {
		return nil, err
	}

	arg_out := model.ArgAdemoQueryOutModel{}
	arg_out.ID = str.ToString(ademo.ID)
	arg_out.Name = ademo.Name
	return &arg_out, nil
}

func (ser *AdemoService) BuildCond(arg model.ArgAdemoQueryInModel) *gorm.DB {
	orm := core.Db
	t := orm.Where("id>0")
	if arg.Name != "" {
		t.Where("name =?", arg.Name)
	}
	return t
}

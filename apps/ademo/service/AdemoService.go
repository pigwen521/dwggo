package service

import (
	"errors"

	"dsjk.com/dwggo/apps/ademo/entity"
	"dsjk.com/dwggo/apps/ademo/model"
	"dsjk.com/dwggo/system/core"
	"dsjk.com/dwggo/system/lib/helper/str"
)

type AdemoService struct{}

func (ser *AdemoService) Query(arg model.ArgAdemoQueryInModel) (*model.ArgAdemoQueryOutModel, error) {
	var obj entity.Ademo
	orm := core.Db

	field_name := core.GetFieldAlias(obj, obj.Name)
	err := orm.Where(field_name+"=?", arg.Name).Find(&obj).Error
	if err != nil {
		core.LogError(err.Error())
		return nil, errors.New("系统错误,DB")
	}

	arg_out := model.ArgAdemoQueryOutModel{}
	arg_out.ID = str.ToString(obj.ID)
	arg_out.Name = obj.Name
	return &arg_out, nil
}

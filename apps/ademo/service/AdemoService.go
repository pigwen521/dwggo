package service

import (
	"dsjk.com/dwggo/apps/ademo/entity"
	"dsjk.com/dwggo/apps/ademo/model"
	"dsjk.com/dwggo/system/core"
	"dsjk.com/dwggo/system/lib/helper/str"
)

type AdemoService struct{}

func (ser *AdemoService) Query(arg model.ArgAdemoQueryInModel) (*model.ArgAdemoQueryOutModel, error) {
	var obj entity.Ademo
	orm := core.Db
	orm.Where("name=?", arg.Name).Find(&obj)

	arg_out := model.ArgAdemoQueryOutModel{}
	arg_out.ID = str.ToString(obj.ID)
	arg_out.Name = obj.Name
	return &arg_out, nil
}

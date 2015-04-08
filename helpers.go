package rest

import (
	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"reflect"
	"strconv"
)

func FindRecordFromIDParameter(c martini.Context, out interface{}) bool {
	r := c.Get(inject.InterfaceOf((*render.Render)(nil))).Interface().(render.Render)
	db := c.Get(reflect.TypeOf((*gorm.DB)(nil))).Interface().(*gorm.DB)
	var params martini.Params
	params = c.Get(reflect.TypeOf(params)).Interface().(martini.Params)

	id, e := strconv.Atoi(params["id"])
	if e != nil {
		RenderErrorNotAllowed(r)
		return false
	}
	db.First(out, id)
	foundId := reflect.ValueOf(out).Elem().FieldByName("Id").Int()
	if foundId == 0 {
		RenderErrorNotAllowed(r)
		return false
	}
	return true
}

func KeyExists(db *gorm.DB, model interface{}, id int64) bool {
	count := 0
	db.Model(model).Where("id = ?", id).Count(&count)
	return count > 0
}

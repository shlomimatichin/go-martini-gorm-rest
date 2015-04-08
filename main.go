package rest

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"github.com/shlomimatichin/go-martini-pagination"
	"io/ioutil"
	"net/http"
	"reflect"
)

func CreateRecordView(jsonRecord interface{}, construct interface{}) func(r render.Render, c martini.Context) {
	return func(r render.Render, c martini.Context) {
		returned, err := c.Invoke(construct)
		if err != nil {
			panic(fmt.Sprint("Unable to invoke Construct:", err))
		}
		id := returned[0].Int()
		if id == 0 {
			RenderErrorNotAllowed(r)
			return
		}
		r.JSON(200, map[string]interface{}{"id": id})
	}
}

func DeleteRecordView(viewParams Params) func(r render.Render, db *gorm.DB, c martini.Context) {
	viewParams.sanity()
	viewParams.sanityMustNotHaveToJSON()
	viewParams.sanityMustNotHaveFilter()
	viewParams.sanityMustNotHaveValid()
	viewParams.sanityMustNotHaveField()
	return func(r render.Render, db *gorm.DB, c martini.Context) {
		zero := viewParams.zero()
		if !FindRecordFromIDParameter(c, zero.Interface()) {
			return
		}
		c.Map(zero.Interface())
		if !viewParams.allowed(c) {
			return
		}
		db.Delete(zero.Interface())
		RenderResultOK(r)
	}
}

func ListRecordsView(viewParams Params) func(p *pagination.Pagination, db *gorm.DB, c martini.Context) {
	viewParams.sanity()
	viewParams.sanityMustHaveToJSON()
	viewParams.sanityMustNotHaveValid()
	viewParams.sanityMustNotHaveField()
	return func(p *pagination.Pagination, db *gorm.DB, c martini.Context) {
		records := viewParams.slice()
		viewParams.filter(c).Offset(int(p.Offset)).Limit(int(p.PerPage)).Find(records.Interface())
		total := records.Elem().Len()
		p.SetTotal(uint(total))
		for i := 0; i < total; i++ {
			record := records.Elem().Index(i)
			c.Map(record.Addr().Interface())
			p.Append(viewParams.toJSON(c))
		}
	}
}

func GetRecordView(viewParams Params) func(r render.Render, db *gorm.DB, c martini.Context) {
	viewParams.sanity()
	viewParams.sanityMustHaveToJSON()
	viewParams.sanityMustNotHaveFilter()
	viewParams.sanityMustNotHaveValid()
	viewParams.sanityMustNotHaveField()
	return func(r render.Render, db *gorm.DB, c martini.Context) {
		record := viewParams.zero()
		if !FindRecordFromIDParameter(c, record.Interface()) {
			return
		}
		c.Map(record.Interface())
		if !viewParams.allowed(c) {
			return
		}
		r.JSON(200, viewParams.toJSON(c))
	}
}

func ModifyRecordFieldView(viewParams Params) func(r render.Render, db *gorm.DB, req *http.Request, c martini.Context) {
	viewParams.sanity()
	viewParams.sanityMustNotHaveToJSON()
	viewParams.sanityMustNotHaveFilter()
	viewParams.sanityMustHaveField()
	return func(r render.Render, db *gorm.DB, req *http.Request, c martini.Context) {
		zero := viewParams.zero()
		if !FindRecordFromIDParameter(c, zero.Interface()) {
			return
		}
		c.Map(zero.Interface())
		if !viewParams.allowed(c) {
			return
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic("Unable to read body")
		}
		var value map[string]interface{}
		err2 := json.Unmarshal(body, &value)
		if err2 != nil {
			RenderError(r, 422, fmt.Sprint("Unable to parse JSON:", err2))
			return
		}
		if value["value"] == nil {
			RenderError(r, 422, fmt.Sprint("JSON must contain field 'value'"))
			return
		}
		field := zero.Elem().FieldByName(viewParams.Field)
		switch field.Interface().(type) {
		default:
			panic(fmt.Sprint("Unknown field type:", field.Type()))
		case float64:
			if reflect.TypeOf(value["value"]).Kind() != reflect.Float64 {
				RenderError(r, 422, fmt.Sprint("Value can not be parsed as float:", reflect.TypeOf(value["value"])))
				return
			}
			final := value["value"].(float64)
			c.Map(final)
			if !viewParams.valid(c) {
				return
			}
			field.SetFloat(final)
		case int64:
			if reflect.TypeOf(value["value"]).Kind() != reflect.Float64 {
				RenderError(r, 422, fmt.Sprint("Value can not be parsed as int:", reflect.TypeOf(value["value"])))
				return
			}
			final := int64(value["value"].(float64))
			c.Map(final)
			if !viewParams.valid(c) {
				return
			}
			field.SetInt(final)
		case string:
			if reflect.TypeOf(value["value"]).Kind() != reflect.String {
				RenderError(r, 422, fmt.Sprint("Value can not be parsed as string:", reflect.TypeOf(value["value"])))
				return
			}
			final := value["value"].(string)
			c.Map(final)
			if !viewParams.valid(c) {
				return
			}
			field.SetString(final)
		}
		db.Save(zero.Interface())
		r.JSON(200, map[string]interface{}{"value": value["value"]})
	}
}

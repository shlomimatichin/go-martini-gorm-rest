package rest

import (
	"fmt"
	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
	"reflect"
)

type Params struct {
	Model   interface{} // model struct zero value, e.g., User{}
	Allowed interface{} // dep injected func returning one of the not allowed reasons strings
	ToJSON  interface{} // dep injected func returning map[string]interface{}
	Filter  interface{} // dep injected func returning db *gorm.DB modified
	Valid   interface{} // dep injected func returning one of the not allowed reasons strings
	Field   string
}

func (self *Params) sanity() {
	if reflect.TypeOf(self.Model).Kind() != reflect.Struct {
		panic(fmt.Sprintf("Model is not a struct"))
	}
	if self.Allowed != nil {
		if reflect.TypeOf(self.Allowed).Kind() != reflect.Func {
			panic(fmt.Sprintf("Allowed is not a func"))
		}
	}
	if self.ToJSON != nil {
		if reflect.TypeOf(self.ToJSON).Kind() != reflect.Func {
			panic(fmt.Sprintf("ToJSON is not a func"))
		}
	}
	if self.Filter != nil {
		if reflect.TypeOf(self.Filter).Kind() != reflect.Func {
			panic(fmt.Sprintf("Filter is not a func"))
		}
	}
	if self.Valid != nil {
		if reflect.TypeOf(self.Valid).Kind() != reflect.Func {
			panic(fmt.Sprintf("Valid is not a func"))
		}
	}
}

func (self *Params) sanityMustNotHaveToJSON() {
	if self.ToJSON != nil {
		panic("ToJSON must not be set")
	}
}

func (self *Params) sanityMustHaveToJSON() {
	if self.ToJSON == nil {
		panic("ToJSON must be set")
	}
}

func (self *Params) sanityMustNotHaveFilter() {
	if self.Filter != nil {
		panic("Filter must not be set")
	}
}

func (self *Params) sanityMustNotHaveValid() {
	if self.Valid != nil {
		panic("Valid must not be set")
	}
}

func (self *Params) sanityMustNotHaveField() {
	if self.Field != "" {
		panic("Field must not be set")
	}
}

func (self *Params) sanityMustHaveField() {
	if self.Field == "" {
		panic("Field must be set")
	}
}

func (self *Params) zero() reflect.Value {
	return reflect.New(reflect.TypeOf(self.Model))
}

func (self *Params) slice() reflect.Value {
	return reflect.New(reflect.SliceOf(reflect.TypeOf(self.Model)))
}

func (self *Params) allowed(c martini.Context) bool {
	if self.Allowed == nil {
		return true
	}
	returned, err := c.Invoke(self.Allowed)
	if err != nil {
		panic(fmt.Sprintf("Unable to invoke Allowed: %s", err))
	}
	reason := returned[0].String()
	if reason != Allowed {
		r := c.Get(inject.InterfaceOf((*render.Render)(nil))).Interface().(render.Render)
		RenderError(r, 405, reason)
		return false
	}
	return true
}

func (self *Params) toJSON(c martini.Context) map[string]interface{} {
	var db *gorm.DB
	db = c.Get(reflect.TypeOf(db)).Interface().(*gorm.DB)
	returned, err := c.Invoke(self.ToJSON)
	if err != nil {
		panic(fmt.Sprintf("Unable to invoke ToJSON: %s", err))
	}
	return returned[0].Interface().(map[string]interface{})
}

func (self *Params) filter(c martini.Context) *gorm.DB {
	var db *gorm.DB
	db = c.Get(reflect.TypeOf(db)).Interface().(*gorm.DB)
	if self.Filter == nil {
		return db
	}
	returned, err := c.Invoke(self.Filter)
	if err != nil {
		panic(fmt.Sprintf("Unable to invoke Filter: %s", err))
	}
	return returned[0].Interface().(*gorm.DB)
}

func (self *Params) valid(c martini.Context) bool {
	if self.Valid == nil {
		return true
	}
	returned, err := c.Invoke(self.Valid)
	if err != nil {
		panic(fmt.Sprintf("Unable to invoke Valid: %s", err))
	}
	reason := returned[0].String()
	if reason != Allowed {
		r := c.Get(inject.InterfaceOf((*render.Render)(nil))).Interface().(render.Render)
		RenderError(r, 405, reason)
		return false
	}
	return true
}

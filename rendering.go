package rest

import (
	"github.com/martini-contrib/render"
)

func RenderError(r render.Render, errorCode int, message string) {
	r.JSON(errorCode, map[string]interface{}{"Error": message})
}

func RenderErrorNotAllowed(r render.Render) {
	RenderError(r, 405, NotAllowed)
}

func RenderResultOK(r render.Render) {
	r.JSON(200, map[string]interface{}{"Result": "OK"})
}

package rest

type StringValueForm struct {
	Value string `form:"value" binding:"required"`
}

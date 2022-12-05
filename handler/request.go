package handler

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type createUpdateListRequest struct {
	Title       string                `json:"title" form:"title" validate:"required"`
	Description string                `json:"description" form:"description" validate:"required"`
	Level0ID    int64                 `json:"level0-id" form:"level0-id"`
	File        *multipart.FileHeader `json:"-" form:"-"`
}

func (r *createUpdateListRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	file, err := c.FormFile("file")
	if err == nil {
		r.File = file
	}

	return nil
}

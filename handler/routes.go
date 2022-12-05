package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
	v1.GET("/ping", h.Ping)

	list := v1.Group("/list")
	list.GET("/root", h.ListRoots)
	list.GET("", h.Lists)
	list.GET("/:id/file", h.GetListFile)
	list.POST("", h.CreateList)
	list.PUT("/:id", h.UpdateList)
	list.DELETE("/:id", h.DeleteList)
}

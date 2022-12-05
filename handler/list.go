package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kevin-luvian/moonlay-test/model"
	"github.com/kevin-luvian/moonlay-test/pkg/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Lists(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 0 {
		page = 0
	}

	perPage, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil {
		perPage = -1
	}

	lists, total, err := h.listUC.GetAll(page, perPage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"lists":    lists,
		"page":     page,
		"per_page": perPage,
		"total":    total,
	})
}

func (h *Handler) ListRoots(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 0 {
		page = 0
	}

	perPage, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil {
		perPage = -1
	}

	show := strings.ToLower(c.QueryParam("show_sublist")) == "true"

	lists, total, err := h.listUC.GetRoots(page, perPage, show)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"lists":    lists,
		"page":     page,
		"per_page": perPage,
		"total":    total,
	})
}

func (h *Handler) CreateList(c echo.Context) error {
	var (
		err error
	)

	req := &createUpdateListRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	list := model.List{
		Title:       req.Title,
		Description: req.Description,
		Level0ID:    sql.NullInt64{Int64: req.Level0ID, Valid: true},
	}

	list, err = h.listUC.Create(list)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	fmt.Println("Created ID", list.ID)

	if req.File != nil {
		list.FileUpload, err = h.listUC.UpdateFile(req.File, list)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}

		list, err = h.listUC.Update(list)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}

		fmt.Println("Updated ID", list.ID)
	}

	return c.JSON(http.StatusOK, list)
}

func (h *Handler) UpdateList(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorStr("invalid id type"))
	}

	req := &createUpdateListRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	list := model.List{
		ID:          uint(id),
		Title:       req.Title,
		Description: req.Description,
		Level0ID:    sql.NullInt64{Int64: req.Level0ID, Valid: true},
	}

	if req.File != nil {
		list.FileUpload, err = h.listUC.UpdateFile(req.File, list)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
	}

	list, err = h.listUC.Update(list)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, list)
}

func (h *Handler) DeleteList(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorStr("invalid id type"))
	}

	err = h.listUC.DeleteByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func (h *Handler) GetListFile(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorStr("invalid id type"))
	}

	list, err := h.listUC.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	reader, ctype, err := h.listUC.GetFile(list)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.Stream(http.StatusOK, ctype, reader)
}

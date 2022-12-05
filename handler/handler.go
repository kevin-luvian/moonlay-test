package handler

import (
	"github.com/kevin-luvian/moonlay-test/usecase"
)

type Handler struct {
	listUC *usecase.ListUC
}

func NewHandler(lu *usecase.ListUC) *Handler {
	return &Handler{
		listUC: lu,
	}
}

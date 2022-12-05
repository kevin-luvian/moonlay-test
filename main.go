package main

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/kevin-luvian/moonlay-test/handler"
	"github.com/kevin-luvian/moonlay-test/pkg/db"
	"github.com/kevin-luvian/moonlay-test/pkg/settings"
	"github.com/kevin-luvian/moonlay-test/repo"
	"github.com/kevin-luvian/moonlay-test/router"
	"github.com/kevin-luvian/moonlay-test/usecase"
)

func main() {
	settings.Setup()
	log.Println("db settings", settings.DatabaseSetting)

	d := db.New()
	db.AutoMigrate(d)

	lr := repo.NewListRepo(d)

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b) + settings.AppSetting.BaseFilePath
	fr := repo.NewFileRepo(basepath)

	luc := usecase.NewListUC(lr, fr)

	h := handler.NewHandler(luc)

	r := router.New()

	v1 := r.Group("/api")
	h.Register(v1)

	r.Logger.Fatal(r.Start(settings.AppSetting.Address))
}

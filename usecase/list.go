package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/jinzhu/gorm"
	"github.com/kevin-luvian/moonlay-test/model"
)

type iListRepo interface {
	GetByID(uint) (model.List, error)
	GetRoots(offset, limit int, preload bool) ([]model.List, int, error)
	GetAll(offset, limit int) ([]model.List, int, error)
	Create(l model.List) (model.List, error)
	Update(l *model.List) (err error)
	DeleteByID(id uint) error
}

type iFileRepo interface {
	MatchExt(filename string, extensions []string) bool
	GetFileContentType(fileID string) string
	SaveFile(file *multipart.FileHeader) (string, error)
	StreamFile(fileID string) (io.Reader, error)
	DeleteFile(fileID string) error
}

type ListUC struct {
	lr iListRepo
	fr iFileRepo
}

func NewListUC(lr iListRepo, fr iFileRepo) *ListUC {
	return &ListUC{
		lr: lr,
		fr: fr,
	}
}

func (u *ListUC) GetRoots(page, perPage int, showSublist bool) ([]model.List, int, error) {
	offset := page * perPage

	lists, count, err := u.lr.GetRoots(offset, perPage, showSublist)
	return lists, count, err
}

func (u *ListUC) GetAll(page, perPage int) ([]model.List, int, error) {
	offset := page * perPage

	lists, count, err := u.lr.GetAll(offset, perPage)
	return lists, count, err
}

func (u *ListUC) GetByID(id uint) (model.List, error) {
	return u.lr.GetByID(id)
}

func (u *ListUC) checkLevel0(l *model.List) error {
	level0ID := uint(l.Level0ID.Int64)

	if level0ID == 0 {
		l.Level0ID = sql.NullInt64{Valid: false}
		return nil
	}

	if level0ID == l.ID {
		return errors.New("not a valid level-0 id")
	}

	// check level-0
	lvl0, err := u.lr.GetByID(level0ID)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errors.New("not a valid level-0 id")
		}
		return err
	}

	if lvl0.Level0ID.Valid {
		return errors.New("can't attach to level-1")
	}

	return nil
}

func (u *ListUC) Create(l model.List) (model.List, error) {
	if err := u.checkLevel0(&l); err != nil {
		return model.List{}, err
	}

	return u.lr.Create(l)
}

func (u *ListUC) Update(l model.List) (model.List, error) {
	err := u.checkLevel0(&l)
	if err != nil {
		return model.List{}, err
	}
	fmt.Println("list.FileUpload", l.FileUpload)

	err = u.lr.Update(&l)
	if err != nil {
		return model.List{}, err
	}

	returned, err := u.lr.GetByID(l.ID)
	return returned, err
}

func (u *ListUC) DeleteByID(id uint) error {
	l, err := u.lr.GetByID(id)
	if err != nil {
		return err
	}

	if l.FileUpload != "" {
		u.fr.DeleteFile(l.FileUpload)
	}

	return u.lr.DeleteByID(id)
}

func (u *ListUC) UpdateFile(file *multipart.FileHeader, l model.List) (string, error) {
	if file == nil {
		return "", errors.New("invalid file")
	}

	file.Filename = fmt.Sprintf("%d-%s", l.ID, file.Filename)

	isMatch := u.fr.MatchExt(file.Filename, []string{".txt", ".pdf"})
	if !isMatch {
		return "", errors.New("invalid file extension")
	}

	fileID, err := u.fr.SaveFile(file)
	if err != nil {
		return "", err
	}

	if l.FileUpload != "" && l.FileUpload != fileID {
		// continue on error
		u.fr.DeleteFile(l.FileUpload)
	}

	return fileID, nil
}

func (u *ListUC) GetFile(l model.List) (io.Reader, string, error) {
	if l.FileUpload == "" {
		return nil, "", errors.New("invalid file")
	}

	contentType := u.fr.GetFileContentType(l.FileUpload)
	file, err := u.fr.StreamFile(l.FileUpload)

	return file, contentType, err
}

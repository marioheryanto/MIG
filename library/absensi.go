package library

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/marioheryanto/MIG/helper"
	"github.com/marioheryanto/MIG/model"
	"github.com/marioheryanto/MIG/repository"
)

type AbsensiLibrary struct {
	Repo repository.AbsensiRepositoryInterface
}

type AbsensiLibraryInterface interface {
	CheckInOut(tokenString string, check string) error
	GetAll(tokenString string, month, year int) (interface{}, error)
}

func NewAbsensiLibrary(repo repository.AbsensiRepositoryInterface) AbsensiLibraryInterface {
	return AbsensiLibrary{
		Repo: repo,
	}
}

func (l AbsensiLibrary) CheckInOut(tokenString string, check string) error {
	var err error
	now := time.Now().In(helper.LoadLocationJakarta()).Format("2006-01-02")

	claims, err := helper.ParseTokenToClaims(tokenString)
	if err != nil {
		return err
	}

	absensi, err := l.Repo.GetAbsensi(now, claims.Issuer)
	if serviceErr, _ := err.(model.ServiceError); serviceErr.Code == http.StatusInternalServerError {
		return err
	}

	if check == "out" && absensi.CheckIn == "" {
		return helper.NewServiceError(http.StatusBadRequest, "must check in first")
	}

	if (absensi.CheckIn != "" && check == "in") || (absensi.CheckOut != "" && check == "out") {
		return helper.NewServiceError(http.StatusBadRequest, fmt.Sprintf("already check %v", check))
	}

	if check == "in" {
		err := l.Repo.CreateAbsensi(check, claims.Issuer)
		if err != nil {
			return err
		}
	} else {
		err = l.Repo.CheckInOut(check, claims.Issuer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l AbsensiLibrary) GetAll(tokenString string, month, year int) (interface{}, error) {
	var err error

	claims, err := helper.ParseTokenToClaims(tokenString)
	if err != nil {
		return nil, err
	}

	from, to := helper.GenerateRange(month, year)

	absensi, err := l.Repo.GetRangeAbsensi(from, to, claims.Issuer)
	if err != nil {
		return nil, err
	}

	if len(absensi) == 0 {
		return nil, errors.New("no data found")
	}

	return absensi, nil
}

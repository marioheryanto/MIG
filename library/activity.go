package library

import (
	"errors"
	"net/http"
	"time"

	"github.com/marioheryanto/MIG/helper"
	"github.com/marioheryanto/MIG/model"
	"github.com/marioheryanto/MIG/repository"
)

type ActivityLibrary struct {
	Repo        repository.ActivityRepositoryInterface
	AbsensiRepo repository.AbsensiRepositoryInterface
}

type ActivityLibraryInterface interface {
	Create(tokenString string, request model.Activity) error
	Edit(tokenString string, request model.Activity) error
	Delete(tokenString string, activityId string) error
	GetAll(tokenString, from, to string) (interface{}, error)
}

func NewActivityLibrary(repo repository.ActivityRepositoryInterface, absensiRepo repository.AbsensiRepositoryInterface) ActivityLibraryInterface {
	return ActivityLibrary{
		Repo:        repo,
		AbsensiRepo: absensiRepo,
	}
}

func (l ActivityLibrary) Create(tokenString string, request model.Activity) error {
	timeNow := time.Now()
	request.CreatedAt = timeNow.Format("2006-01-02 15:04:05")

	claims, err := helper.ParseTokenToClaims(tokenString)
	if err != nil {
		return err
	}

	_, err = l.AbsensiRepo.GetAbsensi(request.Tanggal, claims.Issuer)
	if err != nil {
		serviceErr, _ := err.(model.ServiceError)
		if serviceErr.Code == http.StatusNotFound {
			return helper.NewServiceError(http.StatusBadRequest, "harus melakukan absensi terlebih dahulu")
		}

		return err
	}

	err = l.Repo.CreateActivity(request, claims.Issuer)
	if err != nil {
		return err
	}

	return nil
}

func (l ActivityLibrary) Edit(tokenString string, request model.Activity) error {
	claims, err := helper.ParseTokenToClaims(tokenString)
	if err != nil {
		return err
	}

	_, err = l.AbsensiRepo.GetAbsensi(request.Tanggal, claims.Issuer)
	if err != nil {
		serviceErr, _ := err.(model.ServiceError)
		if serviceErr.Code == http.StatusNotFound {
			return helper.NewServiceError(http.StatusBadRequest, "harus melakukan absensi terlebih dahulu")
		}

		return err
	}

	_, err = l.Repo.GetActivityWithId(request.Id, claims.Issuer)
	if err != nil {
		return err
	}

	timeNow := time.Now()
	request.UpdatedAt = timeNow.Format("2006-01-02 15:04:05")

	err = l.Repo.EditActivity(request)
	if err != nil {
		return err
	}

	return nil
}

func (l ActivityLibrary) Delete(tokenString string, activityId string) error {
	claims, err := helper.ParseTokenToClaims(tokenString)
	if err != nil {
		return err
	}

	data, err := l.Repo.GetActivityWithId(activityId, claims.Issuer)
	if err != nil {
		return err
	}

	_, err = l.AbsensiRepo.GetAbsensi(data.Tanggal, claims.Issuer)
	if err != nil {
		serviceErr, _ := err.(model.ServiceError)
		if serviceErr.Code == http.StatusNotFound {
			return helper.NewServiceError(http.StatusBadRequest, "harus melakukan absensi terlebih dahulu")
		}

		return err
	}

	return l.Repo.DeleteActivity(activityId)
}

func (l ActivityLibrary) GetAll(tokenString, from, to string) (interface{}, error) {
	var err error

	claims, err := helper.ParseTokenToClaims(tokenString)
	if err != nil {
		return nil, err
	}

	activities, err := l.Repo.GetRangeActivity(from, to, claims.Issuer)
	if err != nil {
		return nil, err
	}

	if len(activities) == 0 {
		return nil, errors.New("no data found")
	}

	return activities, nil
}

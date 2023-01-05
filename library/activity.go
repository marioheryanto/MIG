package library

import (
	"errors"
	"time"

	"github.com/marioheryanto/MIG/helper"
	"github.com/marioheryanto/MIG/model"
	"github.com/marioheryanto/MIG/repository"
)

type ActivityLibrary struct {
	Repo repository.ActivityRepositoryInterface
}

type ActivityLibraryInterface interface {
	Create(tokenString string, request model.Activity) error
	Edit(tokenString string, request model.Activity) error
	Delete(tokenString string, activityId string) error
	GetAll(tokenString, from, to string) (interface{}, error)
}

func NewActivityLibrary(repo repository.ActivityRepositoryInterface) ActivityLibraryInterface {
	return ActivityLibrary{
		Repo: repo,
	}
}

func (l ActivityLibrary) Create(tokenString string, request model.Activity) error {
	timeNow := time.Now()
	request.CreatedAt = timeNow.Format("2006-01-02 15:04:05")

	claims, err := helper.ParseTokenToClaims(tokenString)
	if err != nil {
		return err
	}

	err = l.Repo.CreateActivity(request, claims.Issuer)
	if err != nil {
		return err
	}

	return nil
}

func (l ActivityLibrary) Edit(tokenString string, request model.Activity) error {
	_, err := l.Repo.GetActivityWithId(request.Id)
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

	_, err := l.Repo.GetActivityWithId(activityId)
	if err != nil {
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

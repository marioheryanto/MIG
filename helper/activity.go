package helper

import (
	"net/http"

	"github.com/marioheryanto/MIG/model"
)

func ValidateActivity(request model.Activity, checkId bool) error {
	if request.Description == "" {
		return NewServiceError(http.StatusBadRequest, "Description is empty")
	}

	if request.Tanggal == "" {
		return NewServiceError(http.StatusBadRequest, "Tanggal is empty")
	}

	if request.Dari == "" {
		return NewServiceError(http.StatusBadRequest, "Dari is empty")
	}

	if request.Sampai == "" {
		return NewServiceError(http.StatusBadRequest, "Sampai is empty")
	}

	if checkId {
		if request.Id == 0 {
			return NewServiceError(http.StatusBadRequest, "Activity Id is empty")
		}
	}

	return nil
}

func ValidateRegister(request model.UserRequest) error {
	if request.Name == "" {
		return NewServiceError(http.StatusBadRequest, "name is empty")
	}

	if request.Email == "" {
		return NewServiceError(http.StatusBadRequest, "email is empty")
	}

	if request.Password == "" {
		return NewServiceError(http.StatusBadRequest, "password is empty")
	}

	return nil
}

func ValidateLogin(request model.UserRequest) error {
	if request.Email == "" {
		return NewServiceError(http.StatusBadRequest, "email is empty")
	}

	if request.Password == "" {
		return NewServiceError(http.StatusBadRequest, "password is empty")
	}

	return nil
}

package repository

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/marioheryanto/MIG/helper"
	"github.com/marioheryanto/MIG/model"
)

type UserRepository struct {
	DB *sql.DB
}

type UserRepositoryInterface interface {
	PingDB() error
	CreateUser(user model.User) error
	GetUserWithEmail(user *model.User) error
	GetUserWithID(id string, user *model.User) error
	CheckUserExistWithEmail(email string) (bool, error)
}

func NewUserRepository(db *sql.DB) UserRepositoryInterface {
	return UserRepository{
		DB: db,
	}
}

func (r UserRepository) PingDB() error {
	err := r.DB.Ping()
	if err != nil {
		return helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	log.Println("DB Ping success")

	return nil
}

func (r UserRepository) CreateUser(user model.User) error {
	query, args, err := squirrel.Insert("MIG.users").Columns("name,email,password").Values(user.Name, user.Email, user.Password).ToSql()
	if err != nil {
		return helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	result, err := r.DB.Exec(query, args...)
	if err != nil {
		return helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	log.Printf("%v", result)

	return nil
}

func (r UserRepository) GetUserWithEmail(user *model.User) error {
	query, args, err := squirrel.Select("id,name,password").From("MIG.users").Where(squirrel.Eq{"email": user.Email}).ToSql()
	if err != nil {
		return helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	err = r.DB.QueryRow(query, args...).Scan(&user.Id, &user.Name, &user.Password)
	if err != nil && err != sql.ErrNoRows {
		return helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	if err == sql.ErrNoRows {
		return helper.NewServiceError(http.StatusNotFound, "data not found")
	}

	return nil
}

func (r UserRepository) GetUserWithID(id string, user *model.User) error {
	query, args, err := squirrel.Select("email,name,password").From("MIG.users").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	err = r.DB.QueryRow(query, args...).Scan(&user.Email, &user.Name, &user.Password)
	if err != nil && err != sql.ErrNoRows {
		return helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	if err == sql.ErrNoRows {
		return helper.NewServiceError(http.StatusNotFound, "data not found")
	}

	return nil
}

func (r UserRepository) CheckUserExistWithEmail(email string) (bool, error) {
	var userName string

	query, args, err := squirrel.Select("name").From("MIG.users").Where(squirrel.Eq{"email": email}).ToSql()
	if err != nil {
		return false, helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	err = r.DB.QueryRow(query, args...).Scan(&userName)
	if err != nil && err != sql.ErrNoRows {
		return false, helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

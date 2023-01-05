package repository

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/marioheryanto/MIG/helper"
	"github.com/marioheryanto/MIG/model"
)

type ActivityRepository struct {
	DB *sql.DB
}

type ActivityRepositoryInterface interface {
	PingDB() error
	CreateActivity(request model.Activity, userId string) error
	EditActivity(request model.Activity) error
	DeleteActivity(Id string) error
	GetRangeActivity(from, to, id string) ([]model.Activity, error)
	GetActivityWithId(Id interface{}) (model.Activity, error)
}

func NewActivityRepository(db *sql.DB) ActivityRepositoryInterface {
	return ActivityRepository{
		DB: db,
	}
}

func (r ActivityRepository) PingDB() error {
	err := r.DB.Ping()
	if err != nil {
		return helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	log.Println("DB Ping success")

	return nil
}

func (r ActivityRepository) CreateActivity(request model.Activity, userId string) error {
	query, args, err := squirrel.
		Insert("activities").
		Columns("user_id, description, tanggal, dari, sampai, created_at").
		Values(userId, request.Description, request.Tanggal, request.Dari, request.Sampai, request.CreatedAt).
		ToSql()

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

func (r ActivityRepository) EditActivity(request model.Activity) error {
	query, args, err := squirrel.
		Update("activities").
		Set("description", request.Description).
		Set("tanggal", request.Tanggal).
		Set("dari", request.Dari).
		Set("sampai", request.Sampai).
		Set("updated_at", request.UpdatedAt).
		Where("id = ?", request.Id).
		ToSql()

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

func (r ActivityRepository) DeleteActivity(Id string) error {
	query, args, err := squirrel.
		Delete("activities").
		Where("id = ?", Id).
		ToSql()

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

func (r ActivityRepository) GetRangeActivity(from, to, id string) ([]model.Activity, error) {
	activityList := []model.Activity{}

	query, args, err := squirrel.
		Select("description, tanggal, dari, sampai, created_at, updated_at").
		From("activities").
		Where("user_id = ? AND tanggal >= ? AND tanggal <= ?", id, from, to).ToSql()

	if err != nil {
		return activityList, err
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil && err != sql.ErrNoRows {
		return activityList, err
	}

	defer rows.Close()

	for rows.Next() {
		data := model.ActivityDB{}

		err := rows.Scan(&data.Description, &data.Tanggal, &data.Dari, &data.Sampai, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			return activityList, err
		}

		activityList = append(activityList, data.Convert())
	}

	if err := rows.Err(); err != nil {
		return activityList, err
	}

	if len(activityList) == 0 {
		return activityList, helper.NewServiceError(http.StatusNotFound, "data not found")
	}

	return activityList, nil
}

func (r ActivityRepository) GetActivityWithId(Id interface{}) (model.Activity, error) {
	query, args, err := squirrel.
		Select("description, tanggal, dari, sampai, created_at, updated_at").
		From("activities").
		Where("id = ?", Id).ToSql()

	if err != nil {
		return model.Activity{}, err
	}

	dataDB := model.ActivityDB{}
	err = r.DB.QueryRow(query, args...).Scan(&dataDB.Description, &dataDB.Tanggal, &dataDB.Dari, &dataDB.Sampai, &dataDB.CreatedAt, &dataDB.UpdatedAt)
	if err != nil && err != sql.ErrNoRows {
		return model.Activity{}, err
	}

	if err == sql.ErrNoRows {
		return model.Activity{}, helper.NewServiceError(http.StatusNotFound, "data not found")
	}

	return dataDB.Convert(), nil
}

package repository

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/marioheryanto/MIG/helper"
	"github.com/marioheryanto/MIG/model"
)

type AbsensiRepository struct {
	DB *sql.DB
}

type AbsensiRepositoryInterface interface {
	PingDB() error
	CheckInOut(check string, id string) error
	GetAbsensi(tanggal, id string) (model.Absensi, error)
	GetRangeAbsensi(from, to, id string) ([]model.Absensi, error)
	CreateAbsensi(check string, id string) error
}

func NewAbsensiRepository(db *sql.DB) AbsensiRepositoryInterface {
	return AbsensiRepository{
		DB: db,
	}
}

func (r AbsensiRepository) PingDB() error {
	err := r.DB.Ping()
	if err != nil {
		return helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	log.Println("DB Ping success")

	return nil
}

func (r AbsensiRepository) CheckInOut(check string, id string) error {
	now := time.Now().In(helper.LoadLocationJakarta())
	cico := now.Format("2006-01-02 15:04:05")
	tanggal := now.Format("2006-01-02")

	builder := squirrel.Update("MIG.absensi")

	if check == "in" {
		builder = builder.Set("check_in", cico)
	} else {
		builder = builder.Set("check_out", cico)
	}

	query, args, err := builder.Where(squirrel.Eq{"user_id": id, "tanggal": tanggal}).ToSql()
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

func (r AbsensiRepository) GetAbsensi(tanggal, id string) (model.Absensi, error) {
	absensi := model.Absensi{}
	absensiDB := model.AbsensiDB{}

	query, args, err := squirrel.Select("id,check_in, check_out").From("MIG.absensi").Where(squirrel.Eq{"user_id": id, "tanggal": tanggal}).ToSql()
	if err != nil {
		return absensi, helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	err = r.DB.QueryRow(query, args...).Scan(&absensiDB.Id, &absensiDB.CheckIn, &absensiDB.CheckOut)
	if err != nil && err != sql.ErrNoRows {
		return absensi, helper.NewServiceError(http.StatusInternalServerError, err.Error())
	}

	if err == sql.ErrNoRows {
		return absensi, helper.NewServiceError(http.StatusNotFound, "data not found")
	}

	return absensiDB.Convert(), nil
}

func (r AbsensiRepository) GetRangeAbsensi(from, to, id string) ([]model.Absensi, error) {
	absensiList := []model.Absensi{}

	query, args, err := squirrel.Select("id, tanggal, check_in, check_out").From("MIG.absensi").Where("user_id = ? AND tanggal >= ? AND tanggal < ?", id, from, to).ToSql()
	if err != nil {
		return absensiList, err
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil && err != sql.ErrNoRows {
		return absensiList, err
	}

	defer rows.Close()

	for rows.Next() {
		absensiDB := model.AbsensiDB{}

		err := rows.Scan(&absensiDB.Id, &absensiDB.Tanggal, &absensiDB.CheckIn, &absensiDB.CheckOut)
		if err != nil {
			return absensiList, err
		}

		absensiList = append(absensiList, absensiDB.Convert())
	}

	if err := rows.Err(); err != nil {
		return absensiList, err
	}

	if len(absensiList) == 0 {
		return absensiList, helper.NewServiceError(http.StatusNotFound, "data not found")
	}

	return absensiList, nil
}

func (r AbsensiRepository) CreateAbsensi(check string, id string) error {
	now := time.Now().In(helper.LoadLocationJakarta())
	cico := now.Format("2006-01-02 15:04:05")
	tanggal := now.Format("2006-01-02")

	builder := squirrel.Insert("MIG.absensi")

	if check == "in" {
		builder = builder.Columns("tanggal, check_in, user_id").Values(tanggal, cico, id)
	} else {
		builder = builder.Columns("tanggal, check_out, user_id").Values(tanggal, cico, id)
	}

	query, args, err := builder.ToSql()
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

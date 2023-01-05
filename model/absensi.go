package model

import "database/sql"

type Absensi struct {
	Id       uint   `json:"-"`
	Tanggal  string `json:"tanggal"`
	CheckIn  string `json:"check_in"`
	CheckOut string `json:"check_out"`
	UserId   uint   `json:"-"`
}

type AbsensiDB struct {
	Id       uint
	Tanggal  sql.NullString
	CheckIn  sql.NullString
	CheckOut sql.NullString
	UserId   uint
}

func (s AbsensiDB) Convert() Absensi {
	return Absensi{
		Id:       s.Id,
		Tanggal:  s.Tanggal.String,
		CheckIn:  s.CheckIn.String,
		CheckOut: s.CheckOut.String,
		UserId:   s.UserId,
	}
}

package model

import "database/sql"

type Activity struct {
	Id          uint   `json:"-"`
	Description string `json:"description"`
	Tanggal     string `json:"tanggal"`
	Dari        string `json:"dari"`
	Sampai      string `json:"sampai"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ActivityDB struct {
	Id          uint
	Description sql.NullString
	Tanggal     sql.NullString
	Dari        sql.NullString
	Sampai      sql.NullString
	CreatedAt   sql.NullString
	UpdatedAt   sql.NullString
}

func (s ActivityDB) Convert() Activity {
	return Activity{
		Id:          s.Id,
		Description: s.Description.String,
		Tanggal:     s.Tanggal.String,
		Dari:        s.Dari.String,
		Sampai:      s.Sampai.String,
		CreatedAt:   s.CreatedAt.String,
		UpdatedAt:   s.UpdatedAt.String,
	}
}

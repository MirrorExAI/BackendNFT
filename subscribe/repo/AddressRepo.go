package repo

import (
	"database/sql"
	"github.com/flipped-aurora/gin-vue-admin/server/subscribe/model"
)

type AddressRepo struct {
	DB *sql.DB
}

func NewAddressRepo(db *sql.DB) *AddressRepo {
	return &AddressRepo{DB: db}
}

func (repo *AddressRepo) isExit(address string) (bool, error) {
	//
	sql := "select address from pb_users where address = " + address
	var u model.User
	repo.DB.QueryRow(sql).Scan(&u.Address).Error()

	if u.Address == "" {
		return false, nil
	}
	return true, nil
}

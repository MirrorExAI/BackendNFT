package repo

import (
	"database/sql"
	"github.com/flipped-aurora/gin-vue-admin/server/subscribe/model"
	"log"
)

type AddressRepo struct {
	DB *sql.DB
}

func NewAddressRepo(db *sql.DB) *AddressRepo {
	return &AddressRepo{DB: db}
}

func (repo *AddressRepo) isExit(address string) (bool, error) {
	//
	sql := "select address from pb_users where address = ?"
	var u model.User
	err := repo.DB.QueryRow(sql, address).Scan(&u.Address)

	if err == nil {
		log.Println(u.Address)
	}
	if u.Address == "" {
		return false, nil
	}
	return true, nil
}

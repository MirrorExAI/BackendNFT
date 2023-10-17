package repo

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

func TestIsExit(t *testing.T) {
	db, err := sql.Open("mysql", "admin:hjewaWaAs3ffff@tcp(rm-3nsn6d1g079x0d926io.mysql.rds.aliyuncs.com:3306)/nft_opensix")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := NewAddressRepo(db)

	isExit, _ := repo.isExit("0xF9c83685fecEb386F51De53cA3aCb7ed458C4f8E")

	log.Println(isExit)

	isExit2, _ := repo.isExit("0x5BE7Db87Cec1d7E1AD9F6Ae1B20D46511B077416")
	log.Println(isExit2)

}

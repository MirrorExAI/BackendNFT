package utils

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestBcryptHash(t *testing.T) {

	var password string = "123456"
	password_hash := BcryptHash(password)

	log.Println(password)
	log.Println(password_hash)
}

func TestBcryptCheck(t *testing.T) {

	var hash string = "$2a$10$UMhJVS2Cw1aeIgiTMIrgD.87Dd0ClNIi/nK3Eal.b..Wbp/F591ta"
	//var hash string = "$2a$10$ypAHaTF4Up0ctrDgfSr0QuFragqjB7ipCXtbjluXZLtwKsAeqd4FG"
	var password string = "123456"

	log.Println(BcryptCheck(password, hash))

	assert.True(t, BcryptCheck(password, hash), "密码不是123456")
}

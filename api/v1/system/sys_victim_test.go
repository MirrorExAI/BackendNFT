package system

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/tool"
	"log"
	"regexp"
	"testing"
)

func TestGenerateRandom(t *testing.T) {
	result := tool.RandFloat64(0.001234, 1)
	log.Println(result)
}

func TestValidAddres(t *testing.T) {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	fmt.Printf("is valid: %v\n", re.MatchString("0x5052864d46b4B605DD4fC1fd14F65F09302C9325")) // is valid: true
}

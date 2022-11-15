package session

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func GetToken() string {
	return viper.GetViper().GetString("auth.token")
}

func ApplyAuthorization(req *http.Request) {
	token := GetToken()
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

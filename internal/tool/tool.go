package tool

import "regexp"

func CheckPassword(password string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{8,16}$", password); !ok {
		return false
	}
	return true
}

package utils

import "net/http"

func CheckInternet() bool {
	res, err := http.Get("http://connectivitycheck.gstatic.com/generate_204")
	if err != nil {
		return false
	}

	res.Body.Close()

	return res.StatusCode == 204
}

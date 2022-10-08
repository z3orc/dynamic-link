package util

import (
	"errors"
	"fmt"
	"net/http"
)

func CheckUrl(url string) (error) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("No response from request")
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	} else {
		return nil
	}
}
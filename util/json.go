package util

import (
	"fmt"
	"io"
	"net/http"
)

func GetJson(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("No response from request")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, err
}
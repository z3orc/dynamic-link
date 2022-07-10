package fetch

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetJson(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	return string(body), err
}
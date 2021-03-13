package Request

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(Url string, Query map[string]string) (string, error) {
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.26 Safari/537.36")
	if Query != nil {
		q := req.URL.Query()
		for k, v := range Query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

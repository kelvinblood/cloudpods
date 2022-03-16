package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

// huawei fusion compute http client
var (
	endpoint = `https://10.10.1.52:7443`
	user     = `ronly`
	password = `Admin@huawei123$`
)

type CloudProvider struct {
	EndPoint string
	User     string
	Password string
}

func main() {
	cp := CloudProvider{
		EndPoint: endpoint,
		User:     user,
		Password: password,
	}
	dat := _post(cp, "/service/session")
	fmt.Printf("%v\n", string(dat))
	dat = _get(cp, "/service/versions")
	fmt.Printf("%v\n", string(dat))
}

func _get(cp CloudProvider, uri string) []byte {
	req, err := http.NewRequest("GET", cp.EndPoint+uri, nil)
	if err != nil {
		panic(err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
		return nil
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
		return nil
	}
	defer resp.Body.Close()

	return result
}

func _post(cp CloudProvider, uri string) []byte {
	req, err := http.NewRequest("POST", cp.EndPoint+uri, nil)
	if err != nil {
		panic(err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json;version=6.0;charset=UTF-8")
	req.Header.Set("X-Auth-User", cp.User)
	req.Header.Set("X-Auth-Key", cp.Password)
	req.Header.Set("X-ENCRIPT-ALGORITHM", "1")
	req.Header.Set("X-Auth-UserType", "0")
	// 跳过https认证
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
		return nil
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
		return nil
	}
	defer resp.Body.Close()
	fmt.Println(resp.Header.Get("X-Auth-Token"))

	return result
}

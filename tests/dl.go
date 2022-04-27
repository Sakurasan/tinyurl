package main

import (
	"fmt"
	"net/http"
)

func main() {
	dlpath := "http://pan.oneisall.xyz/Outline.apk"
	req, _ := http.NewRequest(http.MethodHead, dlpath, nil)
	rsp, _ := http.DefaultClient.Do(req)
	defer rsp.Body.Close()
	fmt.Println(rsp.Status)
	// dump, _ := httputil.DumpResponse(rsp, false)
	// fmt.Println(string(dump))
	// if rsp.Header != nil {
	// fmt.Printf("%v", rsp.Header)
	for k, v := range rsp.Header {
		fmt.Println(k, v)
	}
	fmt.Println(rsp.Header.Get("Content-Length"))
	fmt.Println(rsp.Header.Get("Accept-Ranges"))

}

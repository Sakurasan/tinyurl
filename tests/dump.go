package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bytereq, _ := httputil.DumpRequest(r, true)
		bytebody, _ := io.ReadAll(r.Body)
		log.Println(string(bytereq))
		m := map[string]interface{}{
			"host":       r.Host,
			"requestUrl": r.RequestURI,
			"proto":      r.Proto,
			"path":       r.URL.String(),
			"method":     r.Method,
			"body":       string(bytebody),
			"header":     r.Header,
			"form":       r.Form,
			"postform":   r.PostForm,
			"dumpreq":    string(bytereq),
		}
		bj, _ := json.Marshal(m)
		fmt.Println(string(bj))
		// json.NewEncoder(w).Encode(m)

		for k, v := range cors {
			w.Header().Set(k, v)
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, string(bj))
		return
	})
	server := http.Server{
		Addr:    ":890",
		Handler: mux,
	}
	server.ListenAndServe()
}

var cors = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Methods":     "POST,OPTIONS,GET",
	"Access-Control-Max-Age":           "3600",
	"Access-Control-Allow-Headers":     "accept,x-requested-with,Content-Type",
	"Access-Control-Allow-Credentials": "true",
}

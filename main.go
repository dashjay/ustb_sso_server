package main

import (
	"fmt"
	"net/http"

	"ustb_sso/auth_hub"
	"ustb_sso/env"
)

func main() {
	// http接口
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("fuck") != "fuck" {
			w.WriteHeader(666)
			return
		}
		unionId := r.FormValue("union_id")
		res := auth_hub.DoAuth(unionId)
		rb, _ := res.MarshalJSON()
		_, _ = w.Write(rb)
		return
	})

	http.HandleFunc("/func", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("fuck") != "fuck" {
			w.WriteHeader(666)
			return
		}
		unionId := r.FormValue("union_id")
		funcName := r.FormValue("func")
		if unionId == "" || funcName == "" {
			w.Write([]byte("unionId or funcName empty"))
			w.WriteHeader(400)
			return
		}
		res, err := auth_hub.Func(funcName, unionId)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(500)
			return
		}
		w.Write(res)
		return
	})

	http.ListenAndServe(fmt.Sprintf(":%s", env.Port), nil)
}

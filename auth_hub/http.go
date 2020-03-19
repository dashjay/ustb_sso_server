package auth_hub

import (
	"net/http"
)

func DoAuthHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("fuck") != "fuck" {
		w.WriteHeader(666)
		return
	}
	unionId := r.FormValue("union_id")
	res := doAuth(unionId)
	rb, _ := res.MarshalJSON()
	_, _ = w.Write(rb)
	return
}

func DoFuncHTTP(w http.ResponseWriter, r *http.Request) {
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
	res, err := doFunc(funcName, unionId)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	w.Write(res)
	return
}

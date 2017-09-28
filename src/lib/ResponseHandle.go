package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/***
*******处理错误信息********
**/

type ResHandle struct {
	Error error    //错误
	Res   Response //返回信息
}
type ReqHandler func(http.ResponseWriter, *http.Request) *ResHandle

func (fn ReqHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		resBytes, err := json.Marshal(err.Res)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(resBytes)
	}
}

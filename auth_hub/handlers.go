package auth_hub

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/boltdb/bolt"

	"ustb_sso/constant"
	"ustb_sso/session"
	"ustb_sso/structs"
)

var hub *Hub

func init() {
	hub = Newhub()
}

var (
	ErrorNotAuth   = errors.New("user do not auth")
	ErrorBoltDB    = errors.New("boltdb cracked")
	ErrorUnmarshal = errors.New("ac unmarshal error")
)

// DoAuth 认证
func DoAuth(unionId string) structs.AuthStruct {

	if unionId == "" {
		return structs.AuthStruct{Code: 1, Msg: "unionId empty"}
	}
	var ac structs.AuthCookie

	// 检查是否已经认证
	err := session.GetDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket(session.DBCookies)
		if b == nil {
			return ErrorBoltDB
		}
		byteac := b.Get([]byte(unionId))
		if byteac == nil {
			return ErrorNotAuth
		}
		err := ac.UnmarshalJSON(byteac)
		if err != nil {
			return ErrorUnmarshal
		}
		return nil
	})

	if err != nil {
		log.Println(err)
		switch err.Error() {
		case ErrorNotAuth.Error():
			{
				// 并没有认证
				break
			}
		default:
			return structs.AuthStruct{Code: 1, Msg: err.Error()}
		}
	} else {
		log.Println("已认证")
		// 已经认证了
		return structs.AuthStruct{Code: 101, URL: "", UID: ac.UID}
	}

	// 开始认证
	reschan, errchan := hub.HandlerAuth(unionId)
	select {
	case resStr := <-reschan:
		{

			return structs.AuthStruct{Code: 102, URL: resStr, UID: ""}
		}
	case errStr := <-errchan:
		{
			return structs.AuthStruct{Code: 1, Msg: errStr.Error()}
		}
	}
}

// Func 使用学生本来的成绩来获取
func Func(funcname, unionId string) ([]byte, error) {

	if unionId == "" {
		return []byte{}, errors.New("unionId empty")
	}

	var ac structs.AuthCookie
	err := session.GetDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket(session.DBCookies)
		if b != nil {
			res := b.Get([]byte(unionId))
			if res == nil {
				return errors.New("not found")
			}
			err := ac.UnmarshalJSON(res)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return []byte{}, err
	}

	var req *http.Request

	switch funcname {

	case "grade":
		req = NewReq("GET", constant.GetGradeURL, ac)
		break
	case "course":
		req = NewReq("GET", constant.GetCourseTable, ac)
		break
	case "unbind":
		err := session.GetDb().Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(session.DBCookies)
			if b == nil {
				return errors.New("bucket not exists")
			} else {
				err := b.Delete([]byte(unionId))
				if err != nil {
					return err
				} else {
					return nil
				}
			}
		})

		if err != nil {

			return []byte{}, err
		} else {

			return json.Marshal(&struct {
				Msg string `json:"msg"`
			}{Msg: "解绑成功，部分功能将会失效。感谢以下开发者为大家提供该服务\n\n1. DJ\n2. XRM"})
		}
	}

	body, err := hub.DoGetBody(req)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

package auth_hub

import (
	"encoding/json"
	"errors"
	"fmt"
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

// doAuth 认证
func doAuth(unionId string) structs.AuthStruct {

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

// doFunc 使用学生本来的成绩来获取
func doFunc(funcName, unionId string) ([]byte, error) {

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

	var res []byte
	err = session.GetDb().View(func(tx *bolt.Tx) error {
		b := tx.Bucket(session.DBCache)
		if b != nil {
			res = b.Get([]byte(fmt.Sprintf("%s:%s", unionId, funcName)))
			if res == nil {
				return errors.New("not found")
			}
		}
		return nil
	})

	if err == nil && res != nil {
		return res, nil
	}

	switch funcName {

	case "grade":
		req = NewReq("GET", constant.GetGradeURL, ac)
		break
	case "course":
		req = NewReq("GET", constant.GetCourseTable, ac)
		break
	case "cet":
		req = NewReq("POST", constant.GetCETScore, ac)
		break
	case "head":
		req = NewReq("POST", constant.GetHeadImageURL, ac)
		break
	case "info":
		req = NewReq("GET", constant.GetStuInfo, ac)
		break
	case "exam":
		req = NewReq("GET", constant.GetExamArrangement, ac)
		break
	case "infomore":
		req = NewReq("GET", constant.GetStuInfoPerfection, ac)
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
	default:
		return []byte("func not exists"), nil
	}

	body, err := hub.DoGetBody(req)
	if err != nil {
		return []byte{}, err
	}

	defer session.GetDb().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(session.DBCache)
		if b == nil {
			return err
		}
		err = b.Put([]byte(fmt.Sprintf("%s:%s", unionId, funcName)), body)
		return err
	})

	return body, nil
}

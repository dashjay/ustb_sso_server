package auth_hub

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/boltdb/bolt"

	"ustb_sso/constant"
	"ustb_sso/session"
	"ustb_sso/structs"
)

type Hub struct {
	Workers *sync.Pool
}

// DoGetBody 从Hub中Worker池中获取一个Worker做请求并且获取Body返回，
// 且归还该Client
func (h *Hub) DoGetBody(req *http.Request) ([]byte, error) {
	// 请求一个HttpClient
	worker := h.Workers.Get().(*Worker)
	defer h.Workers.Put(worker)
	// 做请求
	resp, err := worker.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	// 归还该HttpClient

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}

func (h *Hub) DoGetResp(req *http.Request) (*http.Response, error) {
	worker := h.Workers.Get().(*Worker)
	defer h.Workers.Put(worker)

	resp, err := worker.client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	return resp, nil
}

var (
	// 从页面中提取obj的
	objRegexp = regexp.MustCompile(`objLogin=(.*);`)
	// 将key没有双引号的jsonobj添加上双引号
	// 从教务SIS中获取的js字典是{id:"1",key:"value"}这种形式，golang的解析器无法解析(UnMarshal)这种形式的内容进入结构中
	jsonKey     = regexp.MustCompile(`[\w]+[:]`)
	jsonKeyFunc = func(s []byte) []byte {
		reg2 := regexp.MustCompile(`[\w]+`)
		s = reg2.FindAll(s, -1)[0]
		temp := []byte(`"` + string(s) + `":`)
		temp = bytes.ReplaceAll(temp, []byte(`"https"`), []byte(`https`))
		temp = bytes.ReplaceAll(temp, []byte(`"http"`), []byte(`http`))
		return temp
	}

	// 从回应中获取sid
	sidRegexp = regexp.MustCompile(`qrimg\?sid=(.*)"\s*w`)

	// 从Head中剥离JSESSIONID 和 UID
	jseRegexp = regexp.MustCompile(`JSESSIONID=(.*);\s*P`)
	UIDRegexp = regexp.MustCompile(`UID=(.*);\sV`)
	TSRegexp  = regexp.MustCompile(`TS=(.*);\sV`)
	TKNRegexp = regexp.MustCompile(`TKN=(.*);\sV`)
)

func Newhub() *Hub {
	return &Hub{Workers: &sync.Pool{New: func() interface{} {
		return NewWork()
	}}}
}

func (h *Hub) HandlerAuth(unionId string) (<-chan string, <-chan error) {
	reschan := make(chan string)
	errchan := make(chan error)
	go func() {

		req, err := http.NewRequest("GET", constant.GetRandTokenURL, nil)
		if err != nil {
			errchan <- err
			return
		}
		content, err := h.DoGetBody(req)
		if err != nil {
			errchan <- err
			return
		}
		temp := objRegexp.FindSubmatch(content)
		if len(temp) < 2 {
			errchan <- errors.New("json object not found")
			return
		}

		jsonObj := jsonKey.ReplaceAllFunc(temp[1], jsonKeyFunc)
		var apt structs.AppToken
		err = apt.UnmarshalJSON(jsonObj)
		if err != nil {
			errchan <- err
			return
		}
		URL := fmt.Sprintf("%s?%s", constant.SISGetQr, apt.Encode())
		req, err = http.NewRequest("GET", URL, nil)
		if err != nil {
			errchan <- err
			return
		}
		content, err = h.DoGetBody(req)
		if err != nil {
			errchan <- err
			return
		}

		temp = sidRegexp.FindSubmatch(content)
		if len(temp) < 2 {
			errchan <- errors.New("sid not found")
			return
		}
		sid := temp[1]
		URL = fmt.Sprintf("%s?sid=%s", constant.SISAuth, sid)
		reschan <- URL

		var authCode string

		var s structs.State
		for {
			URL = fmt.Sprintf("%s?sid=%s", constant.SISGetState, sid)
			req, err := http.NewRequest("GET", URL, nil)
			if err != nil {
				log.Printf("获取Client出错：%v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			content, err = h.DoGetBody(req)
			if err != nil {
				log.Printf("请求状态错误：%v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			err = s.UnmarshalJSON(content)
			if err != nil {
				log.Printf("解析状态错误：%v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if s.State == 200 {
				// 得到AuthCode
				authCode = s.Data
				break
			} else if s.State == 101 {
				log.Println("用户未操作")
				time.Sleep(1 * time.Second)
				continue
			} else if s.State == 102 {
				log.Println("用户已扫码")
				time.Sleep(1 * time.Second)
				continue
			} else if s.State == 103 {
				log.Println("sessionid失效")
				return
			} else if s.State == 104 {
				log.Println("用户操作超期")
				return
			}

		}
		URL = fmt.Sprintf(constant.JWSTU, authCode, apt.RandToken)
		req, err = http.NewRequest("GET", URL, nil)
		if err != nil {
			errchan <- err
			return
		}

		resp, err := h.DoGetResp(req)
		if err != nil {
			errchan <- err
			return
		}
		defer resp.Body.Close()
		JSESSIONID := ""
		UID := ""
		TS := ""
		TKN := ""
		for _, v := range resp.Header {

			for _, sk := range v {

				temp := jseRegexp.FindStringSubmatch(sk)
				if len(temp) > 1 {
					JSESSIONID = temp[1]
					continue
				}

				temp = UIDRegexp.FindStringSubmatch(sk)
				if len(temp) > 1 {
					UID = temp[1]
					continue
				}

				temp = TSRegexp.FindStringSubmatch(sk)
				if len(temp) > 1 {
					TS = temp[1]
					continue
				}

				temp = TKNRegexp.FindStringSubmatch(sk)
				if len(temp) > 1 {
					TKN = temp[1]
					continue
				}
			}
		}
		// fmt.Printf("JSESSIONID:%s,UID:%s,TS:%s,TKN:%s\n", JSESSIONID, UID, TS, TKN)
		cookie := structs.AuthCookie{JSESSIONID: JSESSIONID, TS: TS, TKN: TKN, UID: UID}

		res, err := cookie.MarshalJSON()
		if err != nil {
			log.Println(err)
			errchan <- err
			return
		}

		_ = session.GetDb().Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(session.DBCookies)
			if b != nil {
				err := b.Put([]byte(unionId), res)
				if err != nil {
					panic(err)
				}

				log.Printf("已存%s, %d", UID, len(res))
			} else {
				panic("bucket empty")
			}
			return nil
		})

	}()
	return reschan, errchan
}

func NewReq(method, url string, ac structs.AuthCookie) *http.Request {
	var req *http.Request

	req, _ = http.NewRequest(method, url, nil)

	req.AddCookie(&http.Cookie{Name: "JSESSIONID", Value: ac.JSESSIONID, HttpOnly: false})
	req.AddCookie(&http.Cookie{Name: "TS", Value: ac.TS})
	req.AddCookie(&http.Cookie{Name: "TKN", Value: ac.TKN})
	req.AddCookie(&http.Cookie{Name: "UID", Value: ac.UID, HttpOnly: false})

	return req
}

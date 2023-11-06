package HttpClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chaos-star/marvel/Log"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

type SendParamType int

const (
	Bytes SendParamType = iota
	Json
	Query
	Form
)

type HttpClient struct {
	log            Log.ILogger
	DefaultTimeout time.Duration
}

func Initialize(log Log.ILogger) *HttpClient {
	return &HttpClient{log: log, DefaultTimeout: time.Second * 60}
}

type Request struct {
	*http.Request
	log     Log.ILogger
	timeout time.Duration
	params  []byte
}

func (h *HttpClient) NewRequest(url string, method string, timeout time.Duration) *Request {
	req, _ := http.NewRequest(method, url, nil)
	return &Request{req, h.log, timeout, nil}
}

func (r *Request) Bind(paramType SendParamType, params interface{}) (err error, req *Request) {
	var body []byte
	err, body = r.parse(paramType, "", params)
	if err != nil {
		return
	}
	r.params = body
	if len(body) > 0 {
		if paramType == Form {
			formBody, _ := url.ParseQuery(string(body))
			r.Form = formBody
		}

		if paramType == Query {
			newUrl, _ := url.Parse(fmt.Sprintf("%s?%s", r.URL.String(), string(body)))
			r.URL = newUrl
		}

		if paramType == Json || paramType == Bytes {
			reflect.TypeOf(body)
			r.Body = io.NopCloser(bytes.NewReader(body))
		}

		if paramType == Bytes {
			r.Body = io.NopCloser(bytes.NewReader(body))
		}
	}
	req = r
	return
}

func (r *Request) Send(response interface{}, isParse bool) (err error, state int) {
	var (
		client http.Client
	)
	client.Timeout = r.timeout
	headerBody, _ := json.Marshal(r.Header.Clone())
	reqLog := fmt.Sprintf("Url:%s | Method:%s | Header:%s | Params:%s", r.URL.String(), r.Method, string(headerBody), string(r.params))
	r.log.Info(fmt.Sprintf("[HTTP_SEND_REQUEST] %s", reqLog))
	resp, err := client.Do(r.Request)
	if err != nil {
		r.log.Info(fmt.Sprintf("[HTTP_SEND_RESPONSE] [Exception] %s | Error:%s ", reqLog, err.Error()))
		return
	}
	respData, _ := io.ReadAll(resp.Body)
	state = resp.StatusCode
	r.log.Info(fmt.Sprintf("[HTTP_SEND_RESPONSE] %s | HttpStatus:%d | Response:%s", reqLog, state, string(respData)))
	if isParse {
		err = json.Unmarshal(respData, response)
		if err != nil {
			return
		}
	} else {
		response = respData
	}
	return
}

func (r *Request) parse(paramType SendParamType, prefix string, params interface{}) (err error, data []byte) {
	if params == nil {
		return
	}
	switch paramType {
	case Bytes:
		if x, ok := params.([]byte); ok {
			data = x
		} else {
			err = errors.New("incorrect data type")
		}
	case Json:
		var x []byte
		if x, err = json.Marshal(params); err == nil {
			data = x
		} else {
			err = errors.New(fmt.Sprintf("Failed to build json :%s", err.Error()))
		}
	case Query:
		fallthrough
	case Form:
		var (
			rVal    = reflect.ValueOf(params)
			urlArgs = url.Values{}
			part    []byte
			group   []byte
		)
		rType := rVal.Kind().String()
		if rType == "ptr" {
			rType = rVal.Elem().Kind().String()
			rVal = rVal.Elem()
		}
		switch rType {
		case "struct":
			fieldNum := rVal.NumField()
			for i := 0; i < fieldNum; i++ {
				vItem := rVal.Field(i)
				tItem := rVal.Type().Field(i)
				if tItem.Anonymous {
					continue
				}

				itemType := vItem.Kind().String()
				if itemType == "chan" || itemType == "func" {
					continue
				}
				fieldName := tItem.Tag.Get("form")
				if fieldName == "" {
					fieldName = tItem.Name
				}
				if fieldName != "" && prefix != "" {
					fieldName = fmt.Sprintf("%s.%s", prefix, fieldName)
				}
				if itemType == "ptr" || itemType == "struct" || itemType == "map" || itemType == "slice" || itemType == "array" {
					if itemType == "ptr" {
						err, part = r.parse(paramType, fieldName, vItem.Elem().Interface())
					} else {
						err, part = r.parse(paramType, fieldName, vItem.Interface())
					}

					if err != nil {
						return
					}

					part = append([]byte("&"), part...)
					group = append(part, group...)
					continue
				}
				urlArgs.Add(fieldName, fmt.Sprintf("%v", vItem.Interface()))
			}
		case "ptr":
			err, part = r.parse(paramType, prefix, rVal.Elem())
			if err != nil {
				return
			}
		case "map":
			keys := rVal.MapKeys()
			for i := 0; i < rVal.Len(); i++ {
				kv := keys[i]
				if kv.Kind().String() != "string" {
					continue
				}
				vItem := rVal.MapIndex(kv)
				itemType := vItem.Kind().String()
				fieldName := fmt.Sprintf("%v", kv)
				if fieldName != "" && prefix != "" {
					fieldName = fmt.Sprintf("%s.%s", prefix, fieldName)
				}
				if itemType == "ptr" || itemType == "struct" || itemType == "map" || itemType == "slice" || itemType == "array" {
					if itemType == "ptr" {
						err, part = r.parse(paramType, fieldName, vItem.Elem().Interface())
					} else {
						err, part = r.parse(paramType, fieldName, vItem.Interface())
					}

					if err != nil {
						return
					}
					if len(group) > 0 {
						part = append([]byte("&"), part...)
					}
					group = append(group, part...)
					continue
				}
				urlArgs.Add(fieldName, fmt.Sprintf("%v", vItem.Interface()))
			}
		case "array":
			fallthrough
		case "slice":
			for i := 0; i < rVal.Len(); i++ {
				vItem := rVal.Index(i)
				itemType := vItem.Kind().String()
				fieldName := prefix
				if fieldName == "" {
					fieldName = "arr"
				}
				if itemType == "ptr" || itemType == "struct" || itemType == "map" || itemType == "slice" || itemType == "array" {
					if itemType == "ptr" {
						err, part = r.parse(paramType, fieldName, vItem.Elem().Interface())
					} else {
						err, part = r.parse(paramType, fieldName, vItem.Interface())
					}

					if err != nil {
						return
					}
					if len(group) > 0 {
						part = append([]byte("&"), part...)
					}
					group = append(group, part...)
					continue
				}
				urlArgs.Add(fieldName, fmt.Sprintf("%v", vItem.Interface()))
			}
		default:
			return
		}
		group = append([]byte(urlArgs.Encode()), group...)
		data = group
	default:
		err = errors.New("unsupported data transfer type")
	}
	return
}

/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-17 15:24
* Description:
*****************************************************************/

package xhttpServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-xe2/x/encoding/xparser"
	"github.com/go-xe2/x/type/xstring"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type THttpProcessStatusCode int

const (
	HPS_END_WRITE THttpProcessStatusCode = iota
	HPS_UN_SUPPORT
	HPS_INNER_ERROR
)

type HttpProcessStatus interface {
	Code() THttpProcessStatusCode
	Error() error
}

type THttpProcessStatus struct {
	code THttpProcessStatusCode
	err  error
}

func NewHttpErrorStatus(err error) HttpProcessStatus {
	return &THttpProcessStatus{
		code: HPS_INNER_ERROR,
		err:  err,
	}
}

func NewHttpProcessStatus(code THttpProcessStatusCode) HttpProcessStatus {
	return &THttpProcessStatus{
		code: code,
		err:  nil,
	}
}

func (p *THttpProcessStatus) Code() THttpProcessStatusCode {
	return p.code
}

func (p *THttpProcessStatus) Error() error {
	return p.err
}

type THttpRequest struct {
	// 路由名称
	routerName string
	// 路由根路径
	baseRouter string
	// 请求方法
	method string // get, post, delete,option,put
	req    *http.Request
	writer http.ResponseWriter
	// 请求参数
	params       map[string]interface{}
	headers      map[string][]string
	files        []File
	buf          *bytes.Buffer
	headerParams map[string]string
}

func NewHttpRequest(baseRouter string, req *http.Request, writer http.ResponseWriter) *THttpRequest {
	inst := &THttpRequest{
		req:          req,
		baseRouter:   baseRouter,
		writer:       writer,
		params:       make(map[string]interface{}),
		buf:          bytes.NewBuffer([]byte{}),
		headerParams: make(map[string]string),
	}
	return inst.parseParams()
}

func (p *THttpRequest) parseParams() *THttpRequest {
	p.method = p.req.Method
	p.headers = p.req.Header
	contentType := p.req.Header.Get("Content-Type")
	// application/json;utf-8
	// application/x-www-form-urlencoded
	// multipart/form-data
	// 处理域名
	p.routerName = p.req.URL.Path
	if xstring.StartWith(p.routerName, p.baseRouter) {
		p.routerName = p.routerName[len(p.baseRouter):]
	}

	// 处理参数
	// 使用请求头参数
	for k, v := range p.headers {
		c := len(v)
		if c == 0 {
			continue
		}
		if hp, ok := p.headerParams[k]; ok {
			if c == 1 {
				p.params[hp] = v[0]
			} else {
				p.params[hp] = v
			}
		}
	}
	// 处理get参数,get方法提交时，不处理body中的数据
	values := p.req.URL.Query()
	for k, v := range values {
		c := len(v)
		if c == 0 {
			continue
		}
		if c == 1 {
			p.params[k] = v[0]
		} else {
			p.params[k] = v
		}
	}
	//
	//if strings.EqualFold(p.method, "get") {
	//	// 处理get参数,get方法提交时，不处理body中的数据
	//	values := p.req.URL.Query()
	//	for k, v := range values {
	//		c := len(v)
	//		if c == 0 {
	//			continue
	//		}
	//		if c == 1 {
	//			p.params[k] = v[0]
	//		} else {
	//			p.params[k] = v
	//		}
	//	}
	//} else {
	if !strings.EqualFold(p.method, "get") {
		// multipart方式提交
		if strings.Index(contentType, "multipart/form-data") != -1 {
			if err := p.req.ParseMultipartForm(p.req.ContentLength); err == nil {
				form := p.req.MultipartForm
				if form != nil {
					values := form.Value
					for k, item := range values {
						c := len(item)
						if c == 0 {
							continue
						}
						if c == 1 {
							p.params[k] = item[0]
						} else {
							p.params[k] = item
						}
					}
					p.files = make([]File, 0)
					files := form.File
					for _, items := range files {
						for _, hr := range items {
							p.files = append(p.files, NewHttpFile(hr))
						}
					}
				}
			}
		} else {
			// 处理post表单参数
			// application/x-www-form-urlencoded
			if e := p.req.ParseForm(); e == nil {
				values := p.req.Form
				for k, item := range values {
					c := len(item)
					if c == 0 {
						continue
					}
					if c == 1 {
						p.params[k] = item[0]
					} else {
						p.params[k] = item
					}
				}
			}

			if strings.Index(contentType, "application/json") >= 0 || strings.Index(contentType, "application/xml") >= 0 {
				// 处理body的json参数
				body, err := ioutil.ReadAll(p.req.Body)
				if err == nil {
					parser, err := xparser.LoadContent(body)
					if err == nil {
						mp := parser.ToMap()
						for k, v := range mp {
							p.params[k] = v
						}
					} else {
						fmt.Println("parse json error1:", err)
					}
				} else {
					fmt.Println("parse json error:", err)
				}
			}
		}
	}
	return p
}

func (p *THttpRequest) GetFiles() []File {
	return p.files
}

func (p *THttpRequest) GetParams() map[string]interface{} {
	return p.params
}

func (p *THttpRequest) GetMethod() string {
	return p.method
}

func (p *THttpRequest) GetRouterName() string {
	return p.routerName
}

func (p *THttpRequest) GetBaseRouter() string {
	return p.baseRouter
}

func (p *THttpRequest) GetHost() string {
	return p.req.Host
}

func (p *THttpRequest) GetUrl() *url.URL {
	return p.req.URL
}

func (p *THttpRequest) GetRequest() *http.Request {
	return p.req
}

func (p *THttpRequest) GetResponse() http.ResponseWriter {
	return p.writer
}

func (p *THttpRequest) WriteString(str string) {
	p.Write([]byte(str))
}

func (p *THttpRequest) ReturnText(str string) {
	p.SetHeader("Content-Type", "text/plain; charset=utf-8")
	p.WriteString(str)
	p.Exit()
}

func (p *THttpRequest) ReturnJson(obj interface{}) {
	p.SetHeader("Content-Type", "application/json; charset=utf-8")
	bts, err := json.Marshal(obj)
	if err != nil {
		panic(NewHttpErrorStatus(err))
	}
	p.WriteString(string(bts))
	p.Exit()
}

func (p *THttpRequest) Write(data []byte) {
	if _, err := p.buf.Write(data); err != nil {
		panic(NewHttpErrorStatus(err))
	}
}

func (p *THttpRequest) ResponseUnSupport() {
	panic(NewHttpProcessStatus(HPS_UN_SUPPORT))
}

func (p *THttpRequest) Exit() {
	panic(NewHttpProcessStatus(HPS_END_WRITE))
}

func (p *THttpRequest) Clear() {
	p.buf = bytes.NewBuffer([]byte{})
}

func (p *THttpRequest) Buffer() *bytes.Buffer {
	return p.buf
}

func (p *THttpRequest) WriteHeader(statusCode int) {
	p.writer.WriteHeader(statusCode)
}

func (p *THttpRequest) SetHeader(key string, value string) {
	p.writer.Header().Set(key, value)
}

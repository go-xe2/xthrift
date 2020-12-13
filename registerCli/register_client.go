/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-19 11:38
* Description:
*****************************************************************/

package registerCli

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-xe2/x/encoding/xbase64"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
	"io/ioutil"
	"net/http"
)

type TRegisterClient struct {
	host string
}

func NewRegisterClient(host string) *TRegisterClient {
	return &TRegisterClient{
		host: host,
	}
}

func (p *TRegisterClient) Register(host string, port int, project *pdl.FileProject) error {
	if project == nil {
		return errors.New("project为nil")
	}
	writer := bytes.NewBuffer([]byte{})
	if err := project.SaveProject(writer); err != nil {
		return err
	}
	data := xbase64.Encode(writer.Bytes())
	if err := p.postData(host, port, string(data)); err != nil {
		return err
	}
	return nil
}

func (p *TRegisterClient) postData(host string, port int, data string) error {
	postData := make(map[string]interface{})
	postData["host"] = host
	postData["port"] = port
	postData["pdl"] = data
	bts, err := json.Marshal(postData)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(bts)

	resp, err := http.Post(p.host, "application/json; charset=utf-8", body)
	if err != nil {
		return errors.New("网络错误:" + err.Error())
	}
	defer resp.Body.Close()

	resBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(string(resBytes))
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return err
	}
	status := t.Int(result["status"])
	msg := t.String(result["msg"])
	if status > 0 {
		return errors.New("注册出错:" + msg)
	}
	return nil
}

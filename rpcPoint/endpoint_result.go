/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-18 15:19
* Description:
*****************************************************************/

package rpcPoint

type TResResult struct {
	Status  int         `json:"status"`
	Msg     string      `json:"msg"`
	Content interface{} `json:"content,omitempty"`
}

func NewResResult(status int, msg string, data interface{}) *TResResult {
	return &TResResult{
		Status:  status,
		Msg:     msg,
		Content: data,
	}
}

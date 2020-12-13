/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-09-01 17:08
* Description:
*****************************************************************/

package pdlQrySvc

type TResData struct {
	Status  int         `json:"status"`
	Msg     string      `json:"msg"`
	Content interface{} `json:"content,omitempty"`
}

func MakeResData(status int, msg string, data interface{}) *TResData {
	return &TResData{
		Status:  status,
		Msg:     msg,
		Content: data,
	}
}

/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-07 12:15
* Description:
*****************************************************************/

package netstream

import (
	"io"
	"net/http"
)

type StreamWriterFlusher interface {
	io.Writer
	http.Flusher
}

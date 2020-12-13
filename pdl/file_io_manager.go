/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-20 16:09
* Description:
*****************************************************************/

package pdl

import "io"

type FileIOManager interface {
	Create(ns string, fileName string) (io.Writer, error)
	Close(ns string, fileName string)
}

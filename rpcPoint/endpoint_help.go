/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-19 16:04
* Description:
*****************************************************************/

package rpcPoint

import (
	"fmt"
	"github.com/go-xe2/xthrift/xhttpServer"
	"sort"
	"strings"
)

// handler /pdl/help
func (p *TEndPointServer) help(req *xhttpServer.THttpRequest) {
	req.SetHeader("Content-Type", "text/html; charset=utf-8")
	req.WriteString("<html lang=\"zh-cn\">")
	req.WriteString("<head>")
	req.WriteString("<meta charset=\"utf-8\" />")
	req.WriteString("<title>服务接口说明</title>")
	req.WriteString("</head>")
	req.WriteString("<body>")
	p.writeHelpMethods(req)
	req.WriteString("</body>")
	req.WriteString("</html>")
	req.Exit()
}

func (p *TEndPointServer) writeHelpMethods(req *xhttpServer.THttpRequest) {

	req.WriteString("<div class=\"help-title\" style=\"font-size:18;color:#333;\">服务接口说明</div>")
	routers := p.server.Routers()
	i := 0
	size := len(routers)
	methods := make([][]interface{}, 0)
	for _, m := range routers {
		methods = append(methods, []interface{}{m.GetPattern(), m})
	}
	sort.Slice(methods, func(i, j int) bool {
		return strings.Compare(methods[i][0].(string), methods[j][0].(string)) < 0
	})

	for _, items := range methods {
		p.writeHelpMethod(req, items[1].(*xhttpServer.THttpRouterInfo))
		i++
		if i < size {
			p.writeMethodSep(req)
		}
	}
}

func (p *TEndPointServer) writeHelpMethod(req *xhttpServer.THttpRequest, route *xhttpServer.THttpRouterInfo) {
	req.WriteString("<div class=\"m-method\">")
	req.WriteString(fmt.Sprintf("<div class=\"m-title\">%s</div>", route.GetTitle()))
	req.WriteString(fmt.Sprintf("<div class=\"m-addr\"><span>地址:</span><span>%s</span></div>", route.GetPattern()))
	req.WriteString(fmt.Sprintf("<div class=\"m-summary\"><div class=\"m-s-title\">说明:</div><div class=\"m-s-cxt\">%s</div></div>", route.GetSummary()))
	req.WriteString("</div>")
}

func (p *TEndPointServer) writeMethodSep(req *xhttpServer.THttpRequest) {
	req.WriteString("<span class=\"m-h-line\" style=\"font-size:1;height:1;background:#ccc;display:block;\">&nbsp;</span>")
}

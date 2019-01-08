package goplugin

import (
	"fmt"
	"net/http"
	"strings"
	"vuldb/common"
	"vuldb/plugin"
)

type axisWeakPass struct {
	info   common.PluginInfo
	result []common.PluginInfo
}

func init() {
	plugin.Regist("axis", &axisWeakPass{})
}
func (d *axisWeakPass) Init() common.PluginInfo {
	d.info = common.PluginInfo{
		Name:    "Axis2控制台 弱口令",
		Remarks: "攻击者通过此漏洞可以登陆管理控制台，通过部署功能可直接获取服务器权限。",
		Level:   0,
		Type:    "WEAK",
		Author:   "wolf",
		References: common.References{
			URL: "http://sadadas.com",
			CVE: "test",
		},
	}
	return d.info
}
func (d *axisWeakPass) GetResult() []common.PluginInfo {
	return d.result
}
func (d *axisWeakPass) Check(URL string, meta plugin.TaskMeta) bool {
	userList := []string{
		"axis", "admin", "root",
	}
	succFlagList := []string{
		"Administration Page</title>", "System Components", "axis2-admin/upload",
		`include page="footer.inc">`, "axis2-admin/logout",
	}
	for _, user := range userList {
		for _, pass := range PassList {
			pass = strings.Replace(pass, "{user}", user, -1)
			PostStr := fmt.Sprintf("userName=%s&password=%s&submit=+Login+", user, pass)
			request, err := http.NewRequest("GET", URL+"/axis2/axis2-admin/login", strings.NewReader(PostStr))
			resp, err := common.RequestDo(request, true)
			if err != nil {
				return false
			}
			if resp.Other.StatusCode == 404 {
				return false
			}
			if resp.Other.StatusCode == 200 {
				for _, flag := range succFlagList {
					if strings.Contains(resp.ResponseRaw, flag) {
						result := d.info
						result.Response = resp.ResponseRaw
						result.Request = resp.RequestRaw
						result.Remarks = fmt.Sprintf("弱口令：%s,%s,%s", user, pass, result.Remarks)
						d.result = append(d.result, result)
						return true
					}
				}
			}
		}
	}
	return false
}
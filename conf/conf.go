/**
 * global conf
 * ENV:
 *   CONF_FILE      --- 配置文件名
 *   TZ             --- 时区名称"Asia/Shanghai"
 *
 * JSON:
 *  {
      "listen-host": "",
      "listen-port": 7080,
      "services": [
         {
             "name": "echo_server",
             "workerNum": 5,
             "timeout": 0,
             "wx-params": {
                 "token": "hello_rosbit",
                 "app-id": "",
                 "app-secret": "",
                 "aes-key":  null
             },
             "listen-endpoints": {
                 "service-path": "/wx",
                 "redirect-path": "/redirect"
             },
             "msg-proxy-pass": "http://yourhost.or.ip.here",
             "menu-handler": "http://yourhost.or.ip/handler/path"
         }
      ]
   }
 *
 * Rosbit Xu
 */
package gwconf

import (
	"fmt"
	"os"
	"time"
	"io/ioutil"
	"encoding/json"
)

type WxParamsConf struct {
	Token     string `json:"token"`
	AppId     string `json:"app-id"`
	AppSecret string `json:"app-secret"`
	AesKey    string `json:"aes-key"`
}

type EndpointConf struct {
	ServicePath  string `json:"service-path"`
	RedirectPath string `json:"redirect-path"`
}

type WxService struct {
	Name         string
	WorkerNum    int          `json:"workerNum"`
	Timeout      int
	WxParams     WxParamsConf `json:"wx-params"`
	Endpoints    EndpointConf `json:"listen-endpoints"`
	MsgProxyPass string       `json:"msg-proxy-pass"`
	MenuHandler  string       `json:"menu-handler"`
}

type WxServiceConf struct {
	ListenHost     string `json:"listen-host"`
	ListenPort     int    `json:"listen-port"`
	Services       []WxService
}

var (
	ServiceConf WxServiceConf
	Loc = time.FixedZone("UTC+8", 8*60*60)
)

func getEnv(name string, result *string, must bool) error {
	s := os.Getenv(name)
	if s == "" {
		if must {
			return fmt.Errorf("env \"%s\" not set", name)
		}
	}
	*result = s
	return nil
}

func CheckGlobalConf() error {
	var p string
	getEnv("TZ", &p, false)
	if p != "" {
		if loc, err := time.LoadLocation(p); err == nil {
			Loc = loc
		}
	}

	var confFile string
	if err := getEnv("CONF_FILE", &confFile, true); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(confFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &ServiceConf); err != nil {
		return err
	}

	return nil
}

func DumpConf() {
	fmt.Printf("conf: %v\n", ServiceConf)
	fmt.Printf("TZ time location: %v\n", Loc)
}

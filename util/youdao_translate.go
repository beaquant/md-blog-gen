package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"unicode"

	"github.com/golang/glog"

	"github.com/zwh8800/md-blog-gen/conf"
)

// 有道API: http://fanyi.youdao.com/openapi?path=data-mode
var translateCache *AutoClearKVCache

func init() {
	translateCache = NewAutoClearKVCache(100)
}

type youdaoResponse struct {
	Translation []string `json:"translation"`
	Query       string   `json:"query"`
	ErrorCode   int      `json:"errorCode"`
}

func isAscii(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func YoudaoTranslate(title string) string {
	if notename, ok := translateCache.Get(title); ok {
		return notename
	}

	if conf.Conf.Youdao.ApiUrl == "" {
		return ""
	}

	u, err := url.Parse(conf.Conf.Youdao.ApiUrl)
	if err != nil {
		glog.Errorln(err)
		return ""
	}
	q := u.Query()
	q.Set("q", title)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		glog.Errorln(err)
		return ""
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorln(err)
		return ""
	}
	var youdaoData youdaoResponse

	if err := json.Unmarshal(data, &youdaoData); err != nil {
		glog.Errorln(err)
		return ""
	}
	if youdaoData.ErrorCode != 0 {
		glog.Errorln("youdaoData.ErrorCode is ", youdaoData.ErrorCode)
		return ""
	}
	if len(youdaoData.Translation) == 0 {
		glog.Errorln("youdaoData.Translation length is 0")
		return ""
	}
	if !isAscii(youdaoData.Translation[0]) {
		glog.Errorln("youdaoData.Translation is not pure ascii:",
			youdaoData.Translation[0])
		return ""
	}

	translateCache.Put(title, youdaoData.Translation[0])
	return youdaoData.Translation[0]
}

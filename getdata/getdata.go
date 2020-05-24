package getdata

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Vary struct {
	Version string `json:"v"`
	Host    string `json:"host"`
	Path    string `json:"path"`
	TLS     string `json:"tls"`
	Ps      string `json:"ps"`
	Add     string `json:"add"` // url
	Prot    string `json:"port"`
	ID      string `json:"id"`
	Aid     string `json:"aid"`
	Net     string `json:"net"`
	Type    string `josn:"type"`
}

// ExampleScrape get telegarm v2list page data
func ExampleScrape(count string, cors bool) (string, bool) {
	// Request the HTML page.
	var c int
	var err error
	c, err = strconv.Atoi(count)
	var url string
	url = "https://t.me/s/V2List"
	if cors {
		url = strings.Join([]string{"https://cors.zme.ink", url}, "/")
	}
	// fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		return "bad", false
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		// log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return res.Status, false
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	root := doc.Find("body.widget_frame_base > main.tgme_main > div.tgme_container > section.tgme_channel_history > div.tgme_widget_message_wrap")
	length := root.Length()
	findthis := root.Eq(length - c).Find("div.tgme_widget_message_text").Text()
	return findthis, true
}

// MakeList use split to make a array for string
func MakeList(d string) []string {
	x := []string{}
	l := strings.Split(d, "vmess://")
	for i, item := range l {
		var l int
		l = len(item)
		if l > 0 {
			var strHaiCoder string
			var newstr string
			var v string
			var other bool

			other = strings.Contains(item, "?remarks=")
			if other {
				strsss := strings.Split(item, "?remarks=")
				var strtobyte []byte = []byte(strsss[0])
				decodeBytes := make([]byte, base64.StdEncoding.DecodedLen(len(strtobyte))) // 计算解码后的长度
				base64.StdEncoding.Decode(decodeBytes, strtobyte)

				newstr := string(decodeBytes[:])
				// fmt.Println(newstr)
				blen := len(newstr)
				a := strings.Index(newstr, ":")
				b := strings.Index(newstr, "@")
				c := strings.LastIndex(newstr, ":")
				uuid := newstr[a+1 : b]
				host := newstr[b+1 : c]
				port := newstr[c+1 : blen]

				if port == "4" || port == "44" {
					port = "443"
				}

				params := strsss[1]
				plen := len(params)
				e := strings.Index(params, "path=")
				f := strings.Index(params, "obfs=")
				g := strings.Index(params, "tls=")
				path := params[e+5 : f-1]
				obfs := params[f+5 : g-1]
				tls := params[g+4 : plen]

				if obfs == "websocket" {
					obfs = "ws"
				} else {
					obfs = "tcp"
				}

				if tls == "1" {
					tls = "tls"
				} else {
					tls = "tcp"
				}

				// fmt.Println(uuid, host, port, param, path, obfs, tls)
				// log.Println(uuid, host, port, param, path, obfs, tls)
				cumv := strconv.Itoa(i)
				myname := strings.Join([]string{"我的节点", cumv}, "-")
				vjson := &Vary{
					Version: "2",
					Host:    host,
					Path:    path,
					TLS:     tls,
					Ps:      myname,
					Add:     host,
					Prot:    port,
					ID:      uuid,
					Aid:     "1",
					Net:     obfs,
					Type:    "null",
				}
				bytes, err := json.Marshal(vjson)
				if err != nil {
					return x
				}
				// fmt.Println(string(bytes))
				v = strings.Join([]string{"vmess:", base64.StdEncoding.EncodeToString(bytes)}, "//")
				x = append(x, v)
			} else {
				decodeBytes, err := base64.StdEncoding.DecodeString(item)
				if err != nil {
					return x
				}
				strHaiCoder = `"ps" :"翻墙党fanqiangdang.com","" :`
				reg := regexp.MustCompile(strHaiCoder)
				newstr = reg.ReplaceAllString(string(decodeBytes), `"ps" :`)
				var strtobyte []byte = []byte(newstr)
				v = strings.Join([]string{"vmess:", base64.StdEncoding.EncodeToString(strtobyte)}, "//")
				x = append(x, v)
			}
		}
	}
	return x
}

// MakeData is a make Array to BASE64 string function
func MakeData(d []string) string {
	var data string = strings.Join(d[:], "\n")
	var strtobyte []byte = []byte(data)
	return base64.StdEncoding.EncodeToString(strtobyte)
}

// Start this
func Start(n string, w bool) string {
	var d []string
	var dd string = ""
	data, status := ExampleScrape(n, w)
	if status {
		d = MakeList(data)
		dd = MakeData(d)
	}
	return dd
}
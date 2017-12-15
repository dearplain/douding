package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Jpgs struct {
	Zoom string
	Page int
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("error: no url specified!")
		return
	}
	if strings.Contains(os.Args[1], "www.docin.com") {
		downDouding(os.Args[1])
	} else if strings.Contains(os.Args[1], "wenku.baidu.com") {
		downWenkuPPT(os.Args[1])
	}
}

func downWenkuPPT(url string) {

	// https://wkretype.bdimg.com/retype/zoom/7ff454ec0975f46527d3e167?pn=14&raww=1080
	// &rawh=810&o=jpg_6&md5sum=2854096ac319841e6a86ef6d8e931e7c&sign=2dccc30c24&
	// png=230016-246170&jpg=1394303-1518342&aimh=135&aimw=180
	// url := "https://wenku.baidu.com/view/bb52fecf561252d381eb6e65.html"
	buf := downloadA(url)

	urlRe := regexp.MustCompile(`(?U)view/(.*).html`)
	var ss []string
	ss = urlRe.FindStringSubmatch(url)
	urlID := ss[1]

	signRe := regexp.MustCompile(`(?U)(&md5sum=.*)"`)
	ss = signRe.FindStringSubmatch(string(buf))
	signStr := ss[1]

	bcsParam := regexp.MustCompile(`(?U)bcsParam', "([\s\S]*)"`)
	ss = bcsParam.FindStringSubmatch(string(buf))
	var err error
	ss[1] = strings.Replace(ss[1], "\\x22", "\x22", -1)
	var jpgs []Jpgs
	err = json.Unmarshal([]byte(ss[1]), &jpgs)
	fmt.Println(jpgs)

	os.MkdirAll(urlID, 0700)
	for i := range jpgs {
		imgURL := "https://wkretype.bdimg.com/retype/zoom/" + urlID + "?" + "pn=" +
			fmt.Sprintf("%d", jpgs[i].Page) + "&raww=1080&rawh=810&o=jpg_6" + signStr + jpgs[i].Zoom
		fmt.Println(imgURL)
		buf := downloadA(imgURL)
		if err != nil {
			fmt.Println(err)
			break
		}
		if len(buf) < 1024 {
			fmt.Println("download complete!")
			break
		}
		err = ioutil.WriteFile(urlID+"/"+fmt.Sprintf("%04d.jpg", jpgs[i].Page), buf, 0644)
		if err != nil {
			fmt.Println(err)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func downDouding(url string) {
	re := regexp.MustCompile("p-([0-9]*)")
	ss := re.FindStringSubmatch(url)

	id := ss[1]
	os.MkdirAll(id, 0700)
	for page := 1; ; page++ {
		buf, err := downloadBuf("http://211.147.220.164/index.jsp?file=" + id +
			"&width=800&pageno=" + fmt.Sprintf("%d", page))
		if err != nil {
			fmt.Println(err)
			break
		}
		if len(buf) < 1024 {
			fmt.Println("download complete!")
			break
		}
		err = ioutil.WriteFile(id+"/"+fmt.Sprintf("%04d.jpg", page), buf, 0644)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func downloadBuf(url string) ([]byte, error) {
	var buffer bytes.Buffer
	if err := download(&buffer, url); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

func download(w io.Writer, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	return err
}

func downloadA(url string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (iPad; CPU OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

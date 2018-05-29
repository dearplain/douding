package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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

	regID := regexp.MustCompile("view/(.*).html")
	urls := regID.FindAllStringSubmatch(url, -1)

	docID := urls[0][1]
	getUrl := "https://wenku.baidu.com/browse/getbcsurl?doc_id=" + docID + "&pn=1&rn=99999&type=ppt"

	body, _ := downloadBuf(getUrl)

	regPages := regexp.MustCompile(`{"zoom":"(.*?)","page"`)
	pages := regPages.FindAllStringSubmatch(string(body), -1)

	os.Mkdir(docID, 0700)
	for i, p := range pages {
		p[1] = strings.Replace(p[1], "\\", "", -1)
		fmt.Println("downloading:", p[1])
		body, _ := downloadBuf(p[1])
		ioutil.WriteFile(docID+"/img"+fmt.Sprint(i)+".jpg", body, 0644)
	}
	fmt.Println("download complete! total:", len(pages))
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

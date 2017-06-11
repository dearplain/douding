package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("error: no url specified!")
		return
	}

	re := regexp.MustCompile("p-([0-9]*)")
	ss := re.FindStringSubmatch(os.Args[1])
	if len(ss) <= 1 {
		fmt.Println("url incorect, expecting url like: http://www.docin.com/p-172303154.html")
		return
	}

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

package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[2]
	pt("%s\n", id)
	pagePath := fmt.Sprintf("http://www.vvic.com/api/item/%s", id)
	resp, err := http.Get(pagePath)
	ce(err, "get")
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		w.Write([]byte("invalid id"))
		return
	}
	var data struct {
		Code int
		Data struct {
			Imgs string // 图片
			Desc string // 描述html
		}
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	ce(err, "decode")

	buf := new(bytes.Buffer)
	archive := zip.NewWriter(buf)

	for i, imgPath := range strings.Split(data.Data.Imgs, ",") {
		if !strings.HasPrefix(imgPath, "http:") {
			imgPath = "http:" + imgPath
		}
		pt("%s\n", imgPath)

		header := new(zip.FileHeader)
		header.Name = fmt.Sprintf("%s-a-%04d%s", id, i,
			path.Ext(imgPath))
		header.Method = zip.Deflate
		writer, err := archive.CreateHeader(header)
		ce(err, "CreateHeader")
		resp, err := http.Get(imgPath)
		ce(err, "get image")
		defer resp.Body.Close()
		io.Copy(writer, resp.Body)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data.Data.Desc))
	ce(err, "goquery doc")
	doc.Find("img").Each(func(i int, se *goquery.Selection) {
		imgSrc, _ := se.Attr("src")
		pt("%s\n", imgSrc)

		header := new(zip.FileHeader)
		header.Name = fmt.Sprintf("%s-b-%04d%s", id, i,
			path.Ext(imgSrc))
		header.Method = zip.Deflate
		writer, err := archive.CreateHeader(header)
		ce(err, "CreateHeader")
		resp, err := http.Get(imgSrc)
		ce(err, "get image")
		defer resp.Body.Close()
		io.Copy(writer, resp.Body)
	})

	archive.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+id+".zip")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	http.ServeContent(w, r, id+".zip", time.Now(), bytes.NewReader(buf.Bytes()))

}

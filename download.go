package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/reusee/vviccommon"
)

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[2]
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
		if !strings.HasPrefix(imgPath, "http") {
			continue
		}

		resp, err := http.Get(imgPath)
		ce(err, "get image")
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		ce(err, "read body")

		header := new(zip.FileHeader)
		header.Name = fmt.Sprintf("a-%s-%04d%s", id, i,
			path.Ext(imgPath))
		header.Method = zip.Deflate
		writer, err := archive.CreateHeader(header)
		ce(err, "CreateHeader")
		err = vviccommon.ScaleTo800x800(bytes.NewReader(body), writer)
		ce(err, "scale to 800x800")

		header = new(zip.FileHeader)
		header.Name = fmt.Sprintf("b-%s-%04d%s", id, i,
			path.Ext(imgPath))
		header.Method = zip.Deflate
		writer, err = archive.CreateHeader(header)
		ce(err, "CreateHeader")
		err = vviccommon.CompositeLogo(bytes.NewReader(body), writer)
		ce(err, "composite logo")

	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data.Data.Desc))
	ce(err, "goquery doc")
	doc.Find("img").Each(func(i int, se *goquery.Selection) {
		imgSrc, _ := se.Attr("src")
		pt("%s\n", imgSrc)
		if !strings.HasPrefix(imgSrc, "http") {
			return
		}

		resp, err := http.Get(imgSrc)
		ce(err, "get image")
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		ce(err, "read body")

		// origin image
		header := new(zip.FileHeader)
		header.Name = fmt.Sprintf("c-%s-%04d%s", id, i,
			path.Ext(imgSrc))
		header.Method = zip.Deflate
		writer, err := archive.CreateHeader(header)
		ce(err, "CreateHeader")
		_, err = writer.Write(body)
		ce(err, "write")

		// scaled image
		header = new(zip.FileHeader)
		header.Name = fmt.Sprintf("m-%s-%04d%s", id, i,
			path.Ext(imgSrc))
		header.Method = zip.Deflate
		writer, err = archive.CreateHeader(header)
		ce(err, "CreateHeader")
		ce(vviccommon.ScaleForMobile(bytes.NewReader(body), writer), "scale")

		// 770
		header = new(zip.FileHeader)
		header.Name = fmt.Sprintf("w-770px-%s-%04d.jpg", id, i)
		header.Method = zip.Deflate
		writer, err = archive.CreateHeader(header)
		ce(err, "CreateHeader")
		ce(vviccommon.ScaleImageToJpeg(770, 90, bytes.NewReader(body), writer), "scale to 770")

	})

	archive.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+id+".zip")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	http.ServeContent(w, r, id+".zip", time.Now(), bytes.NewReader(buf.Bytes()))

}

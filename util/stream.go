package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/google/gxui/math"
	"io"
	"io/ioutil"
	"mime/multipart"
)

type BufferStream struct {
	bytes  []byte
	cursor int
}

func NewBufferStream() *BufferStream {
	return &BufferStream{bytes: make([]byte, 0)}
}
func (s *BufferStream) Write(b []byte) (int, error) {
	s.bytes = append(s.bytes, b...)
	return len(b), nil
}

func (s *BufferStream) Read(b []byte) (int, error) {
	if s.cursor >= len(s.bytes) {
		return 0, io.EOF
	}
	n := math.Min(len(s.bytes)-s.cursor, len(b))

	for i := 0; i < n; i++ {
		b[i] = s.bytes[s.cursor+i]
	}
	s.cursor += n
	return n, nil
}
func (s *BufferStream) ReadAll() []byte {
	return s.bytes
}
func Read() {
	fmt.Print("read1111")
}
func Upload(url, name string, reader io.Reader) ([]byte, error) {

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//创建第一个需要上传的文件,filepath.Base获取文件的名称
	fileWriter, _ := bodyWriter.CreateFormFile("file", name)
	io.Copy(fileWriter, reader)

	//获取请求Content-Type类型,后面有用
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	res, e := httplib.Post(url).Header("Content-Type", contentType).Body(bodyBuf.Bytes()).DoRequest()
	if e != nil {
		return nil, e
	}
	if res.StatusCode == 200 {
		return ioutil.ReadAll(res.Body)
	} else {
		b, e := ioutil.ReadAll(res.Body)
		if e != nil {
			return nil, e
		}
		return nil, errors.New(string(b))
	}
}
func Download(url string) (io.Reader, error) {
	res, e := httplib.Get(url).DoRequest()
	if e != nil {

		return nil, e
	}
	if res.StatusCode != 200 {
		b, e := ioutil.ReadAll(res.Body)
		if e != nil {
			return nil, e
		}
		return nil, errors.New(string(b))
	}
	return res.Body, nil
}

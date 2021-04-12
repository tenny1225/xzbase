package util

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func CopyInternetFile(url, filepath string) error {
	if e := os.MkdirAll(path.Dir(filepath), os.ModePerm); e != nil {
		fmt.Println(e)
	}

	f, e := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0777)

	if e != nil {
		return e
	}
	defer f.Close()
	res, e := httplib.Get(url).DoRequest()
	if e != nil {
		return e
	}
	if res.StatusCode != 200 {
		b, e := ioutil.ReadAll(res.Body)
		if e != nil {
			return e
		}
		return errors.New(string(b))
	}
	_, e = io.Copy(f, res.Body)
	return e

}
func CopyInternetFileAndClip(url, filepath string,clip[]int) error {
	if e := os.MkdirAll(path.Dir(filepath), os.ModePerm); e != nil {
		fmt.Println(e)
	}

	f, e := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0777)

	if e != nil {
		return e
	}
	defer f.Close()
	res, e := httplib.Get(url).DoRequest()
	if e != nil {
		return e
	}
	if res.StatusCode != 200 {
		b, e := ioutil.ReadAll(res.Body)
		if e != nil {
			return e
		}
		return errors.New(string(b))
	}
	if clip!=nil&&len(clip)==4{
		e = Clip(res.Body,f,clip[0],clip[1],clip[2],clip[3],100)
	}else{
		_, e = io.Copy(f, res.Body)
	}


	return e

}
func Copy(dist, src string) error {
	os.MkdirAll(path.Dir(dist), 0666)
	distFile, e := os.OpenFile(dist, os.O_CREATE|os.O_RDWR, 0666)
	if e != nil {
		return e
	}
	defer distFile.Close()

	srcFile, e := os.Open(src)
	if e != nil {
		return e
	}
	defer srcFile.Close()

	_, e = io.Copy(distFile, srcFile)
	return e

}
func CopyFolder(dist, src string) error {
	list, e := ioutil.ReadDir(src)
	if e != nil {
		return e
	}
	os.RemoveAll(dist)
	os.MkdirAll(dist, 0666)
	for _, l := range list {
		if l.IsDir() {
			fp := path.Join(src, l.Name())
			distPath := path.Join(dist, l.Name())
			if e := CopyFolder(distPath, fp); e != nil {
				os.RemoveAll(dist)
				return e
			}
		} else {
			fp := path.Join(src, l.Name())
			distPath := path.Join(dist, l.Name())

			if e := Copy(distPath, fp); e != nil {
				os.RemoveAll(dist)
				return e
			}
		}

	}
	return nil
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func ReadAll(src string) ([]byte, error) {
	f, e := os.Open(src)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
func WriteAll(src string, b []byte) error {
	f, e := os.OpenFile(src, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if e != nil {
		return e
	}
	defer f.Close()
	_, e = f.Write(b)
	return e
}
func Clip(in io.Reader, out io.Writer, x0, y0, x1, y1, quality int) error {
	origin, fm, err := image.Decode(in)
	if err != nil {
		return err
	}

	switch fm {
	case "jpeg":
		img := origin.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		return jpeg.Encode(out, subImg, &jpeg.Options{quality})
	case "png":
		switch origin.(type) {
		case *image.NRGBA:
			img := origin.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			return png.Encode(out, subImg)
		case *image.RGBA:
			img := origin.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			return png.Encode(out, subImg)
		}
	case "gif":
		img := origin.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		return gif.Encode(out, subImg, &gif.Options{})
	case "bmp":
		img := origin.(*image.RGBA)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
		return bmp.Encode(out, subImg)
	default:
		return errors.New("ERROR FORMAT")
	}
	return nil
}
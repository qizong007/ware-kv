package util

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
)

func ZipBytes(origin []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	_, _ = w.Write(origin)
	_ = w.Flush()
	return b.Bytes()
}

func UnzipBytes(zipContent []byte) (originInfo []byte) {
	var b bytes.Buffer
	b.Write(zipContent)
	r, err := gzip.NewReader(&b)
	if err != nil {
		fmt.Println("gzip.NewReader Failed: ", err)
		return originInfo
	}
	defer r.Close()

	originInfo, err = ioutil.ReadAll(r)
	if err != nil {
		fmt.Println("ioutil.ReadAll Failed: ", err)
	}

	return originInfo
}

package main

import (
	"compress/gzip"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	fPath := flag.String("path", "", "path to checksum")
	flag.Parse()
	sig, err := sha1sum(*fPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sig)
}

// $ gunzip -c _extra/http.log.gz | sha1sum
func sha1sum(fileName string) (string, error) {
	// try opening the file named ‹fileName›
	var reader io.Reader
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	reader = f
	// if the fileName ~= '*.gz' wrappit in a gunzipping reader
	if strings.HasSuffix(f.Name(), ".gz") {
		gzReader, err = gzip.NewReader(f)
		if err != nil {
			return "", err
		}
		defer gzReader.Close()
		reader = gzReader
	}

	w := sha1.New()
	if _, err := io.Copy(w, reader); err != nil {
		return "", err
	}
	sig := w.Sum(nil)
	return fmt.Sprintf("%x", sig), nil
}

/*
Write a function that gets an index file with names of files and sha256
signatures in the following format
0c4ccc63a912bbd6d45174251415c089522e5c0e75286794ab1f86cb8e2561fd  taxi-01.csv
f427b5880e9164ec1e6cda53aa4b2d1f1e470da973e5b51748c806ea5c57cbdf  taxi-02.csv
4e251e9e98c5cb7be8b34adfcb46cc806a4ef5ec8c95ba9aac5ff81449fc630c  taxi-03.csv
...

You should compute concurrently sha256 signatures of these files and see if
they math the ones in the index file.

  - Print the number of processed files
  - If there's a mismatch, print the offending file(s) and exit the program with
    non-zero value

Grab taxi-sha256.zip from the web site and open it. The index file is sha256sum.txt
*/
package main

import (
	"bufio"
	"compress/bzip2"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func fileSig(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, bzip2.NewReader(file))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Parse signature file. Return map of path->signature
func parseSigFile(r io.Reader) (map[string]string, error) {
	sigs := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// Line example
		// 6c6427da7893932731901035edbb9214  nasa-00.log
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			// TODO: line number
			return nil, fmt.Errorf("bad line: %q", scanner.Text())
		}
		sigs[fields[1]] = fields[0]
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sigs, nil
}

func TimeFunc(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

func main() {
	// Change to where to unzipped taxi-sha256.zip
	rootDir, found := os.LookupEnv("TAXITEMP")
	if !found {
		log.Fatal("Set the location of the unzipped files in TAXITEMP. Eg. `set TAXITEMP (mktemp -d) && unzip ../_extras/taxi-sha256.zip -d $TAXITEMP`")
	}
	file, err := os.Open(path.Join(rootDir, "sha256sum.txt"))
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer file.Close()

	sigs, err := parseSigFile(file)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	computedSignatures := make(chan result)
	start := time.Now()
	for name, signature := range sigs {
		fileName := path.Join(rootDir, name) + ".bz2"
		go sigWorker(computedSignatures, fileName, signature)
	}

	allSignaturesMatched := true
	countdown := len(sigs)
	for res := range computedSignatures {
		allSignaturesMatched = res.err != nil && res.matched
		if res.err != nil {
			fmt.Printf("💀 error processing file %q: %v\n", res.filepath, res.err)
		} else if !res.matched {
			fmt.Printf("🛑 %q : file hash mismatch\n", res.filepath)
		} else {
			fmt.Printf("✅ %q\n", res.filepath)
		}
		if countdown -= 1; countdown == 0 {
			close(computedSignatures)
		}
	}
	duration := time.Since(start)
	fmt.Printf("processed %d files in %v\n", len(sigs), duration)
	if !allSignaturesMatched {
		os.Exit(1)
	}
}
func sigWorker(results chan<- result, filePath string, expectedSig string) {
	actualSig, err := fileSig(filePath)
	res := result{
		filepath: filePath,
	}
	if err != nil {
		res.err = fmt.Errorf("error: %s - %s", filePath, err)
	} else {
		res.matched = (actualSig == expectedSig)
	}
	results <- res
}

type result struct {
	filepath string
	matched  bool
	err      error
}

package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"io/ioutil"
	"os"
)

var (
	EOC = errors.New("end of collection")
)

type splitter struct {
	index int
	paths []string
}

func newSplitter() *splitter {
	return &splitter{}
}

func (sp *splitter) split(csv string, rows int) error {
	f, err := os.Open(csv)
	if err != nil {
		return err
	}
	defer f.Close()

	paths := make([]string, 0)
	data := make([]byte, 0)
	buf := bytes.NewBuffer(data)
	count := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		buf.Write(sc.Bytes())
		buf.Write([]byte("\n"))
		count++
		if count == rows {
			path, err := saveFile(buf.Bytes())
			if err != nil {
				return err
			}

			paths = append(paths, path)
			count = 0
			buf.Reset()
		}
	}

	if len(data) > 0 {
		path, err := saveFile(data)
		if err != nil {
			return err
		}
		paths = append(paths, path)
	}

	sp.paths = paths
	return nil
}

func saveFile(data []byte) (string, error) {
	f, err := ioutil.TempFile("", "paack.*.csv")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func (sp *splitter) next() ([][]string, error) {
	if sp.index == len(sp.paths) {
		return nil, EOC
	}

	path := sp.paths[sp.index]
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rd := csv.NewReader(f)
	data, err := rd.ReadAll()
	if err != nil {
		return nil, err
	}

	sp.index++
	return data, nil
}

func (sp *splitter) close() error {
	for _, f := range sp.paths {
		err := os.Remove(f)
		if err != nil {
			return err
		}
	}

	return nil
}

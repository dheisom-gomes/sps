package main

import (
	"net"
	"os"
	"io"
)

func ReadLinesFromBytes(d []byte) [][]byte {
	result := [][]byte{}
	line := []byte{}
	for _, b := range d {
		if b == '\n' {
			result = append(result, line)
			line = []byte{}
			continue
		} else if b == '\r' {
			continue
		}
		line = append(line, b)
	}
	if len(line) > 0 {
		result = append(result, line)
	}
	return result
}

func ReadLineFromConnection(c net.Conn) (string, error) {
	line := []byte{}
	for {
		b := make([]byte, 1)
		s, err := c.Read(b)
		if err != nil {
			return "", err
		} else if b[0] == '\r' || s == 0 {
			continue
		} else if b[0] == '\n' {
			break
		}
		line = append(line, b[0])
	}
	return string(line), nil
}

func AsyncReceiver(c net.Conn, bs uint) (chan []byte, chan error) {
	data := make(chan []byte)
	err := make(chan error)
	go func() {
		received := make([]byte, bs)
		s, e := c.Read(received)
		if e != nil {
			err <- e
		} else {
			data <- received[:s]
		}
	}()
	return data, err
}


// Copied from go1.17/src/os/file.go to run on go 1.13(termux legacy)
func ReadFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var size int
	if info, err := f.Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}
	size++ // one byte for final read at EOF

	// If a file claims a small size, read at least 512 bytes.
	// In particular, files in Linux's /proc claim size 0 but
	// then do not work right if read in small pieces,
	// so an initial read of 1 byte would not work correctly.
	if size < 512 {
		size = 512
	}

	data := make([]byte, 0, size)
	for {
		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
		n, err := f.Read(data[len(data):cap(data)])
		data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return data, err
		}
	}
}

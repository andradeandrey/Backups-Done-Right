package bdrupload

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"os"
	"time"
)

type Upchan_t struct {
	Rowid int64
	Path  string
}

type Downchan_t struct {
	Rowid int
	Err   error
}

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

const bufferSize = 524288

func Server(upchan chan *Upchan_t, done chan bool) {
	var count int
	var size int64
	readBuffer := make([]byte, bufferSize)
	ciphertext := make([]byte, bufferSize)
	key_text := "32o4908go293hohg98fh40ghaidlkjk1"
	c, err := aes.NewCipher([]byte(key_text))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
		os.Exit(-1)
	}
	for f := range upchan {
		fmt.Printf("Server: received rowID=%d path=%s\n", f.Rowid, f.Path)
		//      fmt.Printf("%T %#v\n",f,f)
		size = 0

		// open file and create a reader
		file, err := os.Open(f.Path)
		reader := bufio.NewReader(file)

		// for this file create a cipher and new sha256 state
		cfb := cipher.NewCFBEncrypter(c, commonIV)
		h := sha256.New() // h is a hash.Hash
		// time how long to read, encrypt, and checksum a file
		t0 := time.Now().UnixNano()
		for {
			if count, err = reader.Read(readBuffer); err != nil {
				break
			}
			size = size + int64(count)
			cfb.XORKeyStream(ciphertext[:count], readBuffer[:count])
			h.Write(ciphertext[:count])
		}
		t1 := time.Now().UnixNano()
		file.Close()
		seconds := float64(t1-t0) / 1000000000
		fmt.Printf("%x %s %4.2f MB/sec\n", h.Sum(nil), f.Path, float64(size)/(1024*1024*seconds))

	}
	fmt.Print("Server: Channel closed, existing\n")
	done <- true
}

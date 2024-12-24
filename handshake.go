package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type conn struct {
	net.Conn
}

// Handshake open a handshake for websocket
// rfc6455, 1.3
func Handshake(w http.ResponseWriter, r *http.Request) (*conn, error) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("wrong http method"))
		return nil, errors.New(" Wrong http method")
	}
	key := r.Header.Get("Sec-WebSocket-Key")
	key = strings.Trim(key, " ")
	guid := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	shaHash := sha1.New()
	shaHash.Write([]byte(key))
	shaHash.Write([]byte(guid))
	baseKey := base64.StdEncoding.EncodeToString(shaHash.Sum(nil))
	netConn, brw, err := http.NewResponseController(w).Hijack()
	if err != nil {
		return nil, err
	}
	buf := brw.Writer.AvailableBuffer()
	buf = append(buf, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
	buf = append(buf, baseKey...)
	buf = append(buf, "\r\n"...)
	buf = append(buf, "\r\n"...)
	_, err = netConn.Write(buf)
	if err != nil {
		return nil, err
	}
	if err := netConn.SetDeadline(time.Time{}); err != nil {
		return nil, err
	}
	return &conn{netConn}, nil
}

const (
	defaultBufSize = 1024 * 4
)

func (c *conn) ReadMsg() (string, error) {
	buf := make([]byte, defaultBufSize)
	_, err := c.Read(buf)
	if err != nil {
		panic(err)
	}

	_1b := buf[0]
	_2b := buf[1]
	fin := (_1b & 128) == 128
	_ = fin
	plen := _2b & 127
	plength := uint64(plen) + 1
	mk := buf[2:6]
	payload := buf[6:plength]
	if plen == 126 {
		l := binary.BigEndian.Uint64(append([]byte{0, 0, 0, 0, 0, 0}, buf[2:4]...))
		plength = l
		mk = buf[4:8]
		payload = buf[8:]
		if plength > 4088 {
			for uint64(len(payload)) < plength {
				buf2 := make([]byte, 4096)
				_, err = c.Read(buf2)
				if err != nil {
					panic(err)
				}
				payload = append(payload, buf2...)
				fmt.Println("done", len(payload), plength)
			}
			if plength < uint64(len(payload)) {
				payload = payload[:plength]
			}
		} else {
			payload = buf[8:plength]
		}
	} else if plen == 126 {
		l := binary.BigEndian.Uint64(append([]byte{0, 0, 0, 0, 0, 0}, buf[2:9]...))
		plength = l
		mk = buf[9:13]
		payload = buf[13:plength]
	}
	_ = plength
	ump := []byte{}
	for i, p := range payload {
		ump = append(ump, p^mk[i%4])
	}

	return string(ump), nil
}

func (c *conn) Run(f func(msg string)) error {
	for c != nil {
		msg, err := c.ReadMsg()
		if err != nil {
			return err
		}
		f(msg)
	}
	return nil
}

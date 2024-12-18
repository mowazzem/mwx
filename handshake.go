package main

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// Handshake open a handshake for websocket
// rfc6455, 1.3
func Handshake(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
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
		fmt.Println(err)
		return nil, err
	}
	buf := brw.Writer.AvailableBuffer()
	buf = append(buf, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
	buf = append(buf, baseKey...)
	buf = append(buf, "\r\n"...)
	buf = append(buf, "\r\n"...)
	_, err = netConn.Write(buf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err := netConn.SetDeadline(time.Time{}); err != nil {
		return nil, err
	}
	return netConn, nil
}



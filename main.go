package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
	"log"
)

func wsocket(w http.ResponseWriter, r *http.Request) {
	wskey := r.Header.Get("Sec-WebSocket-Key")
	wskey = strings.Trim(wskey, " ")
	fmt.Println(wskey)

	guid := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

	shaHash := sha1.New()
	shaHash.Write([]byte(wskey))
	shaHash.Write([]byte(guid))

	baseKey := base64.StdEncoding.EncodeToString(shaHash.Sum(nil))
	fmt.Println(baseKey)

	netConn, brw, err := http.NewResponseController(w).Hijack()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := netConn.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	buf := brw.Writer.AvailableBuffer()
	buf = append(buf, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
	buf = append(buf, baseKey...)
	buf = append(buf, "\r\n"...)
	buf = append(buf, "\r\n"...)

	_, err = netConn.Write(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := netConn.SetDeadline(time.Time{}); err != nil {
		panic(err)
	}

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsocket)
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()
}

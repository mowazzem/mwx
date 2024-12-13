package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
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

	netConn, _, err := http.NewResponseController(w).Hijack()
	if err != nil {
		fmt.Println(err)
		return
	}

	//buf := brw.Writer.AvailableBuffer()
	p := []byte{}
	p = append(p, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
	p = append(p, baseKey...)
	p = append(p, "\r\n"...)
	p = append(p, "\r\n"...)

	_, err = netConn.Write(p)
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

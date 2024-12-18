package main

import (
	"fmt"
	"net/http"
)

func wsocket(w http.ResponseWriter, r *http.Request) {
	c, err := Handshake(w, r)
	if err != nil {
		panic(err)
	}
	var mk []byte
	once:=true
	for {
		buf := make([]byte, 10)
		_, err := c.Read(buf)
		if err != nil {
			panic(err)
		}
		firstByte := buf[0]
		fmt.Printf("%08b\n",firstByte)
		fmt.Printf("%08b\n", buf[1])
		pay:=buf
		if once{
			mk = buf[2:6]
			pay=buf[6:]
			once=false
		}
		um := []byte{}


		for i, b := range pay {
			if b == 0 {
				break
			}
			ub := b ^ mk[i%4]
			um = append(um, ub)
		}
		once=true
		fmt.Println(string(um))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsocket)
	s := &http.Server{
		Addr:    ":8089",
		Handler: mux,
	}
	s.ListenAndServe()
}

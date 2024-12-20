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

	err = c.Run(func(msg string) {
		fmt.Println(msg)
	})
	if err != nil {
		fmt.Println(err)
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

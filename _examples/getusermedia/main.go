package main

import (
	"fmt"
	"log"

	"github.com/nobonobo/webrtc"
)

func main() {
	constraints := webrtc.NewConstraints(true, true)
	stream, err := webrtc.GetUserMedia(constraints)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(stream)
}

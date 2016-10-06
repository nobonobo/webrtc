# webrtc

isomorphic package for WebRTC

- WebRTC wrapper for native(github.com/keroserene/go-webrtc)
- WebRTC wrapper for GopherJS

## dependencies for native

- brew install pkg-config
- go get -u github.com/keroserene/go-webrtc

## dependencies for gopherjs

- go get -u github.com/gopherjs/gopherjs

## install

```sh
go get -u github.com/nobonobo/webrtc
```

## usage

getUserMedia sample(gopherjs only)
```go
package main

import "github.com/nobonobo/webrtc"

func main() {
	stream, err := webrtc.GetUserMedia(webrtc.NewConstraints(true, true))
	if err != nil {
		log.Println(err)
		return
	}
    ...
}
```

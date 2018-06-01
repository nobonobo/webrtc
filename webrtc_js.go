// +build js

package webrtc

import (
	"fmt"
	"sync"

	"github.com/gopherjs/gopherjs/js"
)

var (
	navigator      *js.Object
	peerConnection *js.Object
)

func init() {
	navigator = js.Global.Get("navigator")
	peerConnection = js.Global.Get("RTCPeerConnection")
}

// PeerConnection ...
type PeerConnection struct {
	pc *js.Object
}

// NewPeerConnection ...
func NewPeerConnection(config *Configuration) (pc *PeerConnection, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("create peer connection: %s", r)
			pc = nil
		}
	}()
	jpc := peerConnection.New(config.o)
	if jpc == nil {
		return nil, fmt.Errorf("create peer connection: failed")
	}
	pc = &PeerConnection{pc: jpc}
	return
}

// OnNegotiationNeeded ...
func (pc *PeerConnection) OnNegotiationNeeded(cb func()) {
	pc.pc.Call("addEventListener", "negotiationneeded",
		func(arg *js.Object) {
			cb()
		}, false,
	)
}

// OnIceCandidate ...
func (pc *PeerConnection) OnIceCandidate(cb func(*IceCandidate)) {
	pc.pc.Call("addEventListener", "icecandidate",
		func(ev *js.Object) {
			candidate := ev.Get("candidate")
			if candidate != nil {
				cb(NewIceCandidateFromObj(candidate))
			} else {
				ev := js.Global.Get("CustomEvent").New("icegatheringstatechange", js.M{"eventPhase": 2})
				pc.pc.Call("dispatchEvent", ev)
			}
		}, false,
	)
}

// OnIceCandidateError ...
func (pc *PeerConnection) OnIceCandidateError(cb func()) {
	pc.pc.Call("addEventListener", "icecandidateerror",
		func(arg *js.Object) {
			cb()
		}, false,
	)
}

// OnSignalingStateChange ...
func (pc *PeerConnection) OnSignalingStateChange(cb func(SignalingState string)) {
	pc.pc.Call("addEventListener", "signalingstatechange",
		func(ev *js.Object) {
			cb([]string{
				"Stable",
				"HaveLocalOffer", "HaveLocalPrAnswer",
				"HaveRemoteOffer", "HaveRemotePrAnswer",
				"Closed",
			}[ev.Get("eventPhase").Int()])
		}, false,
	)
}

// OnIceConnectionStateChange ...
func (pc *PeerConnection) OnIceConnectionStateChange(cb func(IceConnectionState string)) {
	pc.pc.Call("addEventListener", "iceconnectionstatechange",
		func(ev *js.Object) {
			cb([]string{
				"New", "Checking", "Connected", "Completed",
				"Failed", "Disconnected", "Closed",
			}[ev.Get("eventPhase").Int()])
		}, false,
	)
}

// OnIceGatheringStateChange ...
func (pc *PeerConnection) OnIceGatheringStateChange(cb func(IceGatheringState string)) {
	pc.pc.Call("addEventListener", "icegatheringstatechange",
		func(ev *js.Object) {
			cb([]string{
				"New", "Gathering", "Complete",
			}[ev.Get("eventPhase").Int()])
		}, false,
	)
}

// OnConnectionStateChange ...
func (pc *PeerConnection) OnConnectionStateChange(cb func(PeerConnectionState string)) {
	pc.pc.Call("addEventListener", "connectionstatechange",
		func(ev *js.Object) {
			cb([]string{
				"New", "Connecting", "Connected", "Disconnected", "Failed",
			}[ev.Get("eventPhase").Int()])
		}, false,
	)
}

// OnDataChannel ...
func (pc *PeerConnection) OnDataChannel(cb func(*DataChannel)) {
	pc.pc.Call("addEventListener", "datachannel",
		func(dc *js.Object) {
			cb(&DataChannel{dc: dc})
		}, false,
	)
}

// OnAddStream ...
func (pc *PeerConnection) OnAddStream(cb func(*MediaStream)) {
	pc.pc.Call("addEventListener", "addstream",
		func(ev *js.Object) {
			stream := ev.Get("stream")
			cb(&MediaStream{o: stream})
		}, false,
	)
}

// OnRemoveStream ...
func (pc *PeerConnection) OnRemoveStream(cb func(*MediaStream)) {
	pc.pc.Call("addEventListener", "removestream",
		func(ev *js.Object) {
			stream := ev.Get("stream")
			cb(&MediaStream{o: stream})
		}, false,
	)
}

// AddStream ...
func (pc *PeerConnection) AddStream(stream *MediaStream) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	pc.pc.Call("addStream", stream.o)
	return
}

// AddIceCandidate ...
func (pc *PeerConnection) AddIceCandidate(ic *IceCandidate) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	pc.pc.Call("addIceCandidate", ic.o)
	return
}

// Close ...
func (pc *PeerConnection) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	pc.pc.Call("close")
	return
}

// ConnectionState ...
func (pc *PeerConnection) ConnectionState() string {
	return pc.pc.Get("connectionState").String()
}

// CreateAnswer ...
func (pc *PeerConnection) CreateAnswer() (s *SessionDescription, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	pc.pc.Call(
		"createAnswer",
		func(desc *js.Object) {
			s = NewSessionDescriptionFromObj(desc)
			wg.Done()
		},
		func(e *js.Object) {
			err = fmt.Errorf("create answer failed: %s", e)
			wg.Done()
		},
		nil,
	)
	wg.Wait()
	return
}

// CreateDataChannel ...
func (pc *PeerConnection) CreateDataChannel(label string) (dc *DataChannel, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	jdc := pc.pc.Call("createDataChannel", label)
	if jdc == nil {
		return nil, fmt.Errorf("create data channel: failed")
	}
	dc = &DataChannel{dc: jdc}
	return
}

// CreateOffer ...
func (pc *PeerConnection) CreateOffer() (s *SessionDescription, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	pc.pc.Call(
		"createOffer",
		func(desc *js.Object) {
			s = NewSessionDescriptionFromObj(desc)
			wg.Done()
		},
		func(e *js.Object) {
			err = fmt.Errorf("create offer failed: %s", e)
			wg.Done()
		},
		nil,
	)
	wg.Wait()
	return
}

// IceGatheringState ...
func (pc *PeerConnection) IceGatheringState() string {
	return pc.pc.Get("iceGatheringState").String()
}

// LocalDescription ...
func (pc *PeerConnection) LocalDescription() (sdp *SessionDescription) {
	sd := pc.pc.Get("localDescription")
	return NewSessionDescriptionFromObj(sd)
}

// RemoteDescription ...
func (pc *PeerConnection) RemoteDescription() (sdp *SessionDescription) {
	sd := pc.pc.Get("remoteDescription")
	return NewSessionDescriptionFromObj(sd)
}

// SetLocalDescription ...
func (pc *PeerConnection) SetLocalDescription(sdp *SessionDescription) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	pc.pc.Call("setLocalDescription", sdp.o)
	return
}

// SetRemoteDescription ...
func (pc *PeerConnection) SetRemoteDescription(sdp *SessionDescription) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	pc.pc.Call("setRemoteDescription", sdp.o)
	return
}

// SignalingState ...
func (pc *PeerConnection) SignalingState() string {
	return pc.pc.Get("signalingState").String()
}

// DataChannel ...
type DataChannel struct {
	dc *js.Object
}

// OnOpen ...
func (c *DataChannel) OnOpen(cb func()) {
	c.dc.Call("addEventListener", "open",
		func(ev *js.Object) {
			cb()
		}, false,
	)
}

// OnClose ...
func (c *DataChannel) OnClose(cb func()) {
	c.dc.Call("addEventListener", "close",
		func(ev *js.Object) {
			cb()
		}, false,
	)
}

// OnMessage ...
func (c *DataChannel) OnMessage(cb func([]byte)) {
	c.dc.Call("addEventListener", "message",
		func(ev *js.Object) {
			cb(ev.Get("data").Interface().([]byte))
		}, false,
	)
}

// Close ...
func (c *DataChannel) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	c.dc.Call("close")
	return
}

// ID ...
func (c *DataChannel) ID() int {
	return c.dc.Get("id").Int()
}

// ReadyState ...
func (c *DataChannel) ReadyState() string {
	return c.dc.Get("readyState").String()
}

// Send ...
func (c *DataChannel) Send(data []byte) {
	c.dc.Call("send", data)
}

// Label ...
func (c *DataChannel) Label() string {
	return c.dc.Get("label").String()
}

// MediaStream ...
type MediaStream struct {
	o *js.Object
}

// GetUserMedia ...
func GetUserMedia(constraints *Constraints) (stream *MediaStream, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		mediaDevices := navigator.Get("mediaDevices")
		p := mediaDevices.Call("getUserMedia", constraints)
		p.Call("then", func(ev *js.Object) {
			stream = &MediaStream{o: ev}
			wg.Done()
		}).Call("catch", func(e *js.Object) {
			err = fmt.Errorf("get user media failed: %s", e)
			wg.Done()
		})
	}()
	wg.Wait()
	return
}

// +build js

package webrtc

import (
	"fmt"
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

// IceServer ...
type IceServer struct {
	o              *js.Object
	Urls           []string `js:"urls"`
	Username       string   `js:"username"`
	Credential     string   `js:"credential"`
	CredentialType string   `js:"credentialType"`
}

// NewIceServer ...
func NewIceServer(urls []string, params ...string) *IceServer {
	username := ""
	credential := ""
	if len(params) > 0 {
		username = params[0]
	}
	if len(params) > 1 {
		credential = params[1]
	}
	is := &IceServer{o: js.Global.Get("Object").New()}
	is.Urls = urls
	is.Username = username
	is.Credential = credential
	is.CredentialType = "password"
	return is
}

// Configuration ...
type Configuration struct {
	o                  *js.Object
	IceServers         []*IceServer `js:"iceServers"`
	IceTransportPolicy string       `js:"iceTransportPolicy"`
	BundlePolicy       string       `js:"bundlePolicy"`
	RTCPMuxPolicy      string       `js:"rtcpMuxPolicy"`
	PeerIdentity       string       `js:"peerIdentity"`
}

// NewConfiguration ...
func NewConfiguration() *Configuration {
	c := &Configuration{o: js.Global.Get("Object").New()}
	c.IceServers = []*IceServer{}
	c.IceTransportPolicy = "all"
	c.BundlePolicy = "balanced"
	c.RTCPMuxPolicy = "require"
	c.PeerIdentity = ""
	return c
}

// AddIceServer ...
func (config *Configuration) AddIceServer(params ...string) error {
	urls := strings.Split(params[0], ",")
	for i, url := range urls {
		url = strings.TrimSpace(url)
		if !strings.HasPrefix(url, "stun:") &&
			!strings.HasPrefix(url, "turn:") {
			return fmt.Errorf("IceServer: received malformed url: <%s>", url)
		}
		urls[i] = url
	}
	config.IceServers = append(config.IceServers,
		NewIceServer(urls, params[1:]...),
	)
	return nil
}

// VideoMandatory ...
type VideoMandatory struct {
	o             *js.Object
	MaxWidth      int     `js:"maxWidth"`
	MinWidth      int     `js:"minWidth"`
	MaxHeight     int     `js:"maxHeight"`
	MinHeight     int     `js:"minHeight"`
	MaxFrameRate  float64 `js:"maxFrameRate"`
	MinFrameRate  float64 `js:"minFrameRate"`
	AspectRate    float64 `js:"aspectRate"`
	MaxAspectRate float64 `js:"maxAspectRate"`
	MinAspectRate float64 `js:"minAspectRate"`
}

// AudioMandatory ...
type AudioMandatory struct {
	o                *js.Object
	EchoCancellation bool    `js:"echoCancellation"`
	MaxChannelCount  int     `js:"maxChannelCount"`
	MinChannelCount  int     `js:"minChannelCount"`
	MaxSampleRate    int     `js:"maxSampleRate"`
	MinSampleRate    int     `js:"minSampleRate"`
	MaxLatency       float64 `js:"maxLatency"`
	MinLatency       float64 `js:"minLatency"`
	MaxVolume        float64 `js:"maxVolume"`
	MinVolume        float64 `js:"minVolume"`
}

// VideoConstraints ...
type VideoConstraints struct {
	o         *js.Object
	Mandatory VideoMandatory `js:"mandatory"`
}

// AudioConstraints ...
type AudioConstraints struct {
	o         *js.Object
	Mandatory AudioMandatory `js:"mandatory"`
}

// NewVideoConstraints ...
func NewVideoConstraints() *VideoConstraints {
	c := &VideoConstraints{o: js.Global.Get("Object").New()}
	c.Mandatory = VideoMandatory{o: js.Global.Get("Object").New()}
	return c
}

// NewAudioConstraints ...
func NewAudioConstraints() *AudioConstraints {
	c := &AudioConstraints{o: js.Global.Get("Object").New()}
	c.Mandatory = AudioMandatory{o: js.Global.Get("Object").New()}
	return c
}

// Constraints ...
type Constraints struct {
	o     *js.Object
	Video interface{} `js:"video"` // bool or *VideoConstraints
	Audio interface{} `js:"audio"` // bool or *AudioConstraints
}

// NewConstraints ...
func NewConstraints(video, audio interface{}) *Constraints {
	c := &Constraints{o: js.Global.Get("Object").New()}
	c.Video = video
	c.Audio = audio
	return c
}

// IceCandidate ...
type IceCandidate struct {
	o             *js.Object
	Candidate     string `js:"candidate",json:"candidate"`
	SdpMid        string `js:"sdpMid",json:"sdpMid"`
	SdpMLineIndex int    `js:"sdpMLineIndex",json:"sdpMLineIndex"`
}

// NewIceCandidate ...
func NewIceCandidate(candidate, sdpmid string, sdpMLineIndex int) *IceCandidate {
	ic := &IceCandidate{o: js.Global.Get("Object").New()}
	ic.Candidate = candidate
	ic.SdpMid = sdpmid
	ic.SdpMLineIndex = sdpMLineIndex
	return ic
}

// NewIceCandidateFromObj ...
func NewIceCandidateFromObj(obj *js.Object) *IceCandidate {
	ic := &IceCandidate{o: obj}
	ic.Candidate = obj.Get("candidate").String()
	ic.SdpMid = obj.Get("sdpMid").String()
	ic.SdpMLineIndex = obj.Get("sdpMLineIndex").Int()
	return ic
}

// SessionDescription ...
type SessionDescription struct {
	o    *js.Object
	Type string `js:"type",json:"type"`
	Sdp  string `js:"sdp",json:"sdp"`
}

// NewSessionDescription ...
func NewSessionDescription(tp, sdp string) *SessionDescription {
	sd := &SessionDescription{o: js.Global.Get("Object").New()}
	sd.Type = tp
	sd.Sdp = sdp
	return sd
}

// NewSessionDescriptionFromObj ...
func NewSessionDescriptionFromObj(obj *js.Object) *SessionDescription {
	sd := &SessionDescription{o: obj}
	sd.Type = obj.Get("type").String()
	sd.Sdp = obj.Get("sdp").String()
	return sd
}

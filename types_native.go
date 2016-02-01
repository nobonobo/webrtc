// +build !js

package webrtc

import (
	"fmt"
	"strings"
)

// IceServer ...
type IceServer struct {
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
	is := &IceServer{}
	is.Urls = urls
	is.Username = username
	is.Credential = credential
	is.CredentialType = "password"
	return is
}

// Configuration ...
type Configuration struct {
	IceServers         []*IceServer `js:"iceServers"`
	IceTransportPolicy string       `js:"iceTransportPolicy"`
	BundlePolicy       string       `js:"bundlePolicy"`
	RTCPMuxPolicy      string       `js:"rtcpMuxPolicy"`
	PeerIdentity       string       `js:"peerIdentity"`
}

// NewConfiguration ...
func NewConfiguration() *Configuration {
	c := &Configuration{}
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
	Mandatory VideoMandatory `js:"mandatory"`
}

// AudioConstraints ...
type AudioConstraints struct {
	Mandatory AudioMandatory `js:"mandatory"`
}

// Constraints ...
type Constraints struct {
	Video interface{} `js:"video"` // bool or *VideoConstraints
	Audio interface{} `js:"audio"` // bool or *AudioConstraints
}

// NewConstraints ...
func NewConstraints(video, audio interface{}) *Constraints {
	c := &Constraints{}
	c.Video = video
	c.Audio = audio
	return c
}

// IceCandidate ...
type IceCandidate struct {
	Candidate     string `json:"candidate",js:"candidate"`
	SdpMid        string `json:"sdpMid",js:"sdpMid"`
	SdpMLineIndex int    `json:"sdpMLineIndex",js:"sdpMLineIndex"`
}

// NewIceCandidate ...
func NewIceCandidate(candidate, sdpmid string, sdpMLineIndex int) *IceCandidate {
	return &IceCandidate{
		Candidate:     candidate,
		SdpMid:        sdpmid,
		SdpMLineIndex: sdpMLineIndex,
	}
}

// SessionDescription ...
type SessionDescription struct {
	Type string `json:"type",js:"type"`
	Sdp  string `json:"sdp",js:"sdp"`
}

// NewSessionDescription ...
func NewSessionDescription(tp, sdp string) *SessionDescription {
	return &SessionDescription{
		Type: tp,
		Sdp:  sdp,
	}
}

// +build !js

package webrtc

import (
	"log"
	"strings"

	org "github.com/keroserene/go-webrtc"
)

func init() {
	org.SetLoggingVerbosity(0)
}

func find(s []string, search string) int {
	for i := 0; i < len(s); i++ {
		if strings.ToLower(s[i]) == strings.ToLower(search) {
			return i
		}
	}
	return -1
}

// PeerConnection ...
type PeerConnection struct {
	pc *org.PeerConnection
}

// NewPeerConnection ...
func NewPeerConnection(config *Configuration) (*PeerConnection, error) {
	conf := org.NewConfiguration()
	for _, s := range config.IceServers {
		args := []string{strings.Join(s.Urls, ",")}
		if len(s.Username) > 0 {
			args = append(args, s.Username)
		}
		if len(s.Credential) > 0 {
			args = append(args, s.Credential)
		}
		if err := conf.AddIceServer(args...); err != nil {
			return nil, err
		}
	}
	conf.BundlePolicy = org.BundlePolicy(find(org.BundlePolicyString, config.BundlePolicy))
	conf.IceTransportPolicy = org.IceTransportPolicy(find(
		org.IceTransportPolicyString, config.IceTransportPolicy))
	conf.PeerIdentity = config.PeerIdentity
	pc, err := org.NewPeerConnection(conf)
	if err != nil {
		return nil, err
	}
	return &PeerConnection{pc: pc}, nil
}

// OnNegotiationNeeded ...
func (pc *PeerConnection) OnNegotiationNeeded(cb func()) {
	pc.pc.OnNegotiationNeeded = cb
}

// OnIceCandidate ...
func (pc *PeerConnection) OnIceCandidate(cb func(*IceCandidate)) {
	pc.pc.OnIceCandidate = func(ic org.IceCandidate) {
		log.Println("ice:", ic)
		cb(&IceCandidate{
			Candidate:     ic.Candidate,
			SdpMid:        ic.SdpMid,
			SdpMLineIndex: ic.SdpMLineIndex,
		})
	}
}

// OnIceCandidateError ...
func (pc *PeerConnection) OnIceCandidateError(cb func()) {
	pc.pc.OnIceCandidateError = cb
}

// OnSignalingStateChange ...
func (pc *PeerConnection) OnSignalingStateChange(cb func(SignalingState string)) {
	pc.pc.OnSignalingStateChange = func(s org.SignalingState) {
		cb(org.SignalingStateString[s])
	}
}

// OnIceConnectionStateChange ...
func (pc *PeerConnection) OnIceConnectionStateChange(cb func(IceConnectionState string)) {
	pc.pc.OnIceConnectionStateChange = func(s org.IceConnectionState) {
		cb(org.IceConnectionStateString[s])
	}
}

// OnIceGatheringStateChange ...
func (pc *PeerConnection) OnIceGatheringStateChange(cb func(IceGatheringState string)) {
	pc.pc.OnIceGatheringStateChange = func(s org.IceGatheringState) {
		cb(org.IceGatheringStateString[s])
	}
}

// OnConnectionStateChange ...
func (pc *PeerConnection) OnConnectionStateChange(cb func(PeerConnectionState string)) {
	pc.pc.OnConnectionStateChange = func(s org.PeerConnectionState) {
		cb(org.PeerConnectionStateString[s])
	}
}

// OnDataChannel ...
func (pc *PeerConnection) OnDataChannel(cb func(*DataChannel)) {
	pc.pc.OnDataChannel = func(dc *org.DataChannel) {
		cb(&DataChannel{
			dc: dc,
		})
	}
}

// OnAddStream ...
func (pc *PeerConnection) OnAddStream(cb func(*MediaStream)) {}

// OnRemoveStream ...
func (pc *PeerConnection) OnRemoveStream(cb func(*MediaStream)) {}

// AddStream ...
func (pc *PeerConnection) AddStream(stream *MediaStream) (err error) { return nil }

// AddIceCandidate ...
func (pc *PeerConnection) AddIceCandidate(ic *IceCandidate) error {
	return pc.pc.AddIceCandidate(org.IceCandidate{
		Candidate:     ic.Candidate,
		SdpMid:        ic.SdpMid,
		SdpMLineIndex: ic.SdpMLineIndex,
	})
}

// Close ...
func (pc *PeerConnection) Close() error {
	return pc.pc.Close()
}

// ConnectionState ...
func (pc *PeerConnection) ConnectionState() string {
	return org.PeerConnectionStateString[pc.pc.ConnectionState()]
}

// CreateAnswer ...
func (pc *PeerConnection) CreateAnswer() (*SessionDescription, error) {
	sd, err := pc.pc.CreateAnswer()
	if err != nil {
		return nil, err
	}
	return &SessionDescription{
		Type: sd.Type,
		Sdp:  sd.Sdp,
	}, nil
}

// CreateDataChannel ...
func (pc *PeerConnection) CreateDataChannel(label string) (*DataChannel, error) {
	dc, err := pc.pc.CreateDataChannel(label, org.Init{})
	if err != nil {
		return nil, err
	}
	return &DataChannel{dc}, nil
}

// CreateOffer ...
func (pc *PeerConnection) CreateOffer() (*SessionDescription, error) {
	sd, err := pc.pc.CreateOffer()
	if err != nil {
		return nil, err
	}
	return &SessionDescription{
		Type: sd.Type,
		Sdp:  sd.Sdp,
	}, nil
}

// IceGatheringState ...
func (pc *PeerConnection) IceGatheringState() string {
	return ""
}

// LocalDescription ...
func (pc *PeerConnection) LocalDescription() (sdp *SessionDescription) {
	sd := pc.pc.LocalDescription()
	return &SessionDescription{
		Type: sd.Type,
		Sdp:  sd.Sdp,
	}
}

// RemoteDescription ...
func (pc *PeerConnection) RemoteDescription() (sdp *SessionDescription) {
	sd := pc.pc.RemoteDescription()
	return &SessionDescription{
		Type: sd.Type,
		Sdp:  sd.Sdp,
	}
}

// SetLocalDescription ...
func (pc *PeerConnection) SetLocalDescription(sdp *SessionDescription) error {
	return pc.pc.SetLocalDescription(&org.SessionDescription{
		Type: sdp.Type,
		Sdp:  sdp.Sdp,
	})
}

// SetRemoteDescription ...
func (pc *PeerConnection) SetRemoteDescription(sdp *SessionDescription) error {
	return pc.pc.SetRemoteDescription(&org.SessionDescription{
		Type: sdp.Type,
		Sdp:  sdp.Sdp,
	})
}

// SignalingState ...
func (pc *PeerConnection) SignalingState() string {
	return org.SignalingStateString[pc.pc.SignalingState()]
}

// DataChannel ...
type DataChannel struct {
	dc *org.DataChannel
}

// OnOpen ...
func (c *DataChannel) OnOpen(cb func()) {
	c.dc.OnOpen = cb
}

// OnClose ...
func (c *DataChannel) OnClose(cb func()) {
	c.dc.OnClose = cb
}

// OnMessage ...
func (c *DataChannel) OnMessage(cb func([]byte)) {
	c.dc.OnMessage = cb
}

// Close ...
func (c *DataChannel) Close() error {
	return c.dc.Close()
}

// ID ...
func (c *DataChannel) ID() int {
	return c.dc.ID()
}

// ReadyState ...
func (c *DataChannel) ReadyState() string {
	return ""
}

// Send ...
func (c *DataChannel) Send(data []byte) {
	c.dc.Send(data)
}

// Label ...
func (c *DataChannel) Label() string {
	return c.dc.Label()
}

// MediaStream ...
type MediaStream struct{}

// GetUserMedia ...
func GetUserMedia(constraints *Constraints) (stream *MediaStream, err error) {
	panic("not supported")
}

package webrtc

import (
	"time"

	"github.com/pion/webrtc/v3"
)

// SignalMessage represents the message sent during signalling
type SignalMessage struct {
	Type      string                   `json:"type"`
	SDP       string                   `json:"sdp,omitempty"`
	Candidate *webrtc.ICECandidateInit `json:"candidate,omitempty"`
	UserID    string                   `json:"userId"`
	MeetingID string                   `json:"meetingId"`
	TrackID   string                   `json:"trackId,omitempty"`
	Target    string                   `json:"target,omitempty"` // target user id for p2p messages
}

// Room represents a meeting room with multiple peers
type Room struct {
	ID        string
	Peers     map[string]*Peer
	Tracks    map[string]*webrtc.TrackLocalStaticRTP
	CreatedAt time.Time
	closeChan chan struct{}
}

type Peer struct {
	ID            string
	Connection    *webrtc.PeerConnection
	Tracks        map[string]*webrtc.TrackLocalStaticRTP
	DataChannel   *webrtc.DataChannel
	Room          *Room
	SignalChannel chan *SignalMessage
}

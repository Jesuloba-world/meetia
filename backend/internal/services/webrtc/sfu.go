package webrtc

import (
	"errors"
	"io"
	"log"
	"sync"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
)

type SFUService struct {
	rooms      map[string]*Room
	roomsMutex sync.Mutex
	config     webrtc.Configuration
}

func NewSFUService() *SFUService {
	return &SFUService{
		rooms: make(map[string]*Room),
		config: webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{
						"stun:stun.l.google.com:19302",
						"stun:stun.l.google.com:5349",
						"stun:stun1.l.google.com:3478",
						"stun:stun1.l.google.com:5349",
						"stun:stun2.l.google.com:19302",
						"stun:stun2.l.google.com:5349",
						"stun:stun3.l.google.com:3478",
						"stun:stun3.l.google.com:5349",
						"stun:stun4.l.google.com:19302",
						"stun:stun4.l.google.com:5349",
					},
				},
				// {
				// 	URLs:       []string{"turn:localhost:3478"},
				// 	Username:   "meetia_user",
				// 	Credential: "strong_password",
				// },
			},
			ICETransportPolicy: webrtc.ICETransportPolicyAll,
		},
	}
}

func (s *SFUService) GetOrCreateRoom(roomID string) *Room {
	s.roomsMutex.Lock()
	defer s.roomsMutex.Unlock()

	if room, exists := s.rooms[roomID]; exists {
		return room
	}

	room := &Room{
		ID:        roomID,
		Peers:     make(map[string]*Peer),
		Tracks:    make(map[string]*webrtc.TrackLocalStaticRTP),
		CreatedAt: time.Now(),
		closeChan: make(chan struct{}),
	}

	s.rooms[roomID] = room
	return room
}

func (s *SFUService) RemoveRoom(roomID string) {
	s.roomsMutex.Lock()
	defer s.roomsMutex.Unlock()

	if room, exists := s.rooms[roomID]; exists {
		close(room.closeChan)
		delete(s.rooms, roomID)
	}
}

func (s *SFUService) CreatePeerConnection(roomID string, peerID string) (*Peer, error) {
	room := s.GetOrCreateRoom(roomID)

	// create new peer connection
	peerConnection, err := webrtc.NewPeerConnection(s.config)
	if err != nil {
		return nil, err
	}

	peer := &Peer{
		ID:            peerID,
		Connection:    peerConnection,
		Tracks:        make(map[string]*webrtc.TrackLocalStaticRTP),
		Room:          room,
		SignalChannel: make(chan *SignalMessage, 100),
	}

	// set up data channel for chat messages and signaling
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		peerConnection.Close()
		return nil, err
	}
	peer.DataChannel = dataChannel

	// add peer to room
	room.Peers[peerID] = peer

	// setup ICE connection state handler
	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		log.Printf("Peer %s ICE connection state: %s\n", peerID, state.String())

		if state == webrtc.ICEConnectionStateDisconnected {
			// Start grace period for reconnection
			time.AfterFunc(10*time.Second, func() {
				if peer.Connection.ICEConnectionState() == webrtc.ICEConnectionStateDisconnected {
					log.Printf("Peer %s permanent disconnect, cleaning up\n", peerID)
					s.roomsMutex.Lock()
					delete(room.Peers, peerID)
					s.roomsMutex.Unlock()
					peerConnection.Close()
				}
			})
			return
		}

		if state == webrtc.ICEConnectionStateFailed ||
			state == webrtc.ICEConnectionStateClosed {
			s.roomsMutex.Lock()
			delete(room.Peers, peerID)
			s.roomsMutex.Unlock()
			peerConnection.Close()
		}
	})

	// peerConnection.OnNegotiationNeeded(func() {
	// 	log.Printf("Peer %s negotiation needed", peerID)
	// 	offer, err := peerConnection.CreateOffer(nil)
	// 	if err != nil {
	// 		log.Printf("Failed to create renegotiation offer: %v", err)
	// 		return
	// 	}

	// 	if err := peerConnection.SetLocalDescription(offer); err != nil {
	// 		log.Printf("Failed to set local description: %v", err)
	// 		return
	// 	}

	// 	peer.SignalChannel <- &SignalMessage{
	// 		Type:      "offer",
	// 		SDP:       offer.SDP,
	// 		UserID:    peerID,
	// 		MeetingID: roomID,
	// 	}
	// })

	// handle incoming tracks
	peerConnection.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Peer %s added track: %s\n", peerID, remoteTrack.ID())

		// Create a local track to forward to other peers
		trackLocal, err := webrtc.NewTrackLocalStaticRTP(remoteTrack.Codec().RTPCodecCapability, remoteTrack.ID(), peerID)
		if err != nil {
			log.Printf("Failed to create local track: %v\n", err)
			return
		}

		// add track to room
		room.Tracks[remoteTrack.ID()] = trackLocal
		peer.Tracks[remoteTrack.ID()] = trackLocal

		// send tracks to existing peers in the room
		for otherPeerID, otherPeer := range room.Peers {
			if otherPeerID == peerID {
				continue
			}

			// add track to other peers
			if _, err := otherPeer.Connection.AddTrack(trackLocal); err != nil {
				log.Printf("Failed to add track to peer %s: %v\n", otherPeerID, err)
				continue
			}

			// create an offer for the other peer
			offer, err := otherPeer.Connection.CreateOffer(nil)
			if err != nil {
				log.Printf("Failed to create offer: %v\n", err)
				continue
			}

			// set the local description
			if err := otherPeer.Connection.SetLocalDescription(offer); err != nil {
				log.Printf("Failed to set local description: %v\n", err)
				continue
			}

			// send the offer to the other peer
			otherPeer.SignalChannel <- &SignalMessage{
				Type:      "offer",
				SDP:       offer.SDP,
				UserID:    peerID,
				MeetingID: roomID,
				TrackID:   remoteTrack.ID(),
			}
		}

		// read packets from the track and forward them
		go func() {
			rtpBuf := make([]byte, 1500)
			for {
				select {
				case <-room.closeChan:
					return
				default:
					i, _, readErr := remoteTrack.Read(rtpBuf)
					if readErr != nil {
						return
					}

					// write to local track
					if _, err = trackLocal.Write(rtpBuf[:i]); err != nil && !errors.Is(err, io.ErrClosedPipe) {
						log.Printf("Failed to write to local track: %v\n", err)
						return
					}
				}
			}
		}()

		// send RTCP PLI packets for video tracks
		go func() {
			ticker := time.NewTicker(time.Second * 3)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if remoteTrack.Kind() == webrtc.RTPCodecTypeVideo {
						_ = peerConnection.WriteRTCP([]rtcp.Packet{
							&rtcp.PictureLossIndication{
								MediaSSRC: uint32(remoteTrack.SSRC()),
							},
						})
					}
				case <-room.closeChan:
					return
				}
			}
		}()
	})

	return peer, nil
}

func (s *SFUService) HandleOffer(peer *Peer, offer webrtc.SessionDescription) (webrtc.SessionDescription, error) {
	if peer.Connection.SignalingState() == webrtc.SignalingStateHaveRemoteOffer {
		peer.Connection.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer})
	}

	err := peer.Connection.SetRemoteDescription(offer)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	// create answer
	answer, err := peer.Connection.CreateAnswer(nil)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	// set local description
	err = peer.Connection.SetLocalDescription(answer)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	return answer, nil
}

func (s *SFUService) HandleAnswer(peer *Peer, answer webrtc.SessionDescription) error {
	if peer.Connection.SignalingState() == webrtc.SignalingStateHaveLocalOffer {
		return peer.Connection.SetRemoteDescription(answer)
	}

	log.Printf("Received answer in unexpected state: %s", peer.Connection.SignalingState().String())
	return nil
}

func (s *SFUService) HandleCandidate(peer *Peer, candidate webrtc.ICECandidateInit) error {
	// Validate candidate format
	if candidate.Candidate == "" {
		return errors.New("empty ICE candidate")
	}

	// Check connection state
	if peer.Connection.ICEConnectionState() == webrtc.ICEConnectionStateClosed {
		return errors.New("connection closed")
	}

	return peer.Connection.AddICECandidate(candidate)
}

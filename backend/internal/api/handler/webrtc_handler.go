package handler

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	pion "github.com/pion/webrtc/v3"

	"github.com/meetia/backend/internal/services/webrtc"
)

type WebRTCHandler struct {
	sfuService *webrtc.SFUService
	tokenAuth  *jwtauth.JWTAuth
}

func NewWebRTCHandler(sfuService *webrtc.SFUService, tokenAuth *jwtauth.JWTAuth) *WebRTCHandler {
	return &WebRTCHandler{
		sfuService: sfuService,
		tokenAuth:  tokenAuth,
	}
}

func (h *WebRTCHandler) RegisterRoutes(r chi.Router, api huma.API) {
	r.Route("/api/rtc", func(r chi.Router) {
		r.Use(jwtauth.Verifier(h.tokenAuth))
		r.Use(jwtauth.Authenticator(h.tokenAuth))
		r.Get("/signal/{meetingID}", h.HandleWebSocket)
	})
}

func (h *WebRTCHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	meetingID := chi.URLParam(r, "meetingID")
	if meetingID == "" {
		http.Error(w, "Meeting ID is required", http.StatusBadRequest)
		return
	}

	// get user ID from jwt token
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// accept websocket connection
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		CompressionMode:    websocket.CompressionDisabled,
	})
	if err != nil {
		log.Printf("Websocket accept error: %v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "Connection closed")

	// create peer connection
	peer, err := h.sfuService.CreatePeerConnection(meetingID, userID)
	if err != nil {
		log.Printf("Failed to create peer connection: %v", err)
		c.Close(websocket.StatusInternalError, "Failed to create peer connection")
		return
	}

	// signal channel to coordinate websocket communication
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// wait group to ensure all goroutines finish before closing
	var wg sync.WaitGroup
	wg.Add(2)

	// read messages from websocket
	go func() {
		defer wg.Done()
		defer cancel() // cancel context when this goroutine exits

		for {
			var msg webrtc.SignalMessage
			err := wsjson.Read(ctx, c, &msg)
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				return
			}

			// route based on message type
			switch msg.Type {
			case "offer":
				sdp := pion.SessionDescription{
					Type: pion.SDPTypeOffer,
					SDP:  msg.SDP,
				}
				answer, err := h.sfuService.HandleOffer(peer, sdp)
				if err != nil {
					log.Printf("Handle offer error: %v", err)
					continue
				}

				// send answer back
				responseMsg := webrtc.SignalMessage{
					Type:      "answer",
					SDP:       answer.SDP,
					UserID:    userID,
					MeetingID: meetingID,
				}
				if err := wsjson.Write(ctx, c, responseMsg); err != nil {
					log.Printf("WebSocket write error: %v", err)
					return
				}

			case "answer":
				sdp := pion.SessionDescription{
					Type: pion.SDPTypeAnswer,
					SDP:  msg.SDP,
				}
				if err := h.sfuService.HandleAnswer(peer, sdp); err != nil {
					log.Printf("Handle answer error: %v", err)
				}

			case "candidate":
				if msg.Candidate != nil {
					if err := h.sfuService.HandleCandidate(peer, *msg.Candidate); err != nil {
						log.Printf("Handle candidate error: %v", err)
					}
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		defer cancel() // Cancel context when this goroutine exits

		for {
			select {
			case msg := <-peer.SignalChannel:
				if err := wsjson.Write(ctx, c, msg); err != nil {
					log.Printf("WebSocket write error: %v", err)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()
	c.Close(websocket.StatusNormalClosure, "Connection closed")
}

package handler

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/jwtauth/v5"
	pion "github.com/pion/webrtc/v3"

	"github.com/meetia/backend/internal/api/middleware"
	"github.com/meetia/backend/internal/services/webrtc"
	humagroup "github.com/meetia/backend/lib/humaGroup"
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

func (h *WebRTCHandler) RegisterRoutes(api huma.API) {
	rtcGroup := humagroup.NewHumaGroup(
		api,
		"/api/rtc",
		[]string{"WebRTC"},
		middleware.JWTMiddleware(h.tokenAuth),
		middleware.WithHttpContext,
	)

	humagroup.Get(
		rtcGroup,
		"/signal/{meetingID}",
		h.HandleWebSocket,
		"wsSignal",
		&humagroup.HumaGroupOptions{
			Summary:     "WebRTC Signaling",
			Description: "WebSocket endpoint for WebRTC signaling",
		},
	)
}

type handleWebSocketInput struct {
	AuthParam

	MeetingID string `path:"meetingID" required:"true" doc:"Meeting ID"`
}

func (h *WebRTCHandler) HandleWebSocket(ctx context.Context, input *handleWebSocketInput) (*struct{}, error) {
	meetingID := input.MeetingID

	// get user ID from jwt token
	userID, err := getUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	r, w, ok := middleware.GetHttpContext(ctx)
	if !ok {
		return nil, fmt.Errorf("http context not available")
	}

	// accept websocket connection
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		CompressionMode:    websocket.CompressionDisabled,
	})
	if err != nil {
		log.Printf("Websocket accept error: %v", err)
		return nil, fmt.Errorf("websocket accept error: %v", err)
	}
	defer c.Close(websocket.StatusInternalError, "Connection closed")

	// create peer connection
	peer, err := h.sfuService.CreatePeerConnection(meetingID, userID)
	if err != nil {
		log.Printf("Failed to create peer connection: %v", err)
		c.Close(websocket.StatusInternalError, "Failed to create peer connection")
		return nil, fmt.Errorf("failed to create peer connection: %v", err)
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

	return &struct{}{}, nil
}

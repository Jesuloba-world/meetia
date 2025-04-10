package handler

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/jwtauth/v5"

	"github.com/meetia/backend/internal/api/middleware"
	"github.com/meetia/backend/internal/services/meeting"
	humagroup "github.com/meetia/backend/lib/humaGroup"

)

type MeetingHandler struct {
	meetingService *meeting.MeetingService
	tokenAuth      *jwtauth.JWTAuth
}

func NewMeetinghandler(meetingService *meeting.MeetingService, tokenAuth *jwtauth.JWTAuth) *MeetingHandler {
	return &MeetingHandler{
		meetingService: meetingService,
		tokenAuth:      tokenAuth,
	}
}

func (h *MeetingHandler) RegisterRoutes(api huma.API) {
	meetingGroup := humagroup.NewHumaGroup(api, "/api/meetings", []string{"Meetings"}, middleware.JWTMiddleware(h.tokenAuth))

	humagroup.Post(meetingGroup, "/", h.CreateMeeting, "CreateMeeting", &humagroup.HumaGroupOptions{
		Summary:     "Create a new meeting",
		Description: "Creates a new meeting with the current user as host",
	})
}

type CreateMeetingRequest struct {
	AuthParam

	Body struct {
		Title     string `json:"title" doc:"Meeting title" example:"Team Weekly Sync"`
		IsPrivate bool   `json:"isPrivate" doc:"Whether the meeting requires a password" example:"false"`
		Password  string `json:"password,omitempty" doc:"Password if the meeting is private" example:"securepass123"`
	} `json:"body"`
}

type UserInfoSmall struct {
	ID          string `json:"id" doc:"User ID"`
	DisplayName string `json:"displayName" doc:"User display name"`
}

type MeetingResponse struct {
	ID          string         `json:"id" doc:"Meeting unique identifier"`
	Title       string         `json:"title" doc:"Meeting title"`
	HostID      string         `json:"hostId" doc:"ID of the meeting host"`
	MeetingCode string         `json:"meetingCode" doc:"Unique code to join the meeting"`
	IsPrivate   bool           `json:"isPrivate" doc:"Whether the meeting requires a password"`
	CreatedAt   time.Time      `json:"createdAt" doc:"When the meeting was created"`
	Host        *UserInfoSmall `json:"host,omitempty" doc:"Host details"`
}

type CreateMeetingResponse struct {
	Body struct {
		Meeting MeetingResponse `json:"meeting"`
	} `json:"body"`
}

func (h *MeetingHandler) CreateMeeting(ctx context.Context, input *CreateMeetingRequest) (*CreateMeetingResponse, error) {
	// get user ID from jwt token
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, huma.Error401Unauthorized("invalid token", err)
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, huma.Error401Unauthorized("invalid token claims")
	}

	meeting, err := h.meetingService.CreateMeeting(ctx, input.Body.Title, userID, input.Body.IsPrivate, input.Body.Password)
	if err != nil {
		return nil, huma.Error500InternalServerError("an error occured while creating meeting", err)
	}

	resp := &CreateMeetingResponse{}
	meetingResp := MeetingResponse{
		ID:          meeting.ID,
		Title:       meeting.Title,
		HostID:      meeting.HostID,
		MeetingCode: meeting.MeetingCode,
		IsPrivate:   meeting.IsPrivate,
		CreatedAt:   meeting.CreatedAt,
		Host: &UserInfoSmall{
			ID:          meeting.HostID,
			DisplayName: meeting.Host.DisplayName,
		},
	}
	resp.Body.Meeting = meetingResp

	return resp, nil
}

type JoinMeetingRequest struct {
	MeetingCode string `json:"meetingCode" doc:"Meeting code to join" example:"ABCDEF12345"`
	Password    string `json:"password,omitempty" doc:"Password if the meeting is private" example:"securepass123"`
}

type ChatMessageRequest struct {
	Message string `json:"message" doc:"Chat message content" example:"Hello everyone!"`
}

type ChatMessageResponse struct {
	ID        string    `json:"id" doc:"Message unique identifier"`
	MeetingID string    `json:"meetingId" doc:"ID of the meeting this message belongs to"`
	UserID    string    `json:"userId" doc:"ID of the user who sent the message"`
	Message   string    `json:"message" doc:"Message content"`
	SentAt    time.Time `json:"sentAt" doc:"When the message was sent"`
	User      struct {
		DisplayName string `json:"displayName" doc:"User display name"`
	} `json:"user" doc:"User details"`
}

type ParticipantResponse struct {
	ID     string `json:"id" doc:"Participant record ID"`
	UserID string `json:"userId" doc:"User ID"`
	Role   string `json:"role" doc:"Participant role (host, co-host, participant)"`
	User   struct {
		DisplayName string `json:"displayName" doc:"User display name"`
	} `json:"user" doc:"User details"`
}

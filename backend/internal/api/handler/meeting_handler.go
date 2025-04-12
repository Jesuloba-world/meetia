package handler

import (
	"context"
	"errors"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/jwtauth/v5"

	"github.com/meetia/backend/internal/api/middleware"
	"github.com/meetia/backend/internal/models"
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
	humagroup.Post(meetingGroup, "/join", h.JoinMeeting, "JoinMeeting", &humagroup.HumaGroupOptions{
		Summary:     "Join an existing meeting",
		Description: "Join a meeting using its meeting code and password (if required)",
	})
	humagroup.Get(meetingGroup, "/", h.ListMeetings, "ListMeetings", &humagroup.HumaGroupOptions{
		Summary:     "List user's meetings",
		Description: "Get a list of active meetings where the user is a host or participant",
	})
	humagroup.Get(meetingGroup, "/{id}", h.GetMeeting, "GetMeeting", &humagroup.HumaGroupOptions{
		Summary:     "Get meeting details",
		Description: "Get details about a specific meeting",
	})
	humagroup.Post(meetingGroup, "/{id}/end", h.EndMeeting, "EndMeeting", &humagroup.HumaGroupOptions{
		Summary:     "End a meeting",
		Description: "End a meeting (host only)",
	})
	humagroup.Get(meetingGroup, "/{id}/participants", h.GetParticipants, "GetParticipants", &humagroup.HumaGroupOptions{
		Summary:     "Get meeting participants",
		Description: "Get a list of participants in a meeting",
	})
	humagroup.Post(meetingGroup, "/{id}/chat", h.SendChatMessage, "SendChatMessage", &humagroup.HumaGroupOptions{
		Summary:     "Send a chat message",
		Description: "Send a chat message in a meeting",
	})
	humagroup.Get(meetingGroup, "/{id}/chat", h.GetChatMessages, "GetChatMessages", &humagroup.HumaGroupOptions{
		Summary:     "Get chat messages",
		Description: "Get chat message history for a meeting",
	})
}

type CreateMeetingRequest struct {
	AuthParam

	Body struct {
		Title     string `json:"title" doc:"Meeting title" example:"Team Weekly Sync"`
		IsPrivate bool   `json:"isPrivate" doc:"Whether the meeting requires a password" example:"false"`
		Password  string `json:"password,omitempty" doc:"Password if the meeting is private" example:"securepass123"`
	}
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
	}
}

func (h *MeetingHandler) CreateMeeting(ctx context.Context, input *CreateMeetingRequest) (*CreateMeetingResponse, error) {
	// get user ID from jwt token
	userID, err := getUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	meeting, err := h.meetingService.CreateMeeting(ctx, input.Body.Title, userID, input.Body.IsPrivate, input.Body.Password)
	if err != nil {
		return nil, huma.Error500InternalServerError("an error occured while creating meeting", err)
	}

	resp := &CreateMeetingResponse{}
	meetingResp := meetingToResponse(meeting)
	resp.Body.Meeting = meetingResp
	return resp, nil
}

type JoinMeetingRequest struct {
	AuthParam

	Body struct {
		MeetingCode string `json:"meetingCode" doc:"Meeting code to join" example:"ABCDEF12345"`
		Password    string `json:"password,omitempty" doc:"Password if the meeting is private" example:"securepass123"`
	}
}

type JoinMeetingResponse struct {
	Body struct {
		Meeting MeetingResponse `json:"meeting"`
	}
}

func (h *MeetingHandler) JoinMeeting(ctx context.Context, input *JoinMeetingRequest) (*JoinMeetingResponse, error) {
	userID, err := getUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	joinedmeeting, err := h.meetingService.JoinMeeting(ctx, input.Body.MeetingCode, userID, input.Body.Password)
	if err != nil {
		switch {
		case errors.Is(err, meeting.ErrMeetingNotFound):
			return nil, huma.Error404NotFound("meeting not found", err)
		case errors.Is(err, meeting.ErrInvalidPassword):
			return nil, huma.Error401Unauthorized("Invalid password", err)
		default:
			return nil, huma.Error500InternalServerError("an error occured", err)
		}
	}

	resp := &JoinMeetingResponse{}
	resp.Body.Meeting = meetingToResponse(joinedmeeting)
	return resp, nil
}

type ListMeetingsRequest struct {
	AuthParam
}

type ListMeetingsResponse struct {
	Body struct {
		Meetings []MeetingResponse `json:"meetings" doc:"list of active meetings for user"`
	}
}

func (h *MeetingHandler) ListMeetings(ctx context.Context, input *ListMeetingsRequest) (*ListMeetingsResponse, error) {
	userID, err := getUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	meetings, err := h.meetingService.GetActiveUserMeetings(ctx, userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to list meetings", err)
	}

	response := make([]MeetingResponse, len(meetings))
	for i, meeting := range meetings {
		response[i] = meetingToResponse(meeting)
	}

	resp := &ListMeetingsResponse{}
	resp.Body.Meetings = response
	return resp, nil
}

type GetMeetingRequest struct {
	AuthParam

	ID string `path:"id" doc:"meeting id"`
}

type GetMeetingResponse struct {
	Body struct {
		Meeting MeetingResponse `json:"meeting"`
	}
}

func (h *MeetingHandler) GetMeeting(ctx context.Context, input *GetMeetingRequest) (*GetMeetingResponse, error) {
	meetingID := input.ID

	meetingRes, err := h.meetingService.GetMeeting(ctx, meetingID)
	if err != nil {
		switch {
		case errors.Is(err, meeting.ErrMeetingNotFound):
			return nil, huma.Error404NotFound("meeting not found", err)
		default:
			return nil, huma.Error500InternalServerError("an error occured", err)
		}
	}

	resp := &GetMeetingResponse{}
	resp.Body.Meeting = meetingToResponse(meetingRes)
	return resp, nil
}

type EndMeetingRequest struct {
	AuthParam

	ID string `path:"id" doc:"meeting id"`
}

func (h *MeetingHandler) EndMeeting(ctx context.Context, input *EndMeetingRequest) (*struct{}, error) {
	userID, err := getUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	meetingID := input.ID

	err = h.meetingService.EndMeeting(ctx, meetingID, userID)
	if err != nil {
		switch {
		case errors.Is(err, meeting.ErrMeetingNotFound):
			return nil, huma.Error404NotFound("meeting not found", err)
		case errors.Is(err, meeting.ErrNotAuthorized):
			return nil, huma.Error403Forbidden("cannot end meeting, only host can end meeting", err)
		default:
			return nil, huma.Error500InternalServerError("an error occured", err)
		}
	}

	return &struct{}{}, nil
}

type UserDisplayName struct {
	DisplayName string `json:"displayName" doc:"User display name"`
}

type ParticipantResponse struct {
	ID     string           `json:"id" doc:"Participant record ID"`
	UserID string           `json:"userId" doc:"User ID"`
	Role   string           `json:"role" doc:"Participant role (host, co-host, participant)"`
	User   *UserDisplayName `json:"user" doc:"User details"`
}

type GetParticipantsRequest struct {
	AuthParam

	ID string `path:"id" doc:"meeting id"`
}

type GetParticipantsResponse struct {
	Body struct {
		Participants []ParticipantResponse `json:"participants"`
	}
}

func (h *MeetingHandler) GetParticipants(ctx context.Context, input *GetParticipantsRequest) (*GetParticipantsResponse, error) {
	meetingID := input.ID

	participants, err := h.meetingService.GetMeetingParticipants(ctx, meetingID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to get participants", err)
	}

	response := make([]ParticipantResponse, len(participants))
	for i, p := range participants {
		response[i] = ParticipantResponse{
			ID:     p.ID,
			UserID: p.UserID,
			Role:   string(p.Role),
			User: &UserDisplayName{
				DisplayName: p.User.DisplayName,
			},
		}
	}

	resp := &GetParticipantsResponse{}
	resp.Body.Participants = response
	return resp, nil
}

type SendChatMessageRequest struct {
	AuthParam

	ID   string `path:"id" doc:"meeting id"`
	Body struct {
		Message string `json:"message" required:"true" doc:"Chat message content" example:"Hello everyone!"`
	}
}

func (h *MeetingHandler) SendChatMessage(ctx context.Context, input *SendChatMessageRequest) (*struct{}, error) {
	userID, err := getUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}
	meetingID := input.ID
	message := input.Body.Message

	err = h.meetingService.SaveChatMessage(ctx, meetingID, userID, message)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to send message", err)
	}

	return &struct{}{}, nil
}

type ChatMessageResponse struct {
	ID        string           `json:"id" doc:"Message unique identifier"`
	MeetingID string           `json:"meetingId" doc:"ID of the meeting this message belongs to"`
	UserID    string           `json:"userId" doc:"ID of the user who sent the message"`
	Message   string           `json:"message" doc:"Message content"`
	SentAt    time.Time        `json:"sentAt" doc:"When the message was sent"`
	User      *UserDisplayName `json:"user" doc:"User details"`
}

type GetChatMessagesRequest struct {
	AuthParam

	ID string `path:"id" doc:"meeting id"`
}

type GetChatMessagesResponse struct {
	Body struct {
		ChatMessages []ChatMessageResponse `json:"messages" doc:"messages sent during the meeting"`
	}
}

func (h *MeetingHandler) GetChatMessages(ctx context.Context, input *GetChatMessagesRequest) (*GetChatMessagesResponse, error) {
	meetingID := input.ID

	messages, err := h.meetingService.GetChatHistory(ctx, meetingID)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to get chat messages", err)
	}

	response := make([]ChatMessageResponse, len(messages))
	for i, msg := range messages {
		response[i] = ChatMessageResponse{
			ID:        msg.ID,
			MeetingID: msg.MeetingID,
			UserID:    msg.UserID,
			Message:   msg.Message,
			SentAt:    msg.SentAt,
			User: &UserDisplayName{
				DisplayName: msg.User.DisplayName,
			},
		}
	}

	resp := &GetChatMessagesResponse{}
	return resp, nil
}

func meetingToResponse(meeting *models.Meeting) MeetingResponse {
	response := MeetingResponse{
		ID:          meeting.ID,
		Title:       meeting.Title,
		HostID:      meeting.HostID,
		MeetingCode: meeting.MeetingCode,
		IsPrivate:   meeting.IsPrivate,
		CreatedAt:   meeting.CreatedAt,
	}

	if meeting.Host != nil {
		response.Host = &UserInfoSmall{
			ID:          meeting.Host.ID,
			DisplayName: meeting.Host.DisplayName,
		}
	}

	return response
}

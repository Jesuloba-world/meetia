package api

import (
	"github.com/danielgtaylor/huma/v2"

	"github.com/meetia/backend/internal/api/handler"
	"github.com/meetia/backend/internal/services/auth"
	"github.com/meetia/backend/internal/services/meeting"
	"github.com/meetia/backend/internal/services/webrtc"
)

func SetupRoutes(
	api huma.API,
	authService *auth.AuthService,
	sfuService *webrtc.SFUService,
	meetingService *meeting.MeetingService,
) {
	webrtcHandler := handler.NewWebRTCHandler(sfuService, authService.GetTokenAuth())
	meetingHandler := handler.NewMeetinghandler(meetingService, authService.GetTokenAuth())
	authHandler := handler.NewAuthHandler(authService)

	authHandler.RegisterRoutes(api)
	webrtcHandler.RegisterRoutes(api)
	meetingHandler.RegisterRoutes(api)
}

package meeting

import (
	"context"
	"errors"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/meetia/backend/internal/models"
	"github.com/meetia/backend/internal/repository"
)

var (
	ErrMeetingNotFound = errors.New("meeting not found")
	ErrNotAuthorized   = errors.New("not authorized to access this meeting")
	ErrInvalidPassword = errors.New("invalid meeting password")
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type MeetingService struct {
	meetingRepo *repository.MeetingRepository
	userRepo    *repository.UserRepository
}

func NewMeetingService(meetingRepo *repository.MeetingRepository, userRepo *repository.UserRepository) *MeetingService {
	return &MeetingService{
		meetingRepo: meetingRepo,
		userRepo:    userRepo,
	}
}

func (s *MeetingService) CreateMeeting(ctx context.Context, title string, hostID string, isPrivate bool, password string) (*models.Meeting, error) {
	// generate unique meeting code
	meetingCode := generateMeetingCode(10)

	meeting := &models.Meeting{
		Title:       title,
		HostID:      hostID,
		MeetingCode: meetingCode,
		IsPrivate:   isPrivate,
		Password:    password,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.meetingRepo.Create(ctx, meeting); err != nil {
		return nil, err
	}

	// add host as participant with host role
	participant := &models.MeetingParticipant{
		MeetingID: meeting.ID,
		UserID:    hostID,
		Role:      "host",
	}

	if err := s.meetingRepo.AddParticipant(ctx, participant); err != nil {
		return nil, err
	}

	return meeting, nil
}

func (s *MeetingService) JoinMeeting(ctx context.Context, meetingCode string, userID string, password string) (*models.Meeting, error) {
	meeting, err := s.meetingRepo.GetByCode(ctx, meetingCode)
	if err != nil {
		return nil, ErrMeetingNotFound
	}

	if meeting.IsPrivate && meeting.Password != password {
		return nil, ErrInvalidPassword
	}

	// check if user is already a participant
	participants, err := s.meetingRepo.GetParticipants(ctx, meeting.ID)
	if err != nil {
		return nil, err
	}

	var isParticipant bool
	for _, p := range participants {
		if p.UserID == userID {
			isParticipant = true
			// update joined_at if they're rejoining
			if p.LeftAt.After(time.Time{}) {
				p.JoinedAt = time.Now()
				p.LeftAt = time.Time{}
				if err := s.meetingRepo.UpdateParticipant(ctx, p); err != nil {
					return nil, err
				}
			}
			break
		}
	}

	// if not participant, add then
	if !isParticipant {
		participant := models.MeetingParticipant{
			MeetingID: meeting.ID,
			UserID:    userID,
			Role:      models.MeetingParticipantNormal,
			JoinedAt:  time.Now(),
		}
		if err := s.meetingRepo.AddParticipant(ctx, &participant); err != nil {
			return nil, err
		}
	}

	return meeting, nil
}

func (s *MeetingService) EndMeeting(ctx context.Context, meetingID string, userID string) error {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return ErrMeetingNotFound
	}

	// Only host can end meeting
	if meeting.HostID != userID {
		return ErrNotAuthorized
	}

	return s.meetingRepo.EndMeeting(ctx, meetingID)
}

func (s *MeetingService) GetMeeting(ctx context.Context, meetingID string) (*models.Meeting, error) {
	return s.meetingRepo.GetByID(ctx, meetingID)
}

func (s *MeetingService) GetMeetingByCode(ctx context.Context, code string) (*models.Meeting, error) {
	return s.meetingRepo.GetByCode(ctx, code)
}

func (s *MeetingService) GetActiveUserMeetings(ctx context.Context, userID string) ([]*models.Meeting, error) {
	return s.meetingRepo.GetActiveForUser(ctx, userID)
}

func (s *MeetingService) GetMeetingParticipants(ctx context.Context, meetingID string) ([]*models.MeetingParticipant, error) {
	return s.meetingRepo.GetParticipants(ctx, meetingID)
}

func (s *MeetingService) SaveChatMessage(ctx context.Context, meetingID string, userID string, message string) error {
	chat := &models.MeetingChat{
		MeetingID: meetingID,
		UserID:    userID,
		Message:   message,
		SentAt:    time.Now(),
	}
	return s.meetingRepo.SaveChat(ctx, chat)
}

func (s *MeetingService) GetChatHistory(ctx context.Context, meetingID string) ([]*models.MeetingChat, error) {
	return s.meetingRepo.GetChatHistory(ctx, meetingID)
}

func generateMeetingCode(length int) string {
	return gonanoid.MustGenerate(charset, length)
}

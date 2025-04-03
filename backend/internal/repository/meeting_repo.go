package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/meetia/backend/internal/models"
)

type MeetingRepository struct {
	db *bun.DB
}

func NewMeetingRepository(db *bun.DB) *MeetingRepository {
	return &MeetingRepository{db: db}
}

func (r *MeetingRepository) Create(ctx context.Context, meeting *models.Meeting) error {
	_, err := r.db.NewInsert().Model(meeting).Exec(ctx)
	return err
}

func (r *MeetingRepository) GetByID(ctx context.Context, id string) (*models.Meeting, error) {
	meeting := new(models.Meeting)
	err := r.db.NewSelect().
		Model(meeting).
		Relation("Host").
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return meeting, nil
}

func (r *MeetingRepository) GetByCode(ctx context.Context, code string) (*models.Meeting, error) {
	meeting := new(models.Meeting)
	err := r.db.NewSelect().
		Model(meeting).
		Relation("Host").
		Where("meeting_code = ?", code).
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return meeting, nil
}

func (r *MeetingRepository) GetActiveForUser(ctx context.Context, userID string) ([]*models.Meeting, error) {
	var meetings []*models.Meeting
	err := r.db.NewSelect().
		Model(&meetings).
		Relation("Host").
		Where("host_id = ? OR id IN (SELECT meeting_id FROM meeting_participants WHERE user_id = ?)", userID, userID).
		Where("ended_at IS NULL").
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return meetings, nil
}

func (r *MeetingRepository) Update(ctx context.Context, meeting *models.Meeting) error {
	meeting.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(meeting).Where("id = ?", meeting.ID).Exec(ctx)
	return err
}

func (r *MeetingRepository) EndMeeting(ctx context.Context, id string) error {
	_, err := r.db.NewUpdate().
		Model((*models.Meeting)(nil)).
		Set("ended_at = ?", time.Now()).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Where("ended_at IS NULL").
		Exec(ctx)
	return err
}

func (r *MeetingRepository) AddParticipant(ctx context.Context, participant *models.MeetingParticipant) error {
	_, err := r.db.NewInsert().Model(participant).Exec(ctx)
	return err
}

func (r *MeetingRepository) UpdateParticipant(ctx context.Context, participant *models.MeetingParticipant) error {
	_, err := r.db.NewUpdate().Model(participant).Where("id = ?", participant.ID).Exec(ctx)
	return err
}

func (r *MeetingRepository) GetParticipants(ctx context.Context, meetingID string) ([]*models.MeetingParticipant, error) {
	var participants []*models.MeetingParticipant
	err := r.db.NewSelect().
		Model(&participants).
		Relation("User").
		Where("meeting_id = ?", meetingID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *MeetingRepository) SaveChat(ctx context.Context, chat *models.MeetingChat) error {
	_, err := r.db.NewInsert().Model(chat).Exec(ctx)
	return err
}

func (r *MeetingRepository) GetChatHistory(ctx context.Context, meetingID string) ([]*models.MeetingChat, error) {
	var chats []*models.MeetingChat
	err := r.db.NewSelect().
		Model(&chats).
		Relation("User").
		Where("meeting_id = ?", meetingID).
		Order("sent_at ASC").
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return chats, nil
}

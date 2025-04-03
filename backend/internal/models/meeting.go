package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Meeting struct {
	bun.BaseModel `bun:"table:meetings,alias:m"`

	ID          string    `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Title       string    `bun:"title,notnull" json:"title"`
	HostID      string    `bun:"host_id,notnull" json:"hostId"`
	MeetingCode string    `bun:"meeting_code,notnull,unique" json:"meetingCode"`
	Password    string    `bun:"password" json:"password,omitempty"`
	IsPrivate   bool      `bun:"is_private,notnull" json:"isPrivate"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updatedAt"`
	ScheduledAt time.Time `bun:"scheduled_at" json:"scheduledAt,omitempty"`
	EndedAt     time.Time `bun:"ended_at" json:"endedAt,omitempty"`

	// Relations
	Host         *User                 `bun:"rel:belongs-to,join:host_id=id" json:"host,omitempty"`
	Participants []*MeetingParticipant `bun:"rel:has-many,join:id=meeting_id" json:"participants,omitempty"`
}

type MeetingParticipant struct {
	bun.BaseModel `bun:"table:meeting_participants,alias:mp"`

	ID        string    `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	MeetingID string    `bun:"meeting_id,notnull" json:"meetingId"`
	UserID    string    `bun:"user_id,notnull" json:"userId"`
	Role      string    `bun:"role,notnull" json:"role"` // host, co-host, participant
	JoinedAt  time.Time `bun:"joined_at" json:"joinedAt,omitempty"`
	LeftAt    time.Time `bun:"left_at" json:"leftAt,omitempty"`

	// Relations
	Meeting *Meeting `bun:"rel:belongs-to,join:meeting_id=id" json:"meeting,omitempty"`
	User    *User    `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}

type MeetingChat struct {
	bun.BaseModel `bun:"table:meeting_chats,alias:mc"`

	ID        string    `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	MeetingID string    `bun:"meeting_id,notnull" json:"meetingId"`
	UserID    string    `bun:"user_id,notnull" json:"userId"`
	Message   string    `bun:"message,notnull" json:"message"`
	SentAt    time.Time `bun:"sent_at,notnull,default:current_timestamp" json:"sentAt"`

	// Relations
	Meeting *Meeting `bun:"rel:belongs-to,join:meeting_id=id" json:"meeting,omitempty"`
	User    *User    `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}

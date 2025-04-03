package models

import (
	"time"

	"github.com/uptrace/bun"

)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           string    `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Email        string    `bun:"email,notnull,unique" json:"email"`
	PasswordHash string    `bun:"password_hash,notnull" json:"-"`
	DisplayName  string    `bun:"display_name,notnull" json:"displayName"`
	CreatedAt    time.Time `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt    time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updatedAt"`
}

package jwtsession

import (
	"time"

	"github.com/google/uuid"
)

type JwtSession struct {
	ID            uuid.UUID
	UserEmail     string
	ResfreshToken string
	AccessTokenId uuid.UUID
	IsRevoked     bool
	ExpiresAt     time.Time
}

func (s *JwtSession) ToMap() map[string]string {
	return map[string]string{
		"id":           s.ID.String(),
		"userEmail":    s.UserEmail,
		"refreshToken": s.ResfreshToken,
		"expiresAt":    s.ExpiresAt.String(),
	}
}

func MakeSlice(sessions []*JwtSession) []map[string]string {
	var sessionsMap []map[string]string

	for _, session := range sessions {
		sessionsMap = append(sessionsMap, session.ToMap())
	}

	return sessionsMap
}

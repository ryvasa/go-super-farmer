package domain

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	ReportID  string    `json:"report_id"`
	Email     uuid.UUID `json:"email"`
}

package models

import (
	"bytes"
	"time"
)

type Image struct {
	ID        int64
	Name      string
	OwnerID   int64
	OwnerType string
	Buf       *bytes.Buffer
	CreatedAt time.Time
	UpdatedAt time.Time
}

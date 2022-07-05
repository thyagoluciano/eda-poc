package domain

import "github.com/gofrs/uuid"

func GenerateId() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

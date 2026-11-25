package service

import (
	"deller/internal/dto"
	"deller/internal/schema"
)

func ConvertMails(m []schema.Mails) dto.ContactsResponse {
	mails := make([]dto.ContactMail, len(m))
	for i, e := range m {
		mails[i] = dto.ContactMail{
			Id:          e.Id,
			Name:        e.Name,
			PhoneNumber: e.PhoneNumber,
			Email:       e.Email,
			Comment:     e.Comment,
			IsBanned:    e.IsBanned,
			CreatedAt:   e.CreatedAt,
		}
	}
	return dto.ContactsResponse{
		Mails: mails,
		Total: len(m),
	}
}

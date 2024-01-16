package service

import (
	"context"
)

// GetDomainByEmail get the user related domain by the given email
func (s *Service) GetDomainByEmail(ctx context.Context, email string) (string, error) {
	domain, err := s.db.GetDomainByEmail(ctx, email)

	return domain, err
}

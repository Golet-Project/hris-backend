package service

import "context"

func (s *Service) GetDomainByUid(ctx context.Context, uid string) (string, error) {
	if uid == "" {
		return "", nil
	}

	domain, err := s.db.GetDomainByUid(ctx, uid)

	return domain, err
}

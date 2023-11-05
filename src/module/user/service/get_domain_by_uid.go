package service

import "context"

func (s *Service) GetDomainByUid(ctx context.Context, uid string) (string, error) {
	domain, err := s.db.GetDomainByUid(ctx, uid)

	return domain, err
}

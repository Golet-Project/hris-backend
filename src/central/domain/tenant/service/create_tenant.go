package service

import (
	"context"
	"hroost/central/domain/tenant/model"
	"hroost/shared/primitive"
	"net/http"
)

type CreateTenantIn struct {
	Name   string
	Domain string
}

type CreateTenantOut struct {
	primitive.CommonResult

	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type CreateTenantQueue interface {
	MigrateTenantDB(ctx context.Context, data model.MigrateTenantDBIn) (err *primitive.RepoError)
}

type CreateTenantDb interface {
	CountTenantByDomain(ctx context.Context, domain string) (out model.CountTenantByDomainOut, err *primitive.RepoError)
	CreateTenant(ctx context.Context, data model.CreateTenantIn) (tenant model.CreateTenantOut, err *primitive.RepoError)
}

type CreateTenant struct {
	Db    CreateTenantDb
	Queue CreateTenantQueue
}

// CreateTenant send event to MQ and run tenant database migration
func (s *CreateTenant) Exec(ctx context.Context, in CreateTenantIn) (out CreateTenantOut) {
	// validate request
	if err := s.ValidateCreateTenantIn(in); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// check if given domain already exists
	tenantCount, repoError := s.Db.CountTenantByDomain(ctx, in.Domain)
	if repoError != nil {
		if repoError.Issue != primitive.RepoErrorCodeDataNotFound {
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}
	if tenantCount.Count > 0 {
		out.SetResponse(http.StatusBadRequest, "tenant domain already registered")
		return
	}

	// create tenant
	createdTenant, repoError := s.Db.CreateTenant(ctx, model.CreateTenantIn{
		Name:   in.Name,
		Domain: in.Domain,
	})
	if repoError != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
		return
	}

	// send to MQ and run migration
	repoError = s.Queue.MigrateTenantDB(ctx, model.MigrateTenantDBIn{
		Domain: in.Domain,
	})
	if repoError != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
		return
	}

	out.ID = createdTenant.UID
	out.Name = createdTenant.Name
	out.Domain = createdTenant.Domain

	out.SetResponse(http.StatusCreated, "tenant created")
	return
}

// ValidateCreateTenantIn validate the request body
func (s *CreateTenant) ValidateCreateTenantIn(in CreateTenantIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate name
	if len(in.Name) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "name",
			Message: "name is required",
		})
	} else {
		if len(in.Name) > 100 {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeTooLong,
				Field:   "name",
				Message: "name must less than 100 characters length",
			})
		}
	}

	// validate domain
	if len(in.Domain) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "domain",
			Message: "domain is required",
		})
	} else {
		if len(in.Domain) > 50 {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeTooLong,
				Field:   "domain",
				Message: "domain must be less than 50 characters length",
			})
		}

		if !primitive.TenantDomainPattern.MatchString(in.Domain) {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "domain",
				Message: "invalid tenant domain",
			})
		}
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}

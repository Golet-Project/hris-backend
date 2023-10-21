package service

import (
	"context"
	"errors"
	"hris/module/shared/primitive"
	"hris/module/tenant/repo/queue"
	"hris/module/tenant/repo/tenant"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type Internal_CreateTenantIn struct {
	Name   string
	Domain string
}

type Internal_CreateTenantOut struct {
	primitive.CommonResult

	ID     string
	Name   string
	Domain string
}

func ValidateInternal_CreateTenantIn(in Internal_CreateTenantIn) *primitive.RequestValidationError {
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

func (s *Internal_TenantService) CreateTenant(ctx context.Context, in Internal_CreateTenantIn) (out Internal_CreateTenantOut) {
	// validate request
	if err := ValidateInternal_CreateTenantIn(in); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// check if given domain already exists
	tenantCount, err := s.TenantRepo.CountTenantByDomain(ctx, in.Domain)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}
	if tenantCount.Count > 0 {
		out.SetResponse(http.StatusBadRequest, "tenant domain already registered")
		return
	}

	// create tenant
	createdTenant, err := s.TenantRepo.Internal_CreateTenant(ctx, tenant.Internal_CreateTenantIn{
		Name:   in.Name,
		Domain: in.Domain,
	})
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	// send to MQ and run migration
	err = s.QueueRepo.MigrateTenantDB(ctx, queue.MigrateTenantDBIn{
		Domain: in.Domain,
	})
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	out.ID = createdTenant.UID
	out.Name = createdTenant.Name
	out.Domain = createdTenant.Domain

	out.SetResponse(http.StatusCreated, "tenant created")
	return
}

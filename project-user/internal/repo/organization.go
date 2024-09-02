package repo

import (
	"context"
	"ms_project/project-user/internal/model"
)

type OrganizationRepo interface {
	SaveOrganization(ctx context.Context, organization *model.Organization) error
	FindOrganizationByMemberId(ctx context.Context, id int64) ([]*model.Organization, error)
}

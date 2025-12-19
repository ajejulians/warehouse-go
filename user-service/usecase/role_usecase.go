package usecase

import (
	"context"

	"warehouse-go/user-service/model"
	"warehouse-go/user-service/repository"
)

type RoleUsecaseInterface interface {
	CreateRole(ctx context.Context, role model.Role) error
	UpdateRole(ctx context.Context, role model.Role) error
	DeleteRole(ctx context.Context, id uint) error
	GetAllRoleByID(ctx context.Context, id uint) (*model.Role, error)
	GetAllRoles(ctx context.Context) ([]model.Role, error)
}

type roleUsecase struct {
	roleRepo repository.RoleRepositoryInterface
}


// ────────────────────────────────────────────────────────────────
// CreateRole implements Roleusecase Interface
// ────────────────────────────────────────────────────────────────
func (r *roleUsecase) CreateRole(ctx context.Context, role model.Role) error {
	return r.roleRepo.CreateRole(ctx, role)
}

// ────────────────────────────────────────────────────────────────
// DeleteRole implements Roleusecase Interface
// ────────────────────────────────────────────────────────────────
func (r *roleUsecase) DeleteRole(ctx context.Context,  id uint) error {
	return r.roleRepo.DeleteRole(ctx, id)
}

// ────────────────────────────────────────────────────────────────
// GetAllRoles implements Roleusecase Interface
// ────────────────────────────────────────────────────────────────
func (r *roleUsecase) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	return r.roleRepo.GetAllRoles(ctx)
}

// ────────────────────────────────────────────────────────────────
// GetAllRoleByID implements Roleusecase Interface
// ────────────────────────────────────────────────────────────────
func (r *roleUsecase) GetAllRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	return r.roleRepo.GetAllRoleByID(ctx, id)
}
// ────────────────────────────────────────────────────────────────
// UpdateRole implements Roleusecase Interface
// ────────────────────────────────────────────────────────────────
func (r *roleUsecase) UpdateRole(ctx context.Context, role model.Role) error {
	return r.roleRepo.UpdateRole(ctx, role)
}


func NewRoleUsecase(roleRepo repository.RoleRepositoryInterface) RoleUsecaseInterface {
	return &roleUsecase{roleRepo: roleRepo}
}

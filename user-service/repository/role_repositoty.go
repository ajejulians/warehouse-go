package repository

import (
	"context"
	"warehouse-go/user-service/model"
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	CreateRole(ctx context.Context, role model.Role) error
	UpdateRole(ctx context.Context, role model.Role) error
	DeleteRole(ctx context.Context, id uint) error
	GetAllRoleByID(ctx context.Context, id uint) (*model.Role, error)
	GetAllRoles(ctx context.Context) ([]model.Role, error)
}

type roleRepository struct{
	db *gorm.DB
}

// ────────────────────────────────────────────────────────────────
// CreateRole implements RoleRepository Inteface
// ────────────────────────────────────────────────────────────────
func (r *roleRepository) CreateRole(ctx context.Context, role model.Role) error {
	select {
	case <-ctx.Done():
		log.Errorf("[RoleRepository] CreateRole - 1 %v", ctx.Err())
		return ctx.Err()
	default:
		return r.db.WithContext(ctx).Create(&role).Error
	}
}

// ────────────────────────────────────────────────────────────────
// DeleteRole implements RoleRepository Inteface
// ────────────────────────────────────────────────────────────────
func (r *roleRepository) DeleteRole(ctx context.Context, id uint)  error {
	select {
	case <-ctx.Done():
		log.Errorf("[RoleRepository] DeleteRole - 1 %v", ctx.Err())
		return ctx.Err()
	default:
		modelRole := model.Role{}
		if err := r.db.WithContext(ctx).Preload("Users").Where("id = ?", id).First(&modelRole, id).Error; err != nil {
			log.Errorf("[RoleRepository] DeleteRole - 2: %v", err)
			return err
		}

		if len(modelRole.Users) > 0 {
			log.Errorf("[RoleRepository] DeleteRole - 3: %v", "Role has Users")
			return errors.New("Role has Users")
		}  

		return r.db.WithContext(ctx).Delete(&modelRole).Error
	}
}

// ────────────────────────────────────────────────────────────────
// GetAllRoles implements RoleRepository Inteface
// ────────────────────────────────────────────────────────────────

func (r *roleRepository) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[RoleRepository] GetAllRoles - 1 %v", ctx.Err())
		return nil, ctx.Err()
	default:
		modelRoles := []model.Role{}
		err := r.db.WithContext(ctx).Preload("Users").Find(&modelRoles).Error 
		if err != nil {
			log.Errorf("[RoleRepository] GetAllRoles - 2: %v", err)
			return nil, err
		}

		return modelRoles, nil
	}
}

// ────────────────────────────────────────────────────────────────
// GetAllRoleByID implements RoleRepository Inteface
// ────────────────────────────────────────────────────────────────
func (r *roleRepository) GetAllRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	select {
	case <-ctx.Done():
		log.Errorf("[RoleRepository] GetAllRoleByID - 1 %v", ctx.Err())
		return nil, ctx.Err()
	default:
		modelRole := model.Role{}
		if err := r.db.WithContext(ctx).Preload("Users").Where("id = ?", id).First(&modelRole).Error; err != nil {
			log.Errorf("[RoleRepository] GetAllRoleByID - 2: %v", err)
			return nil, err
		}

		return &modelRole, nil
	}

}

// ────────────────────────────────────────────────────────────────
// UpdateRole implements RoleRepository Inteface
// ────────────────────────────────────────────────────────────────
func (r *roleRepository) UpdateRole(ctx context.Context, role model.Role) error {
	select {
	case <-ctx.Done():
		log.Errorf("[RoleRepository] UpdateRole - 1 %v", ctx.Err())
		return ctx.Err()
	default:
		modelRole := model.Role{}
		if err := r.db.WithContext(ctx).Preload("Users").Where("id = ?", role.ID).First(&modelRole).Error; err != nil {
			log.Errorf("[RoleRepository] UpdateRole - 2: %v", err)
			return err
		}

		modelRole.Name = role.Name
		return r.db.WithContext(ctx).Save(&modelRole).Error
	}
}

func NewRoleRepository(db *gorm.DB) RoleRepositoryInterface {
	return &roleRepository{db: db}
}

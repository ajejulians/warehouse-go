package usecase

import (
	"context"
	"warehouse-go/user-service/model"
	"warehouse-go/user-service/pkg/conv"
	"warehouse-go/user-service/repository"
	"warehouse-go/user-service/service"

	"github.com/gofiber/fiber/v2/log"
)

type UserUsecaseInterface interface {
	CreateUser(ctx context.Context, user model.User) error
	GetAllUsers(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.User, int64, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, id uint) error

	GetUserByRoleName(ctx context.Context, roleName string) ([]model.User, error)

	AssignUserToRole(ctx context.Context, userID uint, roleID uint) error
	EditAssignUserToRole(ctx context.Context, assignRoleID uint, userID uint, roleID uint) error
	GetUserRoleByID(ctx context.Context, assignRoleID uint) (*model.UserRole, error)
	GetAllUserRoles(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.UserRole, int64, error)
}

type userUsecase struct {
	userRepo        repository.UserRepositoryInterface
	rabbitMQService service.RabbitMQServiceInterface
}

// GetUserRoleByID implements UserUsecaseInterface.
func (u *userUsecase) GetUserRoleByID(ctx context.Context, assignRoleID uint) (*model.UserRole, error) {
	return u.userRepo.GetUserByRoleID(ctx, assignRoleID)
}

// ────────────────────────────────────────────────────────────────
// AssignUserToRole implements UserUsecaseInterface
// ────────────────────────────────────────────────────────────────
func (u *userUsecase) AssignUserToRole(ctx context.Context, userID uint, roleID uint) error {
	return u.userRepo.AssignUserToRole(ctx, userID, roleID)
}

// ────────────────────────────────────────────────────────────────
// CreateUser implements UserUsecaseInteface
// ────────────────────────────────────────────────────────────────
func (u *userUsecase) CreateUser(ctx context.Context, user model.User) error {
	password, err := conv.HashPassword(user.Password)
	if err != nil {
		log.Errorf("[UserUsecase] CreateUser - 1: %v", err)
		return err
	}

	uncryptedPassword := user.Password
	user.Password = password
	result, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		log.Errorf("[UserUsecase] CreateUser - 2: %v", err)
		return err
	}

	emailPayload := service.EmailPayload{
		Email:    result.Email,
		Password: uncryptedPassword,
		Type:     "welcome_email",
		UserID:   uint(result.ID),
		Name:     result.Name,
	}

	go func() {
		if err := u.rabbitMQService.PublishEmail(ctx, emailPayload); err != nil {
			log.Errorf("[UserUsecase] CreateUser - 3: %v", err)
		}
	}()

	return nil
}

// ────────────────────────────────────────────────────────────────
// DeleteUser Implements UserRepository Inteface
// ────────────────────────────────────────────────────────────────
func (u *userUsecase) DeleteUser(ctx context.Context, id uint) error {
	_, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		log.Errorf("[UserUsecase] DeleteUser - 1: %v", err)
		return err
	}

	if err := u.userRepo.DeleteUser(ctx, id); err != nil {
		log.Errorf("[UserUsecase] DeleteUser - 2: %v", err)
		return err
	}

	return nil
}

// ────────────────────────────────────────────────────────────────
// EditAssignUserToRole Implements UserRepository Interface
// ────────────────────────────────────────────────────────────────
func (u *userUsecase) EditAssignUserToRole(ctx context.Context, assignRoleID uint, userID uint, roleID uint) error {
	return u.userRepo.EditAssignUserToRole(ctx, assignRoleID, userID, roleID)
}

// ────────────────────────────────────────────────────────────────
// GetAllUserRoles Implements UserRepository Interface
// ────────────────────────────────────────────────────────────────
func (u *userUsecase) GetAllUserRoles(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.UserRole, int64, error) {
	return u.userRepo.GetAllUserRoles(ctx, page, limit, search, sortBy, sortOrder)
}

// ────────────────────────────────────────────────────────────────
// GetAllUsers Implements UserRepository Interface
// ────────────────────────────────────────────────────────────────
func (u *userUsecase) GetAllUsers(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	//Get Users fromr repository
	users, totalRecords, err := u.userRepo.GetAllUsers(ctx, page, limit, search, sortBy, sortOrder)
	if err != nil {
		log.Errorf("[UserUsecase] GetAllUsers - 1: %v", err)
		return nil, 0, err
	}
	return users, totalRecords, nil
}

// ────────────────────────────────────────────────────────────────
// GetUserbyEmail Implements UserRepository Interface
// ────────────────────────────────────────────────────────────────
func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return u.userRepo.GetUserByEmail(ctx, email)
}

// ────────────────────────────────────────────────────────────────
// GetUserbyID Implements UserRepository Interface
// ────────────────────────────────────────────────────────────────
func (u *userUsecase) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}

// ────────────────────────────────────────────────────────────────
// GetUserbyRoleName Implements UserRepository Interface
// ────────────────────────────────────────────────────────────────
func (u *userUsecase) GetUserByRoleName(ctx context.Context, roleName string) ([]model.User, error) {
	return u.userRepo.GetUserByRoleName(ctx, roleName)
}

// ────────────────────────────────────────────────────────────────
// UpdateUser
// ────────────────────────────────────────────────────────────────
func (u *userUsecase) UpdateUser(ctx context.Context, user model.User) error {
	if err := u.userRepo.UpdateUser(ctx, user); err != nil {
		log.Errorf("[UserUsecase] UpdateUser - 1: %v", err)
		return err
	}

	return nil
}

func NewUserUsecase(userRepo repository.UserRepositoryInterface, rabbitMQService service.RabbitMQServiceInterface) UserUsecaseInterface {
	return &userUsecase{userRepo: userRepo, rabbitMQService: rabbitMQService}
}

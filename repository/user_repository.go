package repository

import (
	"context"

	"github.com/rafli-lutfi/go-commerce/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, userID int) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	JoinTableUser(ctx context.Context, userID int) (models.UserInfo, error)
	CreateNewUser(ctx context.Context, user models.NewUser) (int, error)
	AddMoreAddress(ctx context.Context, oldUserData models.User, newUserAddress models.NewAddress) error
	UpdateUser(ctx context.Context, newUserData *models.User) error
	UpdateAddress(ctx context.Context, newUserAddress *models.UserAddress) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserByID(ctx context.Context, userID int) (models.User, error) {
	var user models.User

	err := ur.db.WithContext(ctx).Where("id = ?", userID).Find(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	err := ur.db.WithContext(ctx).Where("email = ?", email).Find(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (ur *userRepository) JoinTableUser(ctx context.Context, userID int) (models.UserInfo, error) {
	var user models.User

	err := ur.db.WithContext(ctx).Preload(clause.Associations).Where("id = ?", userID).Find(&user).Error
	if err != nil {
		return models.UserInfo{}, err
	}

	userInfo := models.UserInfo{
		ID:           user.ID,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		MobileNumber: user.MobileNumber,
		Addresses:    user.UserAddresses,
	}

	return userInfo, nil
}

func (ur *userRepository) CreateNewUser(ctx context.Context, user models.NewUser) (int, error) {
	newUser := models.User{
		Email:        user.Email,
		Username:     user.Username,
		Password:     user.Password,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		MobileNumber: user.MobileNumber,
		UserAddresses: []models.UserAddress{
			{
				Address:    user.NewAddress.Address,
				City:       user.NewAddress.City,
				PostalCode: user.NewAddress.PostalCode,
				Country:    user.NewAddress.Country,
			},
		},
	}

	// create user record
	err := ur.db.WithContext(ctx).CreateInBatches(&newUser, 100).Error
	if err != nil {
		ur.db.Exec("SELECT SETVAL('users_id_seq', (SELECT MAX(id) FROM users));") // Reset autoincrement
		return 0, err
	}

	return int(newUser.ID), nil
}

func (ur *userRepository) AddMoreAddress(ctx context.Context, oldUserData models.User, newUserAddress models.NewAddress) error {
	address := models.UserAddress{
		UserID:     newUserAddress.UserID,
		Address:    newUserAddress.Address,
		City:       newUserAddress.City,
		PostalCode: newUserAddress.PostalCode,
		Country:    newUserAddress.Country,
	}

	err := ur.db.WithContext(ctx).Model(&oldUserData).Association("UserAddresses").Append(&address)
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, newUserData *models.User) error {
	// update user
	err := ur.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", newUserData.ID).Updates(newUserData).Error
	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) UpdateAddress(ctx context.Context, newUserAddress *models.UserAddress) error {
	err := ur.db.WithContext(ctx).Model(&models.UserAddress{}).Where("id = ?", newUserAddress.ID).Omit("user_id").Updates(newUserAddress).Error
	if err != nil {
		return err
	}

	return nil
}

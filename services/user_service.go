package services

import (
	"context"
	"errors"

	"github.com/rafli-lutfi/go-commerce/models"
	"github.com/rafli-lutfi/go-commerce/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID int) (models.UserInfo, error)
	Login(ctx context.Context, creds *models.Login) (int, error)
	Register(ctx context.Context, user models.NewUser) (int, error)
	AddNewAddress(ctx context.Context, newAddress models.NewAddress) error
	UpdateUser(ctx context.Context, user *models.User) error
	UpdateUserAddress(ctx context.Context, userAddress *models.UserAddress) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userService {
	return &userService{userRepository}
}

func hashPassword(password *string) {
	bytePassword := []byte(*password)
	hash, _ := bcrypt.GenerateFromPassword(bytePassword, 12)

	*password = string(hash)
}

func checkPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (us *userService) GetUserByID(ctx context.Context, userID int) (models.UserInfo, error) {
	userDB, err := us.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return models.UserInfo{}, err
	}
	if userDB.ID == 0 {
		return models.UserInfo{}, errors.New("user record not found")
	}

	user, err := us.userRepository.JoinTableUser(ctx, userID)
	if err != nil {
		return models.UserInfo{}, err
	}

	return user, nil
}

func (us *userService) Login(ctx context.Context, creds *models.Login) (int, error) {
	userDB, err := us.userRepository.GetUserByEmail(ctx, creds.Email)
	if err != nil {
		return 0, err
	}

	if userDB.Email == "" || userDB.ID == 0 {
		return 0, errors.New("email not exist")
	}

	check := checkPassword(userDB.Password, creds.Password)
	if !check {
		return 0, errors.New("password not matched")
	}

	return int(userDB.ID), nil
}

func (us *userService) Register(ctx context.Context, user models.NewUser) (int, error) {
	hashPassword(&user.Password)

	userID, err := us.userRepository.CreateNewUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (us *userService) AddNewAddress(ctx context.Context, newAddress models.NewAddress) error {
	userDB, err := us.userRepository.GetUserByID(ctx, int(newAddress.UserID))
	if err != nil {
		return err
	}

	err = us.userRepository.AddMoreAddress(ctx, userDB, newAddress)
	if err != nil {
		return err
	}

	return nil
}

func (us *userService) UpdateUser(ctx context.Context, user *models.User) error {
	userDB, err := us.userRepository.GetUserByID(ctx, int(user.ID))
	if err != nil {
		return err
	}
	if userDB.ID == 0 {
		return errors.New("user record not found")
	}

	err = us.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (us *userService) UpdateUserAddress(ctx context.Context, userAddress *models.UserAddress) error {
	userDB, err := us.userRepository.GetUserByID(ctx, int(userAddress.UserID))
	if err != nil {
		return err
	}
	if userDB.ID == 0 {
		return errors.New("user record not found")
	}

	err = us.userRepository.UpdateAddress(ctx, userAddress)
	if err != nil {
		return err
	}

	return nil
}

package service

import (
	"context"
	"mods/dto"
	"mods/entity"
	"mods/repository"
	"mods/utils"

	"github.com/mashingan/smapping"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	CreateUser(ctx context.Context, userDTO dto.UserCreateDTO) (entity.User, error)
	IsDuplicateEmail(ctx context.Context, email string) (bool, error)
	VerifyCredential(ctx context.Context, email string, password string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	UploadGame(ctx context.Context, gameDTO dto.UploadGame, userid uint64) (entity.Game, error)
	PurchaseGame(ctx context.Context, gameid uint64, userid uint64, metodeBayar string) (entity.Game, error)
	UserProfile(ctx context.Context, userid uint64) (entity.User, error)
	TopUp(ctx context.Context, userid uint64, nominal uint64) (entity.User, error)
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepository: ur,
	}
}

func (us *userService) CreateUser(ctx context.Context, userDTO dto.UserCreateDTO) (entity.User, error) {
	var user entity.User
	if err := smapping.FillStruct(&user, smapping.MapFields(&userDTO)); err != nil {
		return user, err
	}
	return us.userRepository.InsertUser(ctx, user)
}

func (us *userService) IsDuplicateEmail(ctx context.Context, email string) (bool, error) {
	checkUser, err := us.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if checkUser.Email == "" {
		return false, nil
	}

	return true, nil
}

func (us *userService) VerifyCredential(ctx context.Context, email string, password string) (bool, error) {
	checkUser, err := us.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	checkPassword, err := utils.PasswordCompare(checkUser.Password, []byte(password))
	if err != nil {
		return false, err
	}

	if checkUser.Email == email && checkPassword {
		return true, nil
	}
	return false, nil
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	return us.userRepository.GetUserByEmail(ctx, email)
}

func (us *userService) UploadGame(ctx context.Context, gameDTO dto.UploadGame, userid uint64) (entity.Game, error) {
	return us.userRepository.UploadGame(ctx, gameDTO, userid)
}

func (us *userService) PurchaseGame(ctx context.Context, gameid uint64, userid uint64, metodeBayar string) (entity.Game, error) {
	return us.userRepository.PurchaseGame(ctx, gameid, userid, metodeBayar)
}

func (us *userService) UserProfile(ctx context.Context, userid uint64) (entity.User, error) {
	return us.userRepository.UserProfile(ctx, userid)
}

func (us *userService) TopUp(ctx context.Context, userid uint64, nominal uint64) (entity.User, error) {
	return us.userRepository.TopUp(ctx, userid, nominal)
}

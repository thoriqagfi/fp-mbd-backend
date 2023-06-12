package service

import (
	"context"
	"errors"
	"mods/dto"
	"mods/entity"
	"mods/repository"
	"mods/utils"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	// functional
	CreateUser(ctx context.Context, userDTO dto.UserCreateDTO) (entity.User, error)
	IsDuplicateEmail(ctx context.Context, email string) (bool, error)
	VerifyCredential(ctx context.Context, email string, password string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)

	// profiles
	UserProfile(ctx context.Context, userid uint64) (entity.User, error)
	DeveloperProfile(ctx context.Context, devid uint64) (dto.DeveloperReleases, error)

	// Transactional
	UploadGame(ctx context.Context, gameDTO dto.UploadGame, userid uint64) (entity.Game, error)
	PurchaseGame(ctx context.Context, gameid uint64, userid uint64, metodeBayar string) (entity.Game, error)
	TopUp(ctx context.Context, userid uint64, nominal uint64) (entity.User, error)
	UploadDLC(ctx context.Context, dlc dto.UploadDLC) (entity.DLC, error)
	PurchaseDLC(ctx context.Context, dlcid uint64, userid uint64, metodeBayar string) (entity.DLC, error)

	// Add Tags Languages OS
	AddToGame(nama string, gameID uint64, method string) (any, error)
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepository: ur,
	}
}

func (us *userService) CreateUser(ctx context.Context, userDTO dto.UserCreateDTO) (entity.User, error) {
	newUser := entity.User{
		Name:          userDTO.Name,
		Email:         userDTO.Email,
		Password:      userDTO.Password,
		Profile_image: "https://drive.google.com/uc?export=view&id=1GV5u8MnB88S3Hf92-JLnfpHx6kBaOoBU",
		Role:          userDTO.Role,
		Wallet:        0,
	}

	return us.userRepository.InsertUser(ctx, newUser)
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

func (us *userService) DeveloperProfile(ctx context.Context, devid uint64) (dto.DeveloperReleases, error) {
	return us.userRepository.DeveloperProfile(ctx, devid)
}

func (us *userService) UploadDLC(ctx context.Context, dlc dto.UploadDLC) (entity.DLC, error) {
	return us.userRepository.UploadDLC(ctx, dlc)
}

func (us *userService) PurchaseDLC(ctx context.Context, dlcid uint64, userid uint64, metodeBayar string) (entity.DLC, error) {
	return us.userRepository.PurchaseDLC(ctx, dlcid, userid, metodeBayar)
}

func (us *userService) AddToGame(nama string, gameID uint64, method string) (any, error) {
	switch method {
	case "tags":
		return us.userRepository.AddTags(nama, gameID)
	case "language":
		cek, err := us.userRepository.AddBA(nama, gameID)
		if err != nil {
			return nil, err
		}
		us.userRepository.AddBI(nama, gameID)
		us.userRepository.AddBS(nama, gameID)
		return cek, nil
	case "os":
		return us.userRepository.AddOS(nama, gameID)
	default:
		return nil, errors.New("invalid method")
	}
}

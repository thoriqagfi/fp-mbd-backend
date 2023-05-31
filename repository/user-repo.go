package repository

import (
	"context"
	"errors"
	"mods/dto"
	"mods/entity"
	"time"

	"gorm.io/gorm"
)

type userConnection struct {
	connection *gorm.DB
}

type UserRepository interface {
	// functional
	InsertUser(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, userID uint64) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	UploadGame(ctx context.Context, gameDTO dto.UploadGame, userid uint64) (entity.Game, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) InsertUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := db.connection.Create(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (db *userConnection) GetUserByID(ctx context.Context, userID uint64) (entity.User, error) {
	var user entity.User
	if err := db.connection.Where("id = ?", userID).Take(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (db *userConnection) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	if err := db.connection.Where("email = ?", email).Take(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (db *userConnection) UploadGame(ctx context.Context, gameDTO dto.UploadGame, userid uint64) (entity.Game, error) {
	var developer entity.User
	getDev := db.connection.Where("id = ?", userid).Take(&developer)
	if getDev.Error != nil {
		return entity.Game{}, errors.New("invalid user validation")
	}

	newGame := entity.Game{
		Nama:         gameDTO.Nama,
		Deskripsi:    gameDTO.Deskripsi,
		Release_date: time.Now(),
		Harga:        gameDTO.Harga,
		Age_rating:   gameDTO.Age_rating,
		System_min:   gameDTO.System_min,
		System_rec:   gameDTO.System_rec,
		Picture:      gameDTO.Picture,
		Video:        gameDTO.Video,
		Developer:    developer.Name,
	}

	if err := db.connection.Create(&newGame).Error; err != nil {
		return entity.Game{}, err
	}

	return newGame, nil
}

package repository

import (
	"context"
	"errors"
	"mods/dto"
	"mods/entity"
	"mods/utils"

	"gorm.io/gorm"
)

type storeRepository struct {
	db *gorm.DB
}

type StoreRepository interface {
	// functional
	FeaturedInfo(ctx context.Context) ([]dto.StoreFeatured, error)
	CategoriesInfo(ctx context.Context) ([]dto.StoreCategories, error)
	AllGame(ctx context.Context, pagination utils.Pagination) ([]entity.Game, error)
	GamePage(ctx context.Context, gameid uint64) (entity.Game, error)
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{
		db: db,
	}
}

func (r *storeRepository) FeaturedInfo(ctx context.Context) ([]dto.StoreFeatured, error) {
	var listGame []dto.StoreFeatured

	getList := r.db.Model(&entity.Game{}).Limit(5).Find(&listGame)

	if getList.Error != nil {
		return []dto.StoreFeatured{}, errors.New("failed to get featured information")
	}

	return listGame, nil
}

func (r *storeRepository) CategoriesInfo(ctx context.Context) ([]dto.StoreCategories, error) {
	var listCategories []dto.StoreCategories

	getList := r.db.Model(&entity.Tags{}).Limit(19).Find(&listCategories) // ambil semua data tag dahulu
	if getList.Error != nil {
		return []dto.StoreCategories{}, errors.New("failed to get categories information")
	}

	for i := 0; i < len(listCategories); i++ {
		category := &listCategories[i]

		var gameDetails entity.DetailTagGame
		getGameID := r.db.Model(&entity.DetailTagGame{}).Where("tag_id = ?", category.ID).First(&gameDetails) // ambil game id pada tiap tags
		if getGameID.Error != nil {
			return []dto.StoreCategories{}, errors.New("failed to get game from specific tag")
		}

		var game entity.Game
		getPict := r.db.Model(&entity.Game{}).Where("id = ?", gameDetails.GameID).First(&game) // dari game id yang dimiliki, ambil picturenya
		if getPict.Error != nil {
			return []dto.StoreCategories{}, errors.New("failed to get game pict")
		}

		category.Picture = game.Picture // set pict tag

	}

	return listCategories, nil

}

func (r *storeRepository) AllGame(ctx context.Context, pagination utils.Pagination) ([]entity.Game, error) {
	var game entity.Game
	var listGame []entity.Game

	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := r.db.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Model(&entity.Game{}).Where(game).Find(&listGame)

	if result.Error != nil {
		return nil, errors.New("failed to get all game")
	}

	return listGame, nil

}

func (r *storeRepository) GamePage(ctx context.Context, gameid uint64) (entity.Game, error) {
	var game entity.Game

	getGame := r.db.Where("id = ?", gameid).Preload("ListDLC").Find(&game) // ambil semua data tag dahulu
	if getGame.Error != nil {
		return entity.Game{}, errors.New("failed to get game information")
	}

	return game, nil
}

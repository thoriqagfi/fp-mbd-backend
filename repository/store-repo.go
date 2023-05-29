package repository

import (
	"context"
	"errors"
	"mods/dto"
	"mods/entity"

	"gorm.io/gorm"
)

type storeRepository struct {
	db *gorm.DB
}

type StoreRepository interface {
	// functional
	FeaturedInfo(ctx context.Context) ([]dto.StoreFeatured, error)
	CategoriesInfo(ctx context.Context) ([]dto.StoreCategories, error)
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

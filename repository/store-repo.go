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
	DLCGame(ctx context.Context, dlcid uint64) (entity.DLC, error)
	Popular(ctx context.Context) ([]entity.Game, error)
	FilterTags(nama string) ([]entity.Game, error)
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{
		db: db,
	}
}

func (r *storeRepository) FeaturedInfo(ctx context.Context) ([]dto.StoreFeatured, error) {
	var listGame []dto.StoreFeatured

	getList := r.db.Model(&entity.Game{}).Limit(5).Order("release_date desc").Find(&listGame)

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

	getGame := r.db.Where("id = ?", gameid).Preload("ListDLC").Preload("ListBA").Preload("ListBS").Preload("ListBI").Preload("ListOS").Preload("ListTag").Find(&game) // ambil semua data tag dahulu
	if getGame.Error != nil {
		return entity.Game{}, errors.New("failed to get game information")
	}

	return game, nil
}

func (r *storeRepository) DLCGame(ctx context.Context, dlcid uint64) (entity.DLC, error) {
	var dlc entity.DLC

	getDLC := r.db.Where("id = ?", dlcid).Find(&dlc) // ambil semua data tag dahulu
	if getDLC.Error != nil {
		return entity.DLC{}, errors.New("failed to get dlc information")
	}

	return dlc, nil
}

func (r *storeRepository) Popular(ctx context.Context) ([]entity.Game, error) {
	var games []entity.DetailUserGame

	r.db.Debug().Model(&entity.DetailUserGame{}).Select("game_id, count(user_id)").Group("game_id").Order("count(user_id) desc").Limit(10).Find(&games)

	var getGames []entity.Game
	for _, game := range games {
		var findByID entity.Game
		r.db.Where("id = ?", game.GameID).Take(&findByID)
		getGames = append(getGames, findByID)
	}

	return getGames, nil

}

func (r *storeRepository) FilterTags(nama string) ([]entity.Game, error) {
	var tag entity.Tags
	r.db.Where("nama = ?", nama).Take(&tag)

	var details []entity.DetailTagGame
	getDetail := r.db.Where("tags_id = ?", tag.ID).Find(&details)
	if getDetail.Error != nil {
		return []entity.Game{}, errors.New("no games found")
	}

	var games []entity.Game
	var game entity.Game

	for _, detail := range details {
		r.db.Where("game_id = ?", detail.GameID).Take(&game)
		games = append(games, game)
	}

	return games, nil

}

package service

import (
	"context"
	"mods/dto"
	"mods/entity"
	"mods/repository"
	"mods/utils"
)

type storeService struct {
	storeRepository repository.StoreRepository
}

type StoreService interface {
	GetFeatured(ctx context.Context) ([]dto.StoreFeatured, error)
	GetCategories(ctx context.Context) ([]dto.StoreCategories, error)
	GamePage(ctx context.Context, gameid uint64) (entity.Game, error)
	AllGame(ctx context.Context, pagination utils.Pagination) ([]entity.Game, error)
	DLCGame(ctx context.Context, dlcid uint64) (entity.DLC, error)
	Popular(ctx context.Context) ([]entity.Game, error)
}

func NewStoreService(sr repository.StoreRepository) StoreService {
	return &storeService{
		storeRepository: sr,
	}
}

func (ss *storeService) GetFeatured(ctx context.Context) ([]dto.StoreFeatured, error) {
	getFeatured, err := ss.storeRepository.FeaturedInfo(ctx)
	if err != nil {
		return []dto.StoreFeatured{}, err
	}

	return getFeatured, nil
}

func (ss *storeService) GetCategories(ctx context.Context) ([]dto.StoreCategories, error) {
	getCategories, err := ss.storeRepository.CategoriesInfo(ctx)
	if err != nil {
		return []dto.StoreCategories{}, err
	}

	return getCategories, nil
}

func (ss *storeService) GamePage(ctx context.Context, gameid uint64) (entity.Game, error) {
	return ss.storeRepository.GamePage(ctx, gameid)
}

func (ss *storeService) AllGame(ctx context.Context, pagination utils.Pagination) ([]entity.Game, error) {
	return ss.storeRepository.AllGame(ctx, pagination)
}

func (ss *storeService) DLCGame(ctx context.Context, dlcid uint64) (entity.DLC, error) {
	return ss.storeRepository.DLCGame(ctx, dlcid)
}

func (ss *storeService) Popular(ctx context.Context) ([]entity.Game, error) {
	return ss.storeRepository.Popular(ctx)
}

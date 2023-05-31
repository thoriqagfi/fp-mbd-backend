package service

import (
	"context"
	"mods/dto"
	"mods/repository"
)

type storeService struct {
	storeRepository repository.StoreRepository
}

type StoreService interface {
	GetFeatured(ctx context.Context) ([]dto.StoreFeatured, error)
	GetCategories(ctx context.Context) ([]dto.StoreCategories, error)
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
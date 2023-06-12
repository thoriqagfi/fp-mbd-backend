package dto

type StoreFeatured struct {
	ID        uint64 `json:"game_id" binding:"required"`
	Nama      string `json:"game_title" binding:"required"`
	Deskripsi string `json:"game_description" binding:"required"`
	Harga     string `json:"game_price" binding:"required"`
	Picture   string `json:"game_picture" binding:"required"`
}

type StoreCategories struct {
	ID      uint64 `json:"id" binding:"required"`
	Nama    string `json:"tag_name" binding:"required"`
	Picture string `json:"game_picture" binding:"required"`
}

type StorePopular struct {
	GameID    uint64 `json:"game_id"`
	CountUser uint64 `json:"count_user"`
}

type GetTags struct {
	Nama string `json:"nama" binding:"required"`
}

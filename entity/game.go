package entity

import "time"

type Game struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	Nama         string    `json:"nama" binding:"required"`
	Deskripsi    string    `json:"deskripsi" binding:"required"`
	Release_date time.Time `json:"release_date" binding:"required"`
	Harga        uint64    `json:"harga" binding:"required"`
	Age_rating   string    `json:"age_rating" binding:"required"`
	System_min   string    `json:"system_min" binding:"required"`
	System_rec   string    `json:"system_rec" binding:"required"`
	Picture      string    `json:"picture" binding:"required"`
	Video        string    `json:"video" binding:"required"`

	ListUser      []*User            `gorm:"many2many:detail_user_game;" json:"list_user,omitempty"`
	ListTransaksi []*Transaksi       `gorm:"many2many:detail_transaksi_game;" json:"list_transaksi,omitempty"`
	ListTag       []*Tags            `gorm:"many2many:detail_tag_game;" json:"tags,omitempty"`
	ListBA        []*BahasaAudio     `gorm:"many2many:detail_game_ba;" json:"audio_language,omitempty"`
	ListBI        []*BahasaInterface `gorm:"many2many:detail_game_bi;" json:"interface_language,omitempty"`
	ListBS        []*BahasaSubtitle  `gorm:"many2many:detail_game_bs;" json:"subtitle_language,omitempty"`
	ListOS        []*OperatingSystem `gorm:"many2many:detail_game_os;" json:"operating_systems,omitempty"`

	ListDLC    []DLC    `json:"list_dlc,omitempty"`
	ListReview []Review `json:"list_review,omitempty"`
}

package entity

type DLC struct {
	ID         uint64 `json:"id" gorm:"primaryKey"`
	Nama       string `json:"nama" binding:"required"`
	Deskripsi  string `json:"deskripsi" binding:"required"`
	Harga      uint64 `json:"harga" binding:"required"`
	System_min string `json:"system_min" binding:"required"`
	System_rec string `json:"system_rec" binding:"required"`

	ListUser      []*User      `gorm:"many2many:detail_user_dlc;" json:"list_user,omitempty"`
	ListTransaksi []*Transaksi `gorm:"many2many:detail_transaksi_dlc;" json:"list_transaksi,omitempty"`

	ListReview []Review `json:"list_review,omitempty"`

	GameID uint64 `gorm:"foreignKey" json:"game_id"`
	Game   *Game  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"game,omitempty"`
}

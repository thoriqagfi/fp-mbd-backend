package entity

type Transaksi struct {
	ID           uint64 `json:"id" gorm:"primaryKey"`
	MetodeBayar  string `json:"metode_bayar" binding:"required"`
	TglTransaksi string `json:"tgl_transaksi" binding:"required"`

	ListGame []*Game `gorm:"many2many:detail_transaksi_game;" json:"list_game,omitempty"`
	ListDLC  []*DLC  `gorm:"many2many:detail_transaksi_dlc;" json:"list_dlc,omitempty"`

	UserID uint64 `gorm:"foreignKey" json:"user_id"`
	User   *User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

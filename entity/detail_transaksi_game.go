package entity

type DetailTransaksiGame struct {
	TransaksiID uint64 `gorm:"primaryKey" json:"transaksi_id"`
	GameID      uint64 `gorm:"primaryKey" json:"game_id"`
}

func (DetailTransaksiGame) TableName() string {
	return "detail_transaksi_game"
}

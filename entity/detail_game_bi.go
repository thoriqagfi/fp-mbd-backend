package entity

type DetailGameBI struct {
	GameID uint64 `gorm:"primaryKey" json:"game_id"`
	BiID   uint64 `gorm:"primaryKey" json:"bi_id"`
}

func (DetailGameBI) TableName() string {
	return "detail_game_bi"
}

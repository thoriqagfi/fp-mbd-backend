package entity

type DetailGameBA struct {
	GameID uint64 `gorm:"primaryKey" json:"game_id"`
	BaID   uint64 `gorm:"primaryKey" json:"ba_id"`
}

func (DetailGameBA) TableName() string {
	return "detail_game_ba"
}

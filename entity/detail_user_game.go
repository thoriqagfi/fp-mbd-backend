package entity

type DetailUserGame struct {
	UserID uint64 `gorm:"primaryKey" json:"user_id"`
	GameID uint64 `gorm:"primaryKey" json:"game_id"`
}

func (DetailUserGame) TableName() string {
	return "detail_user_game"
}

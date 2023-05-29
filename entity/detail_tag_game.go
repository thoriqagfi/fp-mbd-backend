package entity

type DetailTagGame struct {
	TagID  uint64 `gorm:"primaryKey" json:"tag_id"`
	GameID uint64 `gorm:"primaryKey" json:"game_id"`
}

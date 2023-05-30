package entity

type DetailGameBS struct {
	GameID uint64 `gorm:"primaryKey" json:"game_id"`
	BsID   uint64 `gorm:"primaryKey" json:"bs_id"`
}

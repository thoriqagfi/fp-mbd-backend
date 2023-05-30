package entity

type DetailGameOS struct {
	GameID uint64 `gorm:"primaryKey" json:"game_id"`
	OsID   uint64 `gorm:"primaryKey" json:"os_id"`
}

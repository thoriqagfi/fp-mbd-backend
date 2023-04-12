package entity

type DetailUserDLC struct {
	UserID uint64 `gorm:"primaryKey" json:"user_id"`
	DLCID  uint64 `gorm:"primaryKey" json:"dlc_id"`
}

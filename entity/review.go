package entity

type Review struct {
	ID      uint64 `json:"id" gorm:"primaryKey"`
	Rating  string `json:"rating" binding:"required"`
	Comment string `json:"comment" binding:"required"`

	GameID uint64 `gorm:"foreignKey" json:"game_id"`
	Game   *Game  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"game,omitempty"`

	DLCID uint64 `gorm:"foreignKey" json:"dlc_id"`
	DLC   *DLC   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"dlc,omitempty"`

	UserID uint64 `gorm:"foreignKey" json:"user_id"`
	User   *User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

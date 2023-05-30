package entity

type BahasaAudio struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Nama string `json:"nama" binding:"required"`

	ListGame []*Game `gorm:"many2many:detail_game_ba;" json:"list_game,omitempty"`
}

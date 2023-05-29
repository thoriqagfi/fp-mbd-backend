package entity

type Tags struct {
	ID   uint64 `json:"id" gorm:"primaryKey"`
	Nama string `json:"nama" binding:"required"`

	ListGame []*Game `gorm:"many2many:detail_tag_game;" json:"list_game,omitempty"`
}

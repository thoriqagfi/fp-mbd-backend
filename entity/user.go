package entity

import (
	"mods/utils"

	"github.com/go-mail/mail"
	"gorm.io/gorm"
)

type User struct {
	ID            uint64 `json:"id" gorm:"primaryKey"`
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password" binding:"required"`
	Profile_image string `json:"profile_image"`
	Role          string `json:"role"`
	Wallet        uint64 `json:"wallet"`

	ListGame      []*Game     `gorm:"many2many:detail_user_game;" json:"list_game,omitempty"`
	ListDLC       []*DLC      `gorm:"many2many:detail_user_dlc;" json:"list_dlc,omitempty"`
	ListTransaksi []Transaksi `json:"list_transaksi,omitempty"`
	ListReview    []Review    `json:"list_review,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.Password, err = utils.PasswordHash(u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	var err error
	if u.Password != "" {
		u.Password, err = utils.PasswordHash(u.Password)
	}
	if err != nil {
		return err
	}
	return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
	m := mail.NewMessage()
	m.SetHeader("From", "steammbd@gmail.com")
	m.SetHeader("To", u.Email)
	m.SetHeader("Subject", "Welcome To Steam MBD!")
	m.SetBody("text/html", "Hello <b>"+u.Name+"</b> Welcome To Our Website and Enjoy!")

	d := mail.NewDialer("smtp.gmail.com", 587, "steammbd@gmail.com", "UPaSS001")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (User) TableName() string {
	return "users"
}

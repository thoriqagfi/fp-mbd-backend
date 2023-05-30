package dto

type UserLoginDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserCreateDTO struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

type UploadGame struct {
	Nama       string `json:"nama" binding:"required"`
	Deskripsi  string `json:"deskripsi" binding:"required"`
	Harga      uint64 `json:"harga" binding:"required"`
	Age_rating string `json:"age_rating" binding:"required"`
	System_min string `json:"system_min" binding:"required"`
	System_rec string `json:"system_rec" binding:"required"`
	Picture    string `json:"picture" binding:"required"`
	Video      string `json:"video" binding:"required"`
}

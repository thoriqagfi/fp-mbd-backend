package entity

type DetailTransaksiDLC struct {
	TransaksiID uint64 `gorm:"primaryKey" json:"transaksi_id"`
	DLCID       uint64 `gorm:"primaryKey" json:"dlc_id"`
}

package repository

import (
	"context"
	"errors"
	"mods/dto"
	"mods/entity"
	"time"

	"gorm.io/gorm"
)

type userConnection struct {
	connection *gorm.DB
}

type UserRepository interface {
	// functional
	InsertUser(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, userID uint64) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)

	// profiles
	UserProfile(ctx context.Context, userid uint64) (entity.User, error)
	DeveloperProfile(ctx context.Context, devid uint64) (dto.DeveloperReleases, error)

	// Transactional
	PurchaseGame(ctx context.Context, gameID uint64, userID uint64, metodeBayar string) (entity.Game, error)
	UploadGame(ctx context.Context, gameDTO dto.UploadGame, userid uint64) (entity.Game, error)
	UploadDLC(ctx context.Context, dlc dto.UploadDLC) (entity.DLC, error)
	PurchaseDLC(ctx context.Context, dlcid uint64, userid uint64, metodeBayar string) (entity.DLC, error)
	TopUp(ctx context.Context, userid uint64, nominal uint64) (entity.User, error)

	// Add Tags Languages OS
	AddTags(tagID uint64, gameID uint64) (entity.Tags, error)
	AddBA(baID uint64, gameID uint64) (entity.BahasaAudio, error)
	AddBI(biID uint64, gameID uint64) (entity.BahasaInterface, error)
	AddBS(bsID uint64, gameID uint64) (entity.BahasaSubtitle, error)
	AddOS(osID uint64, gameID uint64) (entity.OperatingSystem, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) InsertUser(ctx context.Context, user entity.User) (entity.User, error) {
	if err := db.connection.Create(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (db *userConnection) GetUserByID(ctx context.Context, userID uint64) (entity.User, error) {
	var user entity.User
	if err := db.connection.Where("id = ?", userID).Take(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (db *userConnection) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	if err := db.connection.Where("email = ?", email).Take(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (db *userConnection) UserProfile(ctx context.Context, userid uint64) (entity.User, error) {
	var user entity.User
	getDetails := db.connection.Where("id = ?", userid).Preload("ListGame").Preload("ListDLC").Preload("ListTransaksi").Preload("ListReview").Take(&user)
	if getDetails.Error != nil {
		return entity.User{}, errors.New("failed to get user profile")
	}
	return user, nil
}

func (db *userConnection) UploadGame(ctx context.Context, gameDTO dto.UploadGame, userid uint64) (entity.Game, error) {
	var developer entity.User
	getDev := db.connection.Where("id = ?", userid).Take(&developer)
	if getDev.Error != nil {
		return entity.Game{}, errors.New("invalid user validation")
	}

	newGame := entity.Game{
		Nama:         gameDTO.Nama,
		Deskripsi:    gameDTO.Deskripsi,
		Release_date: time.Now(),
		Harga:        gameDTO.Harga,
		Age_rating:   gameDTO.Age_rating,
		System_min:   gameDTO.System_min,
		System_rec:   gameDTO.System_rec,
		Picture:      gameDTO.Picture,
		Video:        gameDTO.Video,
		Developer:    developer.Name,
	}

	if err := db.connection.Create(&newGame).Error; err != nil {
		return entity.Game{}, err
	}

	return newGame, nil
}

func (db *userConnection) PurchaseGame(ctx context.Context, gameID uint64, userID uint64, metodeBayar string) (entity.Game, error) {
	var user entity.User
	getUser := db.connection.Where("id = ?", userID).Take(&user)
	if getUser.Error != nil {
		return entity.Game{}, errors.New("invalid user")
	}

	var game entity.Game
	getGame := db.connection.Where("id = ?", gameID).Take(&game)
	if getGame.Error != nil {
		return entity.Game{}, errors.New("game not found")
	}

	var detail entity.DetailUserGame
	getDetail := db.connection.Debug().Where("user_id = ? AND game_id = ?", userID, gameID).Take(&detail)
	if getDetail.Error == nil {
		return entity.Game{}, errors.New("game already exist in library")
	}

	if metodeBayar == "Steam Wallet" {
		if user.Wallet < game.Harga {
			return entity.Game{}, errors.New("not enough steam wallet")
		}
		db.connection.Model(&user).Where(entity.User{ID: userID}).Update("wallet", (user.Wallet)-game.Harga)
	}

	newTransaksi := entity.Transaksi{
		MetodeBayar:  metodeBayar,
		TglTransaksi: time.Now(),
		UserID:       userID,
	}

	db.connection.Debug().Model(&entity.Transaksi{}).Create(&newTransaksi)

	newDetail := entity.DetailUserGame{
		UserID: userID,
		GameID: game.ID,
	}

	db.connection.Debug().Model(&detail).Create(&newDetail)

	db.connection.Model(&user).Association("ListGame").Append(&game)
	return game, nil
}

func (db *userConnection) TopUp(ctx context.Context, userid uint64, nominal uint64) (entity.User, error) {
	var user entity.User
	getUser := db.connection.Debug().Where("id = ?", userid).Take(&user)
	if getUser.Error != nil {
		return entity.User{}, errors.New("failed to load user")
	}

	db.connection.Model(&user).Where(entity.User{ID: userid}).Update("wallet", (user.Wallet)+nominal)
	return user, nil
}

func (db *userConnection) DeveloperProfile(ctx context.Context, devid uint64) (dto.DeveloperReleases, error) {
	var dev_releases dto.DeveloperReleases
	var developer entity.User
	getDev := db.connection.Where("id = ?", devid).Take(&developer)
	if getDev.Error != nil {
		return dto.DeveloperReleases{}, errors.New("failed to get developer")
	}

	var games []entity.Game
	getGames := db.connection.Where("developer = ?", developer.Name).Find(&games)
	if getGames.Error != nil {
		return dto.DeveloperReleases{}, errors.New("failed to get games")
	}

	for _, game := range games {
		dev_releases.ListGames = append(dev_releases.ListGames, game)
		var listDLC []entity.DLC
		db.connection.Model(&entity.DLC{}).Where("game_id = ?", game.ID).Find(&listDLC)
		dev_releases.ListDLC = append(dev_releases.ListDLC, listDLC...)
	}

	db.connection.Preload("ListGames").Preload("ListDLC").Take(&dev_releases)
	return dev_releases, nil
}

func (db *userConnection) UploadDLC(ctx context.Context, dlc dto.UploadDLC) (entity.DLC, error) {

	newDLC := entity.DLC{
		Nama:       dlc.Nama,
		Deskripsi:  dlc.Deskripsi,
		Harga:      dlc.Harga,
		System_min: dlc.System_min,
		System_rec: dlc.System_rec,
		Picture:    dlc.Picture,
		GameID:     dlc.GameID,
	}

	if err := db.connection.Create(&newDLC).Error; err != nil {
		return entity.DLC{}, errors.New("failed to upload DLC")
	}

	return newDLC, nil
}

func (db *userConnection) PurchaseDLC(ctx context.Context, dlcid uint64, userid uint64, metodeBayar string) (entity.DLC, error) {
	var user entity.User
	getUser := db.connection.Where("id = ?", userid).Take(&user)
	if getUser.Error != nil {
		return entity.DLC{}, errors.New("invalid user")
	}

	var dlc entity.DLC
	getGame := db.connection.Where("id = ?", dlcid).Take(&dlc)
	if getGame.Error != nil {
		return entity.DLC{}, errors.New("dlc not found")
	}

	var detail entity.DetailUserDLC
	getDetail := db.connection.Debug().Where("user_id = ? AND dlc_id = ?", userid, dlcid).Take(&detail)
	if getDetail.Error == nil {
		return entity.DLC{}, errors.New("dlc already exist in library")
	}

	if metodeBayar == "Steam Wallet" {
		if user.Wallet < dlc.Harga {
			return entity.DLC{}, errors.New("not enough steam wallet")
		}
		db.connection.Model(&user).Where(entity.User{ID: userid}).Update("wallet", (user.Wallet)-dlc.Harga)
	}

	newTransaksi := entity.Transaksi{
		MetodeBayar:  metodeBayar,
		TglTransaksi: time.Now(),
		UserID:       userid,
	}

	db.connection.Debug().Model(&entity.Transaksi{}).Create(&newTransaksi)

	newDetail := entity.DetailUserDLC{
		UserID: userid,
		DLCID:  dlc.ID,
	}

	db.connection.Debug().Model(&detail).Create(&newDetail)

	db.connection.Model(&user).Association("ListDLC").Append(&dlc)
	return dlc, nil
}

func (db *userConnection) AddTags(tagID uint64, gameID uint64) (entity.Tags, error) {
	var game entity.Game
	var tag entity.Tags
	var detail entity.DetailTagGame

	if err := db.connection.Where("tag_id = ? AND game_id = ?", tagID, gameID).Take(&detail).Error; err == nil {
		return entity.Tags{}, errors.New("selected tag already exist")
	}

	db.connection.Where("id = ?", gameID).Take(&game)
	db.connection.Where("id = ?", tagID).Take(&tag)

	newDetail := entity.DetailTagGame{
		TagID:  tag.ID,
		GameID: game.ID,
	}
	db.connection.Debug().Model(&detail).Create(&newDetail)

	db.connection.Model(&game).Association("ListTag").Append(&tag)
	return tag, nil

}

func (db *userConnection) AddBA(baID uint64, gameID uint64) (entity.BahasaAudio, error) {
	var detail entity.DetailGameBA
	var game entity.Game
	var ba entity.BahasaAudio

	if err := db.connection.Where("game_id = ? AND ba_id = ?", gameID, baID).Take(&detail).Error; err == nil {
		return entity.BahasaAudio{}, errors.New("selected bahasa already exist")
	}
	db.connection.Where("id = ?", gameID).Take(&game)
	db.connection.Where("id = ?", baID).Take(&ba)

	newDetail := entity.DetailGameBA{
		GameID: game.ID,
		BaID:   ba.ID,
	}
	db.connection.Debug().Model(&detail).Create(&newDetail)

	db.connection.Model(&game).Association("ListBA").Append(&ba)
	return ba, nil
}

func (db *userConnection) AddBI(biID uint64, gameID uint64) (entity.BahasaInterface, error) {
	var detail entity.DetailGameBI
	var game entity.Game
	var bi entity.BahasaInterface

	if err := db.connection.Where("game_id = ? AND bi_id = ?", gameID, biID).Take(&detail).Error; err == nil {
		return entity.BahasaInterface{}, errors.New("selected bahasa already exist")
	}
	db.connection.Where("id = ?", gameID).Take(&game)
	db.connection.Where("id = ?", biID).Take(&bi)

	newDetail := entity.DetailGameBI{
		GameID: game.ID,
		BiID:   bi.ID,
	}
	db.connection.Debug().Model(&detail).Create(&newDetail)

	db.connection.Model(&game).Association("ListBI").Append(&bi)
	return bi, nil
}

func (db *userConnection) AddBS(bsID uint64, gameID uint64) (entity.BahasaSubtitle, error) {
	var detail entity.DetailGameBS
	var game entity.Game
	var bs entity.BahasaSubtitle

	if err := db.connection.Where("game_id = ? AND bs_id = ?", gameID, bsID).Take(&detail).Error; err == nil {
		return entity.BahasaSubtitle{}, errors.New("selected bahasa already exist")
	}
	db.connection.Where("id = ?", gameID).Take(&game)
	db.connection.Where("id = ?", bsID).Take(&bs)

	newDetail := entity.DetailGameBS{
		GameID: game.ID,
		BsID:   bs.ID,
	}
	db.connection.Debug().Model(&detail).Create(&newDetail)

	db.connection.Model(&game).Association("ListBS").Append(&bs)
	return bs, nil
}

func (db *userConnection) AddOS(osID uint64, gameID uint64) (entity.OperatingSystem, error) {
	var detail entity.DetailGameOS
	var game entity.Game
	var os entity.OperatingSystem

	if err := db.connection.Where("game_id = ? AND os_id = ?", gameID, osID).Take(&detail).Error; err == nil {
		return entity.OperatingSystem{}, errors.New("selected os already exist")
	}
	db.connection.Where("id = ?", gameID).Take(&game)
	db.connection.Where("id = ?", osID).Take(&os)

	newDetail := entity.DetailGameOS{
		GameID: game.ID,
		OsID:   os.ID,
	}
	db.connection.Debug().Model(&detail).Create(&newDetail)

	db.connection.Model(&game).Association("ListOS").Append(&os)
	return os, nil
}

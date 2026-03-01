package sql

type DiscordGambleGamer struct {
	Id    string `gorm:"primaryKey;size:36"          json:"id"`    // 36字元 UUID 主鍵
	Name  string `gorm:"size:100;not null"           json:"name"`  // 最多100字，不可為 NULL
	Chips int64  `gorm:"not null;default:0"          json:"chips"` // 起始籌碼 0，必填
}

type User struct {
	ID    string `gorm:"primaryKey;size:36"            json:"id"`
	Name  string `gorm:"size:100;not null"             json:"name"`
	Email string `gorm:"size:255;uniqueIndex;not null" json:"email"`
}

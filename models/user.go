package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"size:255;not null;unique" json:"email"`
	Password     string `gorm:"size:255;" json:"password"`
	IsGithubUser bool   `gorm:"not null" json:"is_github_user"`
	CanSendEmail bool   `gorm:"not null" json:"can_send_email"`
}

// User オブジェクトをデータベースに保存する
func (u *User) Save(db *gorm.DB) (*User, error) {
	err := db.Create(u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

// パスワードを空文字にして出力準備をする
func (u *User) PrepareOutput() *User {
	u.Password = ""
	return u
}

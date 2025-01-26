package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"size:255;not null;unique" json:"email"`
	Password     string `gorm:"size:255;" json:"password"`
	IsGithubUser bool   `gorm:"not null" json:"is_github_user"`
}

// User オブジェクトをデータベースに保存する
func (u *User) Save(db *gorm.DB) (*User, error) {
	err := db.Create(u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Userオブジェクトが保存される前に実行する
func (u *User) BeforeSave(*gorm.DB) error {
	// パスワードをハッシュ化する
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

// パスワードを空文字にして出力準備をする
func (u *User) PrepareOutput() *User {
	u.Password = ""
	return u
}

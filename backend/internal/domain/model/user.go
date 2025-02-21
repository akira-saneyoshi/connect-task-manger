package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string //ハッシュ化されたパスワード
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser は新しい User エンティティを作成します。
func NewUser(id, name, email, rawPassword string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: string(hashedPassword), // ハッシュ化されたパスワードを保存
	}, nil
}

// Authenticate は提供されたパスワードがユーザーのハッシュ化されたパスワードと一致するかを検証します。
func (u *User) Authenticate(rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(rawPassword))
}

// Update はユーザーの情報を更新します。
func (u *User) Update(name, email, rawPassword string) error {
	u.Name = name
	u.Email = email

	if rawPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}

	return nil
}

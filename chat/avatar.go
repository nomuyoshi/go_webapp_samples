package main

import (
	"errors"
	"fmt"
	"path/filepath"

	gomniauthcommon "github.com/stretchr/gomniauth/common"
)

// ErrNoAvatarURL はインスタンスがアバターURLを返すことが出来ないときに発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターURLを取得できません")

// ChatUser はチャットユーザーのインターフェイス
type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

// UniqueID はチャットユーザーのユニークIDを返す
func (u chatUser) UniqueID() string {
	return u.uniqueID
}

type Avatar interface {
	// GetAvatarURL はクライアントのアバターURLを取得
	GetAvatarURL(u *chatUser) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}

	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if u.UniqueID() == "" {
		return "", ErrNoAvatarURL
	}

	return fmt.Sprintf("//www.gravatar.com/avatar/%s", u.UniqueID()), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if u.UniqueID() == "" {
		return "", ErrNoAvatarURL
	}

	matches, err := filepath.Glob(filepath.Join("avatars", u.UniqueID()+"*"))
	if err != nil || len(matches) == 0 {
		return "", ErrNoAvatarURL
	}

	return matches[0], nil
}

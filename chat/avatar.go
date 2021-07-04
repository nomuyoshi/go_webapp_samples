package main

import (
	"errors"
)

// ErrNoAvatarURL はインスタンスがアバターURLを返すことが出来ないときに発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターURLを取得できません")

type Avatar interface {
	// GetAvatarURL はクライアントのアバターURLを取得
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}

	return "", ErrNoAvatarURL
}

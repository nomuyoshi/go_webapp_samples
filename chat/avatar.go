package main

import (
	"errors"
	"fmt"
)

// ErrNoAvatarURL はインスタンスがアバターURLを返すことが出来ないときに発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターURLを取得できません")

type Avatar interface {
	// GetAvatarURL はクライアントのアバターURLを取得
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}

	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userID, ok := c.userData["user_id"]; ok {
		if userIDStr, ok := userID.(string); ok {
			return fmt.Sprintf("//www.gravatar.com/avatar/%s", userIDStr), nil
		}
	}

	return "", ErrNoAvatarURL
}

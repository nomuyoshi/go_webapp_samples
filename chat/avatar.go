package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
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
	if url, ok := c.userData["avatarURL"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}

	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userID, ok := c.userData["userID"]; ok {
		if userIDStr, ok := userID.(string); ok {
			return fmt.Sprintf("//www.gravatar.com/avatar/%s", userIDStr), nil
		}
	}

	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userID, ok := c.userData["userID"]; ok {
		if userIDStr, ok := userID.(string); ok {
			if files, err := ioutil.ReadDir("avatars"); err == nil {
				for _, file := range files {
					if file.IsDir() {
						continue
					}

					if strings.HasPrefix(file.Name(), userIDStr) {
						return "/avatars/" + file.Name(), nil
					}
				}
			}
		}
	}

	return "", ErrNoAvatarURL
}

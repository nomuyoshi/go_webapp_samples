package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、AuthAvatar.GetAvatarURLはErrNoAvatarURLを返すべきです")
	}

	testURL := "http://url-to-avatar/"
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return(testURL, nil)
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != nil && err != ErrNoAvatarURL {
		t.Error("値が存在する場合、AuthAvatar.GetAvatarURLはエラーを返すべきではありません。unexpected error: ", err)
	} else {
		if url != testURL {
			t.Error("AuthAvatar.GetAvatarURLは正しいURLを返すべきです")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	chatUser := &chatUser{uniqueID: "0bc83cb571cd1c50ba6f3e8a78ef1346"}
	url, err := gravatarAvatar.GetAvatarURL(chatUser)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではありません。error: ", err)
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvitar.GetAvatarURLが%sという誤った値を返しました", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	var fileSystemAvatar FileSystemAvatar
	filename := filepath.Join("avatars", "0bc83cb571cd1c50ba6f3e8a78ef1346.jpg")
	ioutil.WriteFile(filename, []byte{}, 0644)
	defer func() {
		os.Remove(filename)
	}()

	chatUser := &chatUser{uniqueID: "0bc83cb571cd1c50ba6f3e8a78ef1346"}
	url, err := fileSystemAvatar.GetAvatarURL(chatUser)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURLはエラーを返すべきではありません。error: ", err)
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURLが%sという誤った値を返しました", url)
	}
}

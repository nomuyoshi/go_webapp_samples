package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	_, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、AuthAvatar.GetAvatarURLはErrNoAvatarURLを返すべきです")
	}

	testURL := "http://url-to-avatar/"
	client.userData = map[string]interface{}{"avatarURL": testURL}
	url, err := authAvatar.GetAvatarURL(client)
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
	client := new(client)
	client.userData = map[string]interface{}{"userID": "0bc83cb571cd1c50ba6f3e8a78ef1346"}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURLはエラーを返すべきではありません。error: ", err)
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvitar.GetAvatarURLが%sという誤った値を返しました", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	var fileSystemAvatar FileSystemAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userID": "abc"}
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0644)
	defer func() {
		os.Remove(filename)
	}()

	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURLはエラーを返すべきではありません。error: ", err)
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURLが%sという誤った値を返しました", url)
	}
}

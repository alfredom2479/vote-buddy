package main

import (
	"errors"
	"strings"
)

func getCommentFullName(shareLink string) (string, error) {

	fullNameIndex := strings.Index(shareLink, "comment/")

	if fullNameIndex == -1 {
		return "", errors.New(" 'comment' not found in share link path")
	}

	fullNameIndex += 8

	if len(shareLink) < fullNameIndex+7 {
		return "", errors.New("share link was cut off too soon")
	}

	return "t1_" + shareLink[fullNameIndex:fullNameIndex+7], nil

}

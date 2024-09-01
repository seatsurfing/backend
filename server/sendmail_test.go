package main

import (
	"path/filepath"
	"testing"
)

func TestGetEmailTemplatePathExists(t *testing.T) {
	res, err := getEmailTemplatePath(EmailTemplateSignup, "de")
	checkStringNotEmpty(t, res)
	checkTestBool(t, true, err == nil)
}

func TestGetEmailTemplatePathFallback(t *testing.T) {
	res, err := getEmailTemplatePath(EmailTemplateSignup, "notexists")
	checkStringNotEmpty(t, res)
	checkTestBool(t, true, err == nil)
}

func TestGetEmailTemplatePathNotExists(t *testing.T) {
	path, _ := filepath.Abs("./res/notexisting.txt")
	res, err := getEmailTemplatePath(path, "en")
	checkTestString(t, "", res)
	checkTestBool(t, true, err != nil)
}

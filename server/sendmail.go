package main

import (
	"io/ioutil"
	"net/smtp"
	"path/filepath"
	"strings"
)

var EmailTemplateSignup, _ = filepath.Abs("./res/email-signup.txt")
var EmailTemplateConfirm, _ = filepath.Abs("./res/email-confirm.txt")
var EmailTemplateResetpassword, _ = filepath.Abs("./res/email-resetpw.txt")
var SendMailMockContent = ""

func sendEmail(recipient, sender, templateFile, language string, vars map[string]string) error {
	templateFile = strings.ReplaceAll(templateFile, ".txt", "_"+language+".txt")
	body, err := compileEmailTemplate(templateFile, vars)
	if err != nil {
		return err
	}
	if GetConfig().MockSendmail {
		SendMailMockContent = body
		return nil
	}
	to := []string{recipient}
	msg := []byte(body)
	err = smtp.SendMail(GetConfig().SMTPHost, nil, sender, to, msg)
	return err
}

func compileEmailTemplate(templateFile string, vars map[string]string) (string, error) {
	data, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return "", err
	}
	s := string(data)
	for key, val := range vars {
		s = strings.ReplaceAll(s, "{{"+key+"}}", val)
	}
	return s, nil
}

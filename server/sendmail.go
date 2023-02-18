package main

import (
	"crypto/tls"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
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
	err = smtpDialAndSend(sender, to, msg)
	return err
}

func compileEmailTemplate(templateFile string, vars map[string]string) (string, error) {
	data, err := os.ReadFile(templateFile)
	if err != nil {
		return "", err
	}
	s := string(data)
	c := GetConfig()
	vars["frontendUrl"] = c.FrontendURL
	vars["publicUrl"] = c.PublicURL
	vars["senderAddress"] = c.SMTPSenderAddress
	for key, val := range vars {
		s = strings.ReplaceAll(s, "{{"+key+"}}", val)
	}
	return s, nil
}

func smtpDialAndSend(from string, to []string, msg []byte) error {
	config := GetConfig()
	addr := config.SMTPHost + ":" + strconv.Itoa(config.SMTPPort)
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if config.SMTPStartTLS {
		if ok, _ := c.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{
				ServerName:         config.SMTPHost,
				InsecureSkipVerify: config.SMTPInsecureSkipVerify,
			}
			if err = c.StartTLS(tlsConfig); err != nil {
				return err
			}
		}
	}
	if config.SMTPAuth {
		auth := smtp.PlainAuth("", config.SMTPAuthUser, config.SMTPAuthPass, config.SMTPHost)
		if err = c.Auth(auth); err != nil {
			return err
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

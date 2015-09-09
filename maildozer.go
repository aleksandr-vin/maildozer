package main

import (
	"bytes"
	"fmt"
	"net/smtp"
	"runtime"
	"strings"
	"text/template"
)

func main() {
	fmt.Println("Maildozer is here to send your nasty mails, BOSS...")

	SendEmail(
		"ginger.vinpn",
		25,
		"alex.vinokurov@evecon.co",
		"",
		[]string{"john.the.tester@evecon.co"},
		"Testing subject",
		"<html><body><h1>Test message</h1></body></html>")

	fmt.Println("Accomplished")
}

func printStack() {
	// Capture the stack trace
	buf := make([]byte, 10000)
	bytesWritten := runtime.Stack(buf, false)
	lBuf := bytes.NewBuffer(buf)
	lBuf.Truncate(bytesWritten)

	fmt.Printf("Stack Trace: %s", lBuf.String())
}

func catchPanic(err *error, functionName string) {
	if r := recover(); r != nil {
		fmt.Printf("%s: PANIC Defered: %v\n", functionName, r)
		printStack()
		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	} else if err != nil && *err != nil {
		fmt.Printf("%s: ERROR: %v\n", functionName, *err)
		printStack()
	}
}

func SendEmail(host string, port int, userName string, password string, to []string, subject string, message string) (err error) {
	defer catchPanic(&err, "SendEmail")

	parameters := struct {
		From    string
		To      string
		Subject string
		Message string
	}{
		userName,
		strings.Join([]string(to), ","),
		subject,
		message,
	}

	buffer := new(bytes.Buffer)

	template := template.Must(template.New("emailTemplate").Parse(emailScript()))
	template.Execute(buffer, &parameters)

	//auth := smtp.PlainAuth("", userName, password, host)
	var auth smtp.Auth = nil

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		userName,
		to,
		buffer.Bytes())

	return err
}

func emailScript() (script string) {
	return `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

{{.Message}}`
}

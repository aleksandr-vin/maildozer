package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	htmltpl "html/template"
	"log"
	"net/smtp"
	"os"
	"runtime"
	"strings"
	"text/template"
)

var doSend bool = false // Send emails or not?

func readConfig(filename string) (data []byte, err error) {
	fl, err := os.Open(filename)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	defer fl.Close()
	buf := make([]byte, 10000000)
	bytesRead, err := fl.Read(buf)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	lBuf := bytes.NewBuffer(buf)
	lBuf.Truncate(bytesRead)
	data = lBuf.Bytes()
	return
}

func main() {
	fmt.Println("Maildozer is here to send your nasty mails, BOSS...")

	if len(os.Args) != 2 {
		fmt.Printf("Specify config file name when calling the app!")
		return
	}

	configData, err := readConfig(os.Args[1])
	t := make(map[interface{}]interface{})

	err = yaml.Unmarshal(configData, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var subjTplStr = t["subject-template"].(string)
	subjTpl := template.New("subject")
	subjTpl, err = subjTpl.Parse(subjTplStr)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var bodyTplFilename = t["body-template"].(string)

	doSend = t["do-send"].(bool)

	if t["debug"].(bool) {
		log.Printf("Config data: %v\n", t)
		log.Printf("Mail server: %s\n", t["mail-server"])
		log.Printf("From: %s\n", t["from"])
		log.Printf("Subject template: %s\n", subjTplStr)
		log.Printf("Body template file name: %s\n", bodyTplFilename)
		log.Printf("Do send emails: %v\n", doSend)
	}

	bodyTpl := htmltpl.New("body")
	bodyTpl, err = htmltpl.ParseFiles(bodyTplFilename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, val := range t["to"].([]interface{}) {
		var recipient map[interface{}]interface{} = val.(map[interface{}]interface{})

		var email string = recipient["email"].(string)

		subject, err := makeSubj(subjTpl, recipient)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		body, err := makeBody(bodyTpl, recipient)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		if t["debug"].(bool) {
			log.Printf("Body: %v\n", body)
		}

		SendEmail(
			t["mail-server"].(string),
			t["from"].(string),
			"",
			[]string{email},
			subject,
			body)
		log.Printf("Message '%s' sent to %v <%v>\n",
			subject,
			recipient["fullname"], email)
	}

	fmt.Println("Accomplished")
}

func makeSubj(tpl *template.Template, params map[interface{}]interface{}) (subject string, err error) {
	buffer := new(bytes.Buffer)
	err = tpl.Execute(buffer, params)
	subject = buffer.String()
	return
}

func makeBody(tpl *htmltpl.Template, params map[interface{}]interface{}) (body string, err error) {
	buffer := new(bytes.Buffer)
	err = tpl.Execute(buffer, params)
	body = buffer.String()
	return
}

func printStack() {
	// Capture the stack trace
	buf := make([]byte, 10000)
	bytesWritten := runtime.Stack(buf, false)
	lBuf := bytes.NewBuffer(buf)
	lBuf.Truncate(bytesWritten)

	log.Fatalf("Stack Trace: %s", lBuf.String())
}

func catchPanic(err *error, functionName string) {
	if r := recover(); r != nil {
		log.Fatalf("%s: PANIC Defered: %v\n", functionName, r)
		printStack()
		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	} else if err != nil && *err != nil {
		log.Fatalf("%s: ERROR: %v\n", functionName, *err)
		printStack()
	}
}

func SendEmail(mailServer string, userName string, password string, to []string, subject string, message string) (err error) {
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

	if doSend {
		log.Printf("Sending mail to %v\n", userName)
		err = smtp.SendMail(
			mailServer,
			auth,
			userName,
			to,
			buffer.Bytes())
	}

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

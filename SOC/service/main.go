package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
)

var flag = os.Getenv("FLAG")
var dbLogin = os.Getenv("LOGIN")
var dbPassword = os.Getenv("PASS")
var ctx = logrus.Entry{}

func renderIndex(w http.ResponseWriter, r *http.Request) {
	var form = `
	<div><form action="/login" method="POST">
	<p>Username: <input type="text" name="username" /></p>
	<p>Password: <input type="password" name="password" /></p>
	<p><input type="submit" value="Login" /></p>
	</form></div>`
	fmt.Fprint(w, form)
}

func Login(w http.ResponseWriter, r *http.Request) {
	log := logrus.New()
	conn, err := net.Dial("tcp", "logstash:5000")
	if err != nil {
		log.Fatal(err)
	}

	login := r.FormValue("username")
	password := r.FormValue("password")
	hashPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{
		"type":         "loginService",
		"rAddr":        r.Header.Get("X-FORWARDED-FOR"),
		"login":        login,
		"hashPassword": hashPassword,
	}))
	log.Hooks.Add(hook)
	ctx := log.WithFields(logrus.Fields{
		"method": "main",
	})
	if login == dbLogin && password == dbPassword {
		ctx.Info("Success authentication")
		fmt.Fprintf(w, "Flag{%s}\n", flag)
	} else {
		ctx.Warn("Wrong credentials")
		fmt.Fprint(w, "Auth Error")
	}
}

func main() {
	log := logrus.New()
	conn, err := net.Dial("tcp", "logstash:5000")
	if err != nil {
		log.Fatal(err)
	}
	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "loginService"}))
	log.Hooks.Add(hook)
	ctx := log.WithFields(logrus.Fields{
		"method": "main",
	})
	ctx.Info("Web Server Started")
	http.HandleFunc("/login", Login)
	http.HandleFunc("/", renderIndex)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

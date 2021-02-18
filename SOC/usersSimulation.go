package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

const adminPass = "admin"
const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var numOfRequests = 0

func generatePassword() (randPass string) {
	b := make([]byte, 5)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)

}

func generateIp() (ip string) {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

type UserPattern struct {
	password string
	route    string
	ip       string
}

func (u *UserPattern) sendRequest() {
	endpoint := "http://localhost:8080" + u.route
	data := url.Values{}
	data.Set("username", "admin")
	data.Set("password", u.password)

	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Add("X-Forwarded-For", u.ip)
	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	numOfRequests++
	if numOfRequests%100 == 0 {
		fmt.Println(numOfRequests, "requests sent")
	}
	//_, err = ioutil.ReadAll(res.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

type Simulation struct {
	hackersPercent       int
	hackersAttempts      int
	usersErrorAttempts   int
	numOfUsers           int
	usersIp              []string
	usersDelay           int
	usersAttemptsDelay   int
	hackersAttemptsDelay int
	hackersIP            string
	influxOfUsersPercent int
	intensityOfUsers     int
}

func (s *Simulation) GenerateUsersIp() {
	for i := 0; i < s.numOfUsers; i++ {
		s.usersIp = append(s.usersIp, generateIp())
	}
}

func (s *Simulation) Go() {
	for {

		hackersCondition := rand.Intn(100)
		influxOfUsers := rand.Intn(100)
		if hackersCondition <= s.hackersPercent {
			go func() {
				hacker := UserPattern{route: "/login"}
				hacker.ip = s.hackersIP

				for i := 0; i < s.hackersAttempts; i++ {
					hacker.password = generatePassword()
					hacker.sendRequest()
					time.Sleep(time.Duration(s.hackersAttemptsDelay) * time.Millisecond)
				}
			}()
		}
		if influxOfUsers <= s.influxOfUsersPercent {
			for i := 0; i <= s.intensityOfUsers; i++ {
				go func() {
					user := UserPattern{}
					user.route = "/login"
					user.ip = s.usersIp[rand.Intn(len(s.usersIp))]

					for i := 1; i <= rand.Intn(s.usersErrorAttempts); i++ {
						user.password = generatePassword()
						user.sendRequest()
						time.Sleep(time.Duration(s.usersAttemptsDelay) * time.Millisecond)
					}
					user.password = adminPass
					user.sendRequest()
					time.Sleep(time.Duration(s.hackersAttemptsDelay) * time.Millisecond)
				}()
			}
		} else {
			go func() {
				user := UserPattern{}
				user.route = "/login"
				user.ip = s.usersIp[rand.Intn(len(s.usersIp))]

				for i := 1; i <= rand.Intn(s.usersErrorAttempts); i++ {
					user.password = generatePassword()
					user.sendRequest()
					time.Sleep(time.Duration(s.usersAttemptsDelay) * time.Millisecond)
				}
				user.password = adminPass
				user.sendRequest()
				time.Sleep(time.Duration(s.usersDelay) * time.Millisecond)
			}()

		}
		time.Sleep(time.Duration(s.usersDelay) * time.Millisecond)
	}
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c

		fmt.Println("\nSimulation stopped...")
		os.Exit(0)
	}()
	rand.Seed(time.Now().UnixNano())
	sim := Simulation{
		hackersPercent:       10,
		hackersAttempts:      10,
		hackersAttemptsDelay: 300,
		hackersIP:            generateIp(),

		usersErrorAttempts:   2,
		numOfUsers:           10,
		usersAttemptsDelay:   2000,
		usersDelay:           1000,
		influxOfUsersPercent: 20,
		intensityOfUsers:     10,
	}
	sim.GenerateUsersIp()
	fmt.Println("Simulation started...")
	sim.Go()
}

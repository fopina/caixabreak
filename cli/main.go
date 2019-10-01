package main

import (
	"flag"
	"fmt"
	"handler/function"
	"log"
)

func main() {
	loginPtr := flag.Bool("login", false, "Login")
	tokenPtr := flag.String("token", "", "Token")
	logoutPtr := flag.Bool("logout", false, "Terminate session")

	flag.Parse()

	var token string
	var err error

	if *loginPtr {
		if len(flag.Args()) < 2 {
			log.Fatal("Missing username and password args")
		}

		token, err = function.Login(flag.Arg(0), flag.Arg(1))

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Token: " + token)
	}

	if *tokenPtr != "" {
		token = *tokenPtr
	}

	if *logoutPtr {
		err = function.Logout(token)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("logged out")
		return
	}

	days, err := function.GetData(token)

	if err != nil {
		log.Fatalln(err)
	}

	for _, day := range days.History {
		fmt.Println(day)
	}

	days, err = function.GetDataForMonth(token, days.ViewState, days.CardNumber, days.PreviousExtracts[0])

	if err != nil {
		log.Fatalln(err)
	}

	for _, day := range days.History {
		fmt.Println(day)
	}
}

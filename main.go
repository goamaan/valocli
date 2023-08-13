package main

import (
	"fmt"

	"github.com/goamaan/valocli/internal/core"
)

func main() {
	client := core.New(nil)

	data, err := client.Authorize("amaan18", "amaangokak18")
	var multifactorCode string

	if err != nil {
		if err == core.ErrorRiotMultifactor {
			fmt.Println("Seems like you have Multi factor set up. Input the code sent to your email: ")
			_, err := fmt.Scanln(&multifactorCode)
			if err != nil {
				return
			}

			data, err = client.SubmitTwoFactor(multifactorCode)
		} else {
			panic(err)
		}
	}

	fmt.Println(data.AccessToken)
}

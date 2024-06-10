package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"interparkTester/internal/interpark"
	"interparkTester/pkg/custom"
)

const (
	appVersion = "0.0.1"
	appAuthor  = "Zhin"
	appYear    = "2024"
)

func main() {
	for {
		custom.PrintLogo(appVersion, appAuthor, appYear)

		reader := bufio.NewReader(os.Stdin)

		fmt.Print(string(custom.Green), "Enter your user ID: ", string(custom.Reset))
		userID, _ := reader.ReadString('\n')
		userID = strings.TrimSpace(userID)

		fmt.Print(string(custom.Red), "Enter your password: ", string(custom.Reset))
		userPwd, _ := reader.ReadString('\n')
		userPwd = strings.TrimSpace(userPwd)

		err := interpark.LoginProcess(userID, userPwd)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goamaan/valocli/internal/core"
)

type AuthConfiguration struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	ConfigFileDirectory = ".valocli"
	ConfigFilePath      = "valocli_config.json"
)

func main() {
	client := core.New(nil)

	config := readFromConfig()

	data, err := client.Authorize(config.Username, config.Password)

	if err != nil {
		if err == core.ErrorRiotMultifactor {
			fmt.Println("Seems like you have Multi factor set up. Enter the code sent to your email: ")
			var multifactorCode string
			fmt.Scan(&multifactorCode)

			data, err = client.SubmitTwoFactor(multifactorCode)
		} else {
			panic(err)
		}
	}

	fmt.Println(data)
}

func readFromConfig() AuthConfiguration {
	var config AuthConfiguration
	configPath := getConfigPath()

	fmt.Println("Welcome to the Valorant CLI:")
	fmt.Println()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Configuration file doesn't exist, prompt for username and password
		userAuthInput(&config)

		saveConfiguration(configPath, config)
	} else {
		// Configuration file exists, read file
		config = loadConfiguration(configPath)
		fmt.Printf("Use previously saved username (%s) and password?: Y/n - ", config.Username)
		var usePrevious string
		fmt.Scan(&usePrevious)
		if usePrevious != "Y" && usePrevious != "y" {
			userAuthInput(&config)
			saveConfiguration(configPath, config)
		}
	}

	return config
}

func userAuthInput(config *AuthConfiguration) {
	fmt.Println("Please enter your username:")
	fmt.Scan(&config.Username)
	fmt.Println("Please enter your password:")
	fmt.Scan(&config.Password)
}

func saveConfiguration(path string, config AuthConfiguration) {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling configuration:", err)
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(getConfigDirectory(), 0700)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Println("Error writing configuration file:", err)
	}
}

func loadConfiguration(path string) AuthConfiguration {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading configuration file:", err)
		return AuthConfiguration{}
	}

	var config AuthConfiguration
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshaling configuration:", err)
		return AuthConfiguration{}
	}

	return config
}

func getConfigPath() string {
	configDir := getConfigDirectory()

	configPath := filepath.Join(configDir, ConfigFilePath)

	return configPath
}

func getConfigDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return ""
	}

	configDir := filepath.Join(homeDir, ConfigFileDirectory)

	return configDir
}

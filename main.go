package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/goamaan/valocli/internal/core"
)

type AuthConfiguration struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	ConfigFileDirectory  = ".valocli"
	ConfigFilePath       = "valocli_config.json"
	AuthSaveDataFilePath = "valocli_auth_save.json"
)

func main() {
	client := core.New(nil)

	config, saveData := readFromConfig()
	if saveData != nil {
		client.AuthData = saveData
		// ensure that saved auth data actually works
		err := client.SetUserId()
		if err != nil {
			fmt.Printf("Got error: %s. Previous tokens have expired. Logging in again...\n", err)
		} else {
			saveAuthSaveData(getSaveDataPath(), client.AuthData)
			cliLoop(client)
			return
		}
	}

	err := client.Authorize(config.Username, config.Password)
	if err != nil {
		if err == core.ErrorRiotMultifactor {
			fmt.Println("Seems like you have Multi factor set up. Enter the code sent to your email: ")
			var multifactorCode string
			fmt.Scan(&multifactorCode)

			err = client.MultiFactorAuth(multifactorCode)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	saveAuthSaveData(getSaveDataPath(), client.AuthData)
	cliLoop(client)
}

func cliLoop(c *core.Client) {
	fmt.Println("authData: ", c.AuthData)
	var response string
	for {
		fmt.Scan(&response)
		fmt.Println(response)
	}
}

func readFromConfig() (AuthConfiguration, *core.AuthSaveData) {
	var config AuthConfiguration
	configPath := getConfigPath()

	fmt.Println("Welcome to the Valorant CLI:")
	fmt.Println()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Configuration file doesn't exist, prompt for username and password
		userAuthInput(&config)

		saveConfiguration(configPath, config)
		return config, nil
	}
	// Configuration file exists, read file
	config = loadConfiguration(configPath)
	fmt.Printf("Use previously saved username (%s) and password?: Y/n - ", config.Username)
	var usePrevious string
	fmt.Scan(&usePrevious)

	if usePrevious == "n" || usePrevious == "N" {
		userAuthInput(&config)
		saveConfiguration(configPath, config)
		return config, nil
	}

	// using previous username, password, so try saved auth data
	saveData := readFromSaveData()
	return config, saveData
}

func userAuthInput(config *AuthConfiguration) {
	fmt.Println("Please enter your username:")
	fmt.Scan(&config.Username)
	fmt.Println("Please enter your password:")
	fmt.Scan(&config.Password)
}

func readFromSaveData() *core.AuthSaveData {
	var saveData *core.AuthSaveData
	saveDataPath := getSaveDataPath()

	if _, err := os.Stat(saveDataPath); os.IsNotExist(err) {
		return nil
	}

	// Save data file exists, read file
	saveData = loadAuthSaveData(saveDataPath)
	oneHourAgo := time.Now().Add(-time.Hour)

	// Compare SavedAt time with one hour ago
	if saveData.SavedAt.Before(oneHourAgo) {
		// discard old save data if it's been more than an hour
		fmt.Println("Last login was more than an hour ago, trying to login again...")
		return nil
	}

	fmt.Printf("Attempting to use previous login tokens...")
	return saveData
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

func getSaveDataPath() string {
	configDir := getConfigDirectory()
	saveDataPath := filepath.Join(configDir, AuthSaveDataFilePath)

	return saveDataPath
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

func saveAuthSaveData(path string, authData *core.AuthSaveData) {
	data, err := json.MarshalIndent(authData, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling auth save data:", err)
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(getConfigDirectory(), 0700)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Println("Error writing auth save data:", err)
	}
}

func loadAuthSaveData(path string) *core.AuthSaveData {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading auth save data:", err)
		return nil
	}

	var config *core.AuthSaveData
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshaling auth save data:", err)
		return nil
	}

	return config
}

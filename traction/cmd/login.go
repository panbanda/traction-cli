package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type tokenResponse struct {
	AccessToken    string `json:"access_token"`
	ExpiresAt      string `json:".expires"`
	Username       string `json:"userName"`
	OrganizationID int    `json:"organization_id"`
	UserID         int    `json:"user_id"`
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to your traction tools account",
	Run: func(cmd *cobra.Command, args []string) {
		validateEmail := func(input string) error {
			_, err := mail.ParseAddress(input)
			if err != nil {
				return errors.New("please provide a valid email address")
			}

			return nil
		}

		emailPrompt := promptui.Prompt{
			Label:    "Email",
			Validate: validateEmail,
		}

		email, err := emailPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		passwordPrompt := promptui.Prompt{
			Label: "Password",
			Mask:  '*',
		}

		password, err := passwordPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		postBody, _ := json.Marshal(map[string]string{
			"grant_type": "password",
			"userName":   email,
			"password":   password,
		})

		responseBody := bytes.NewBuffer(postBody)

		resp, err := http.Post("https://app.bloomgrowth.com/token", "application/json", responseBody)
		if err != nil {
			log.Fatalf("Could not log in: %v", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var payload tokenResponse

		err = json.Unmarshal(body, &payload)
		if err != nil {
			log.Fatalln(err)
		}

		viper.Set("auth.token", payload.AccessToken)
		viper.Set("auth.expires_at", payload.ExpiresAt)
		viper.Set("user.id", payload.UserID)
		viper.Set("user.name", payload.Username)
		viper.Set("organization.id", payload.OrganizationID)
		viper.WriteConfig()

		fmt.Println("Successfully logged in.")
	},
}

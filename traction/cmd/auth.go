package cmd

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

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

		fmt.Printf("You choose %s %s\n", email, password)
	},
}

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/panbanda/traction-cli/internal/session"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

type owner struct {
	ID   int    `json:"Id"`
	Name string `json:"Name"`
}

type todo struct {
	ID         int    `json:"Id"`
	Name       string `json:"Name"`
	DetailsURL string `json:"DetailsUrl"`
	Origin     string `json:"Origin"`
	OriginID   int    `json:"OriginId"`
	DueDate    string `json:"DueDate"`
	Owner      owner  `json:"Owner"`
	Complete   bool   `json:"Complete"`
}

var todoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the associated todos",
	Run: func(cmd *cobra.Command, args []string) {
		req, err := http.NewRequest("GET", "https://app.bloomgrowth.com/api/v1/todo/users/mine", nil)
		if err != nil {
			log.Fatal(err)
		}

		session.ApplyAuthorization(req)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Could not log in: %v", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var payload []todo

		err = json.Unmarshal(body, &payload)
		if err != nil {
			log.Fatalln(err)
		}

		// Output the table
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		tbl := table.New("ID", "Name", "Complete", "Details")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, item := range payload {
			// For some reason they have negative IDs for deleted items
			if item.ID > 0 {
				tbl.AddRow(item.ID, item.Name, item.Complete, item.DetailsURL)
			}
		}

		tbl.Print()
	},
}

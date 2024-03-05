package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Short: "Logger main commands",
	Long:  `The Logger commands to manage log files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hand Shake!")
	},
}

func Exec() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(messageGeneratorCmd)
}

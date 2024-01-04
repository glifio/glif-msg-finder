package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	msgfinder "github.com/jimpick/glif-msg-finder"
	"github.com/spf13/cobra"
)

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "find-messages <agent-id>",
	Short: "Find the messages sent to an agent",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		/*
			maxEpoch, err := cmd.Flags().GetUint64("max-epoch")
			if err != nil {
				log.Fatal(err)
			}
			minEpoch, err := cmd.Flags().GetUint64("min-epoch")
			if err != nil {
				log.Fatal(err)
			}
		*/

		agentID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}

		agentAddress, err := msgfinder.GetAgentAddress(agentID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Address: %v\n", agentAddress)
	},
}

func init() {
	rootCmd.Flags().Uint64("max-epoch", 0, "The minimum epoch")
	rootCmd.Flags().Uint64("min-epoch", 0, "The minimum epoch")
	rootCmd.Flags().Bool("strict", false, "Fail if node doesn't have enough data")
}

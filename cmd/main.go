package main

import (
	"fmt"
	"log"
	"math"
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
		ctx := cmd.Context()

		maxEpoch, err := cmd.Flags().GetUint64("max-epoch")
		if err != nil {
			log.Fatal(err)
		}
		minEpoch, err := cmd.Flags().GetUint64("min-epoch")
		if err != nil {
			log.Fatal(err)
		}

		agentID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}

		agentAddress, err := msgfinder.GetAgentAddress(ctx, agentID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Address: %v\n", agentAddress)

		txs, err := msgfinder.GetTransactions(ctx, agentAddress)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Transactions:")
		seen := make(map[string]bool)
		for _, tx := range txs {
			if tx.Height > maxEpoch {
				continue
			}
			if tx.Height < minEpoch {
				break
			}
			if seen[tx.CID] {
				continue
			}

			txDetail, err := msgfinder.GetTransactionDetail(ctx, tx.SearchID)
			if err != nil {
				log.Fatal(err)
			}
			if txDetail.Level > 0 {
				continue
			}
			seen[tx.CID] = true

			method, params, err := txDetail.ParseParams()
			if err != nil {
				log.Fatal(err)
			}
			if method != nil {
				fmt.Printf("%d %s %s %s %+v\n", tx.Height, tx.CID, tx.Status, method.Name, params)
			}
		}
	},
}

func init() {
	rootCmd.Flags().Uint64("max-epoch", math.MaxUint64, "The minimum epoch")
	rootCmd.Flags().Uint64("min-epoch", 0, "The minimum epoch")
	rootCmd.Flags().Bool("strict", false, "Fail if node doesn't have enough data")
}

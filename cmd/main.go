package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
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
			if method == nil {
				if tx.Type == "Send" {
					fmt.Printf("%d %s incoming %0.2f from %v\n", tx.Height, tx.CID, ToFIL(tx.Amount), tx.From)
				} else {
					fmt.Printf("%d %s unknown: %s\n", tx.Height, tx.CID, tx.Type)
				}
			} else {
				paramStr := ""
				sc, _ := params["sc"].(struct {
					Vc struct {
						Issuer          common.Address "json:\"issuer\""
						Subject         *big.Int       "json:\"subject\""
						EpochIssued     *big.Int       "json:\"epochIssued\""
						EpochValidUntil *big.Int       "json:\"epochValidUntil\""
						Value           *big.Int       "json:\"value\""
						Action          [4]uint8       "json:\"action\""
						Target          uint64         "json:\"target\""
						Claim           []uint8        "json:\"claim\""
					} "json:\"vc\""
					V uint8     "json:\"v\""
					R [32]uint8 "json:\"r\""
					S [32]uint8 "json:\"s\""
				})
				vc := sc.Vc
				switch name := method.Name; name {
				case "addMiner":
					paramStr = fmt.Sprintf("f0%d", vc.Target)
				case "removeMiner":
					paramStr = fmt.Sprintf("f0%d -> New owner: f0%d", vc.Target, params["newMinerOwner"])
				case "pay":
					paramStr = fmt.Sprintf("%0.2f", ToFIL(vc.Value))
				case "borrow":
					paramStr = fmt.Sprintf("%0.2f", ToFIL(vc.Value))
				case "pullFunds":
					paramStr = fmt.Sprintf("%0.2f from f0%d", ToFIL(vc.Value), vc.Target)
				case "pushFunds":
					paramStr = fmt.Sprintf("%0.2f to f0%d", ToFIL(vc.Value), vc.Target)
				case "withdraw":
					paramStr = fmt.Sprintf("%0.2f to %v", ToFIL(vc.Value), params["receiver"])
				case "setRecovered":
				case "refreshRoutes":
				case "confirmChangeMinerWorker":
					paramStr = fmt.Sprintf("f0%d", params["miner"])
				case "changeMinerWorker":
					paramStr = fmt.Sprintf("f0%d -> New worker: f0%v, New control addresses: %v", params["miner"], params["worker"], params["controlAddresses"])
				default:
					paramStr = fmt.Sprintf("%+v", params)
				}
				fmt.Printf("%d %s %s %s\n", tx.Height, tx.CID, method.Name, paramStr)
			}
		}
	},
}

func init() {
	rootCmd.Flags().Uint64("max-epoch", math.MaxUint64, "The minimum epoch")
	rootCmd.Flags().Uint64("min-epoch", 0, "The minimum epoch")
}

func ToFIL(atto *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)
	return f.Quo(f.SetInt(atto), big.NewFloat(params.Ether))
}

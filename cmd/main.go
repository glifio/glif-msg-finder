package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	lotusapi "github.com/filecoin-project/lotus/api"
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
	Use:   "find-blocks <miner-id>",
	Short: "Find the blocks won by a miner",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rpcUrl := cmd.Flag("rpc-url").Value.String()
		maxEpoch, err := cmd.Flags().GetUint64("max-epoch")
		if err != nil {
			log.Fatal(err)
		}
		minEpoch, err := cmd.Flags().GetUint64("min-epoch")
		if err != nil {
			log.Fatal(err)
		}

		minerAddr, err := address.NewFromString(args[0])
		if err != nil {
			log.Fatal(err)
		}
		strict, _ := cmd.Flags().GetBool("strict")

		// Use JSON-RPC API to get miner info
		headers := http.Header{}
		var api lotusapi.FullNodeStruct
		closer, err := jsonrpc.NewMergeClient(context.Background(), rpcUrl, "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
		if err != nil {
			log.Fatalf("connecting with lotus failed: %s", err)
		}
		defer closer()

		head, err := api.ChainHead(cmd.Context())
		if err != nil {
			log.Fatal(err)
		}
		headTSK := head.Key()
		fmt.Println("Height:", head.Height())
		fmt.Println("Tipset:", head.Key())

		maxHeight := abi.ChainEpoch(maxEpoch)
		if maxHeight == 0 {
			maxHeight = head.Height()
		}
		minHeight := abi.ChainEpoch(minEpoch)

		blocks, err := revenue.MinerFindBlocks(cmd.Context(), api, headTSK,
			minerAddr, maxHeight, minHeight, !strict)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Results:\n")
		total := new(big.Int)
		for _, block := range blocks {
			total.Add(total, block.Reward)
			fmt.Printf("%d: %0.09f %s\n", block.Epoch, toFIL(block.Reward), block.CID)
		}
		fmt.Printf("Total: %0.09f\n", toFIL(total))
	},
}

func toFIL(atto *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)
	return f.Quo(f.SetInt(atto), big.NewFloat(1e18))
}

func init() {
	rootCmd.Flags().String("rpc-url", "https://api.node.glif.io/rpc/v1", "Lotus endpoint")
	rootCmd.Flags().Uint64("max-epoch", 0, "The minimum epoch")
	rootCmd.Flags().Uint64("min-epoch", 0, "The minimum epoch")
	rootCmd.Flags().Bool("strict", false, "Fail if node doesn't have enough data")
}

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"

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
	Use:   "find-messages <agent-id>",
	Short: "Find the messages sent to an agent",
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

		agentID, err := strconv.Atoi(args[0])
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

		msgs, err := msgfinder.AgentFindMessages(cmd.Context(), api, headTSK,
			agentID, maxHeight, minHeight, !strict)
		if err != nil {
			log.Fatal(err)
		}
		for _, msg := range msgs {
			fmt.Printf("%+v\n", msg)
		}
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

package msgfinder

import (
	"context"
	"errors"
	"math/big"
	"regexp"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
)

type MinerBlockReward struct {
	Miner    address.Address
	Epoch    abi.ChainEpoch
	WinCount int64
	Reward   *big.Int
	CID      cid.Cid
}

// AgentFindMessages gets messages sent to an agent
func AgentFindMessages(ctx context.Context, api lotusapi.FullNodeStruct,
	headTSK types.TipSetKey, agentID int,
	maxEpoch abi.ChainEpoch, minEpoch abi.ChainEpoch,
	ignoreMissingData bool) ([]MinerBlockReward, error) {

	if maxEpoch-minEpoch <= 1 {
		return []MinerBlockReward{}, errors.New("invalid range")
	}
	var step abi.ChainEpoch = 1000

	checkHeight := maxEpoch

	// Agent 2 (Jim)
	agentAddr, _ := address.NewFromString("f410f6dy45thxrvar53m4ugimu7yvofamzmwtxrc4aaq")

	vestingCID, ts, err := getVestingCIDAndTipSet(ctx, api, agentAddr, headTSK, checkHeight)
	if err != nil {
		return []MinerBlockReward{}, err
	}

	count := 0
	blocks := make([]MinerBlockReward, 0)

	for {
		if checkHeight != ts.Height() {
			// fmt.Printf("%d null round detected\n", checkHeight)
		} else {
			// fmt.Println("Checking at", checkHeight)
			// log.Println("Blocks:")
			for _, block := range ts.Blocks() {
				// log.Printf("  %+v\n", block.Miner)
				if block.Miner == minerAddr {
					count++
					// fmt.Printf("Block %d: %d\n", count, checkHeight)

					winCount := block.ElectionProof.WinCount
					rewardState, err := api.StateReadState(ctx, rewardActor, ts.Key())
					if err != nil {
						return []MinerBlockReward{}, err
					}
					state := rewardState.State.(map[string]interface{})
					thisEpochReward := new(big.Int)
					thisEpochReward.SetString(state["ThisEpochReward"].(string), 10)
					reward :=
						new(big.Int).Div(
							new(big.Int).Mul(thisEpochReward, (new(big.Int).SetInt64(winCount))),
							new(big.Int).SetInt64(5))

					blockReward := MinerBlockReward{
						Miner:    minerAddr,
						Epoch:    checkHeight,
						WinCount: winCount,
						Reward:   reward,
						CID:      block.Cid(),
					}
					blocks = append(blocks, blockReward)
					break
				}
			}
		}
		checkHeight, ts, vestingCID, err = scanPriorVestingTransition(ctx, api,
			minerAddr, headTSK, checkHeight-1, minEpoch, step, ignoreMissingData,
			vestingCID)
		if err != nil {
			return []MinerBlockReward{}, err
		}
		if ts == nil {
			// fmt.Println("No more data.")
			return blocks, nil
		}
	}
}

func scanPriorVestingTransition(ctx context.Context, api lotusapi.FullNodeStruct,
	minerAddr address.Address, headTSK types.TipSetKey, maxEpoch abi.ChainEpoch,
	minEpoch abi.ChainEpoch, step abi.ChainEpoch, ignoreMissingData bool,
	alreadyHaveVestingCID string) (epoch abi.ChainEpoch, ts *types.TipSet,
	vestingCID string, err error) {
	top := maxEpoch
	scanHeight := max(maxEpoch-step, minEpoch)
	var scanVestingCID string
	for {
		scanVestingCID, ts, err = getVestingCIDAndTipSet(ctx, api, minerAddr,
			headTSK, scanHeight)
		if err != nil {
			if ignoreMissingData && isNoDataError(err) {
				// No data on node, search earlier
				// fmt.Printf("No data at epoch %d\n", scanHeight)
				// Use binary search to find later epoch for vesting transition
				// fmt.Printf("Binary search backwards, top: %d bottom: %d\n", top, scanHeight+1)
				laterEpoch, laterTs, laterVestingCID, err := findVestingTransition(ctx,
					api, minerAddr, headTSK, top, scanHeight+1, ignoreMissingData,
					alreadyHaveVestingCID)
				if err != nil {
					return 0, nil, "", err
				}
				if laterEpoch > scanHeight {
					scanHeight = laterEpoch
					scanVestingCID = laterVestingCID
					ts = laterTs
				}
				return scanHeight, ts, scanVestingCID, nil
			}

			return 0, nil, "", err
		}
		// fmt.Printf("Scan %d %s (Top: %d)\n", scanHeight, scanVestingCID, top)
		if scanVestingCID == alreadyHaveVestingCID {
			if scanHeight == minEpoch {
				return 0, nil, "", nil
			}
			top = scanHeight - 1
			scanHeight = max(scanHeight-step, minEpoch)
		} else {
			// Use binary search to find later epoch for vesting transition
			// fmt.Printf("Binary search backwards, top: %d bottom: %d\n", top, scanHeight+1)
			laterEpoch, laterTs, laterVestingCID, err := findVestingTransition(ctx,
				api, minerAddr, headTSK, top, scanHeight+1, ignoreMissingData,
				alreadyHaveVestingCID)
			if err != nil {
				return 0, nil, "", err
			}
			if laterEpoch > scanHeight {
				scanHeight = laterEpoch
				scanVestingCID = laterVestingCID
				ts = laterTs
			}
			return scanHeight, ts, scanVestingCID, nil
		}
	}
}

func findVestingTransition(ctx context.Context, api lotusapi.FullNodeStruct,
	minerAddr address.Address, headTSK types.TipSetKey, maxEpoch abi.ChainEpoch,
	minEpoch abi.ChainEpoch, ignoreMissingData bool, alreadyHaveVestingCID string) (
	height abi.ChainEpoch, ts *types.TipSet, vestingCID string, err error) {
	var foundEpoch abi.ChainEpoch
	var foundTs *types.TipSet
	var foundVestingCID string
	top := maxEpoch
	bottom := minEpoch
	for bottom < top {
		middle := bottom + (top-bottom)/2
		vestingCID, ts, err := getVestingCIDAndTipSet(ctx, api, minerAddr, headTSK, middle)
		if err != nil {
			if ignoreMissingData && isNoDataError(err) {
				// No data on node, search earlier
				// fmt.Printf("Binary search: %d %d %d no data, searching higher\n", top, middle, bottom)
				bottom = middle + 1
				continue
			}
			return 0, nil, "", err
		}
		// fmt.Printf("Binary search: %d %d %d %s\n", top, middle, bottom, vestingCID)
		if vestingCID == alreadyHaveVestingCID {
			// fmt.Println("Already have, searching lower")
			top = middle
		} else {
			// fmt.Printf("Found new vesting data at %d, searching higher\n", middle)
			foundEpoch = middle
			foundTs = ts
			foundVestingCID = vestingCID
			bottom = middle + 1
		}
	}

	return foundEpoch, foundTs, foundVestingCID, nil
}

func getVestingCIDAndTipSet(ctx context.Context, api lotusapi.FullNodeStruct,
	minerAddr address.Address, headTSK types.TipSetKey, height abi.ChainEpoch) (
	string, *types.TipSet, error) {
	ts, err := api.ChainGetTipSetByHeight(ctx, height, headTSK)
	if err != nil {
		return "", nil, err
	}
	state, err := api.StateReadState(ctx, minerAddr, ts.Key())
	if err != nil {
		return "", nil, err
	}
	vestingFundsCID := state.State.(map[string]interface{})["VestingFunds"]
	vestingFunds := vestingFundsCID.(map[string]interface{})["/"].(string)
	return vestingFunds, ts, nil
}

var noDataRegex *regexp.Regexp
var glifLimitRegex *regexp.Regexp

func isNoDataError(err error) bool {
	return noDataRegex.MatchString(err.Error()) ||
		glifLimitRegex.MatchString(err.Error())
}

func init() {
	noDataMatchStr := "^getting actor: load state tree: failed to load state tree \\S+: failed to load hamt node: ipld: could not find \\S+$"
	noDataRegex = regexp.MustCompile(noDataMatchStr)

	glifLimitMatchStr := "^bad tipset height: lookbacks of more than .* are disallowed$"
	glifLimitRegex = regexp.MustCompile(glifLimitMatchStr)
}

package msgfinder

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

type AgentRecord struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
}

const agentURL = "https://events.glif.link/agent"

func GetAgentAddress(ctx context.Context, agentID int) (common.Address, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, agentURL, nil)
	if err != nil {
		return common.Address{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return common.Address{}, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return common.Address{}, err
	}

	var agents []AgentRecord

	err = json.Unmarshal(resBody, &agents)
	if err != nil {
		return common.Address{}, err
	}

	for _, agent := range agents {
		if agent.ID == agentID {
			return common.HexToAddress(agent.Address), nil
		}
	}

	return common.Address{}, errors.New("agent id not found")
}

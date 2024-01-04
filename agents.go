package msgfinder

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type AgentRecord struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
}

const agentURL = "https://events.glif.link/agent"

func GetAgentAddress(agentID int) (string, error) {
	req, err := http.NewRequest(http.MethodGet, agentURL, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var agents []AgentRecord

	err = json.Unmarshal(resBody, &agents)
	if err != nil {
		return "", err
	}

	for _, agent := range agents {
		if agent.ID == agentID {
			return agent.Address, nil
		}
	}

	return "", errors.New("agent id not found")
}

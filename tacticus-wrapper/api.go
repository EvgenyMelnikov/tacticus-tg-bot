package tacticus_wrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	yandex "tacticus-tg-bot/yandex-tables-wrapper"
)

func FetchPlayer(apiKey string) (*TacticusPlayerResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.tacticusgame.com/api/v1/player", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Api-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Tacticus API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var playerResp TacticusPlayerResponse
	if err := json.Unmarshal(body, &playerResp); err != nil {
		return nil, err
	}

	return &playerResp, nil
}

func FetchPlayers(players []yandex.Player) []PlayerData {
	var tacticusPlayers []PlayerData

	for _, player := range players {
		resp, err := FetchPlayer(player.Api)
		if err != nil {
			log.Println(player.Name, err)
			continue
		}

		tacticusPlayers = append(tacticusPlayers, resp.Player)
	}

	sort.Slice(tacticusPlayers, func(i, j int) bool {
		return tacticusPlayers[i].Progress.GuildRaid.BombTokens.NextTokenInSeconds <
			tacticusPlayers[j].Progress.GuildRaid.BombTokens.NextTokenInSeconds
	})

	return tacticusPlayers
}

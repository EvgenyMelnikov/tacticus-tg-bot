package yandex_tables_wrapper

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/xuri/excelize/v2"
)

func ReadPlayersFromXLSX(file *excelize.File) ([]Player, error) {
	rows, err := file.GetRows("Лист1")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var players []Player

	for i, row := range rows {
		if i == 0 {
			continue // заголовок
		}

		players = append(players, Player{
			Name: row[0],
			Id:   row[1],
			Api:  row[2],
		})
	}

	return players, nil
}

func GetYandexDownloadLink(publicLink string) (string, error) {
	api := "https://cloud-api.yandex.net/v1/disk/public/resources/download?public_key=" +
		url.QueryEscape(publicLink)

	resp, err := http.Get(api)
	if err != nil {
		log.Panicln(err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var data YandexDownload
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Panicln(err)
		return "", err
	}

	return data.Href, nil
}

func LoadXLSX(url string) (*excelize.File, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	tmpFile, err := os.CreateTemp("", "*.xlsx")
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	defer func(name string) {
		_ = os.Remove(name)
	}(tmpFile.Name())

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}

	return excelize.OpenFile(tmpFile.Name())
}

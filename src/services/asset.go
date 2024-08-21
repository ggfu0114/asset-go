package assetserv

import (
	dbmodel "asset-go/src/models"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const AssetPyHost = "http://127.0.0.1:8080"

type CurrentValueResponse struct {
	CurrentValue float32 `json:"current_value"`
}
type CurrentValueRequest struct {
	Type   string  `json:"type"`
	Market string  `json:"market"`
	Code   string  `json:"code"`
	Amount float32 `json:"amount"`
}

func QueryAssetValue(asset dbmodel.Asset) {

	jsonData := CurrentValueRequest{
		Type:   asset.AssetType,
		Market: asset.Market,
		Code:   asset.Code,
		Amount: asset.Amount,
	}

	jsonByte, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest(http.MethodGet, AssetPyHost+"/current_value", bytes.NewBuffer(jsonByte))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	//Convert the body to type string
	sb := string(body)
	log.Println("response Body:", string(sb))
	var result CurrentValueResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("Can not unmarshal JSON")
	}
	log.Println("result value:", result.CurrentValue)
	asset.Value = result.CurrentValue
	dbmodel.UpdateCurrentValue(asset)
}

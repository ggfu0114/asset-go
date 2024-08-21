package dbmodel

import (
	"log"
	"strconv"
	"time"
)

type Asset struct {
	Aid          int
	Uid          int
	AssetType    string
	AssetSubType string
	Code         string
	Amount       float32
	Label        string
	Market       string
	Value        float32
	UpdatedAt    time.Time
}

func ListMyAsset(Uid string) []Asset {
	var assets []Asset
	db := GetDb()
	rows, err := db.Query(`SELECT
								m.*,
								av.Value,
								av.UpdatedAt
							FROM
								myasset m
							INNER JOIN asset_value av ON
								m.Aid = av.Aid
							WHERE
								Uid = ?`, Uid)
	if err != nil {
		log.Println("Failed to query DB", err)
	}
	defer rows.Close()

	for rows.Next() {
		var asset Asset
		err := rows.Scan(
			&asset.Aid, &asset.AssetType, &asset.AssetSubType, &asset.Code, &asset.Amount, &asset.Label, &asset.Market, &asset.Uid, &asset.Value, &asset.UpdatedAt)
		if err != nil {
			log.Fatalln(err)
		} else {
			assets = append(assets, asset)
		}

	}
	return assets
}
func AddAsset(Uid string, asset Asset) string {
	// TODO: Trigger request to get current value.
	db := GetDb()
	stmt, err := db.Prepare("INSERT INTO myasset SET AssetType=?, AssetSubType=?, Code=?, Amount=?, Label=?, Market=?, Uid=?")
	if err != nil {
		return ""
	}

	res, queryError := stmt.Exec(asset.AssetType,
		asset.AssetSubType, asset.Code, asset.Amount, asset.Label, asset.Market, Uid)
	id, err := res.LastInsertId()
	var aid string
	if err != nil {
		log.Println("Error:", err.Error())
	} else {
		log.Println("LastInsertId:", id)
		aid = strconv.FormatInt(id, 10)
	}
	if queryError != nil {
		log.Fatalln(queryError)
		return ""
	}
	return aid
}

func UpdateAsset(uid string, aid string, asset Asset) int {
	db := GetDb()
	stmt, err := db.Prepare("UPDATE myasset SET AssetType=?, AssetSubType=?, Code=?, Amount=?, Label=?, Market=? WHERE Aid=? AND Uid=?")
	if err != nil {
		return -1
	}
	res, queryError := stmt.Exec(asset.AssetType,
		asset.AssetSubType, asset.Code, asset.Amount, asset.Label, asset.Market, aid, uid)
	log.Println("res", res)
	if queryError != nil {
		log.Fatalln(queryError)
		return -1
	}
	return 1
}

func DeleteAsset(uid string, aid string) bool {
	db := GetDb()
	_, err = db.Exec("DELETE FROM `myasset` WHERE `aid` = ? AND `uid`=?;", aid, uid)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	return true
}

func UpdateCurrentValue(asset Asset) bool {

	db := GetDb()
	log.Println("asset.Aid:", asset.Aid)
	log.Println("asset.Value:", asset.Value)
	_, err = db.Exec(`
		INSERT INTO asset_value 
		(Aid, Value) VALUES(?, ?) 
		ON DUPLICATE KEY UPDATE  Aid=?
	`, asset.Aid, asset.Value, asset.Aid)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	return true
}

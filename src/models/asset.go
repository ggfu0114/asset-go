package dbmodel

import (
	"log"
	"strconv"
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
}

func ListMyAsset(Uid string) []Asset {
	var assets []Asset
	db := GetDb()
	rows, err := db.Query("SELECT * FROM myasset WHERE Uid=?", Uid)
	if err != nil {
		log.Println("Failed to query DB", err)
	}
	defer rows.Close()

	for rows.Next() {
		var asset Asset
		err := rows.Scan(
			&asset.Aid, &asset.AssetType, &asset.AssetSubType, &asset.Code, &asset.Amount, &asset.Label, &asset.Market, &asset.Uid)
		if err != nil {
			log.Fatalln(err)
		} else {
			assets = append(assets, asset)
		}
	}
	return assets
}
func AddAsset(Uid string, asset Asset) string {
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

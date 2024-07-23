package dbmodel

import(
	"log"
	"strconv"
)

type Asset struct{
	Aid int
	AssetType string
	AssetSubType string
	Code string 
	Amount float32
	Label string
	Market string
}

func ListMyAsset() ([]Asset) {
	var assets []Asset
	db := GetDb()
	rows, err := db.Query("SELECT * FROM myasset")
	if err!= nil{
		log.Println("Failed to query DB", err)
	}
	for rows.Next() {
		var asset Asset
		err := rows.Scan(
			&asset.Aid, &asset.AssetType, &asset.AssetSubType, &asset.Code, &asset.Amount, &asset.Label, &asset.Market)
		if err != nil {
			log.Fatalln(err)
		}else{
			assets=append(assets,asset)
		}	
	}
	return assets
}
func AddAsset(asset Asset)(string){
	db := GetDb()
	stmt, err := db.Prepare("INSERT INTO myasset SET AssetType=?, AssetSubType=?, Code=?, Amount=?, Label=?, Market=?")
	if err != nil {
		return ""
	}
	
	res, queryError := stmt.Exec(asset.AssetType, 
		asset.AssetSubType, asset.Code, asset.Amount, asset.Label, asset.Market)
	id, err := res.LastInsertId()
	var aid string
	if err != nil {
		log.Println("Error:", err.Error())
	} else {
		log.Println("LastInsertId:", id)
		aid =  strconv.FormatInt(id, 10)
	}
	if queryError != nil {
		log.Fatalln(queryError)
		return ""
	}	
	return aid
}

func UpdateAsset(aid string, asset Asset)(int){
	db := GetDb()
	stmt, err := db.Prepare("UPDATE myasset SET AssetType=?, AssetSubType=?, Code=?, Amount=?, Label=?, Market=? WHERE Aid=?")
	if err != nil {
		return -1
	}
	res, queryError := stmt.Exec(asset.AssetType, 
		asset.AssetSubType, asset.Code, asset.Amount, asset.Label, asset.Market, aid)
	log.Println("res", res)
	if queryError != nil {
		log.Fatalln(queryError)
		return -1
	}	
	return 1
}
func DeleteAsset(aid string)(bool){
	db := GetDb()
	_, err=db.Exec("DELETE FROM `myasset` WHERE `aid` = ?;", aid)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	return true
}
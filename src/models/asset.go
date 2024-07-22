package dbmodel

import(
	"log"
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
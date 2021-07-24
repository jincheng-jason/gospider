package main

import (
	"context"
	"log"
	"time"

)

func main() {

	//init context
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	

	url := "https://www.bunnings.co.nz/products/building-hardware/timber/mouldings/primed-mouldings"
	selector := "div.SearchComponent__StyledSearchComponent-sc-17x15nm-0.dEnjhg > div.SearchComponentstyle__SearchComponentWrapper-sc-1l60lhw-11.iXiHCd > article"
	
	items, err := GetItems(url, selector)
	if err != nil {
		log.Fatalf("get data error %v\n", err)
	}
	log.Printf("Items: %v\n", len(items))

	//init mongo, get collection
	mongoUrl := "mongodb://root:root@localhost:27020/price?authSource=admin"
	collection, err := InitMongo(ctx, mongoUrl, 10)
	if err != nil {
		log.Fatalf("init mongo error %v\n", err)
	}
	res, err := collection.InsertMany(ctx, items)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

}

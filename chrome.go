package main

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetItems(url string, selector string) ([]interface{}, error) {
	// create context
	c, _ := chromedp.NewExecAllocator(
		context.Background(), 
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
		)...,
	)

	chromeCtx, _ := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	// defer cancel()

	timeoutCtx, _ := context.WithTimeout(chromeCtx, 40 * time.Second)
	// defer timeoutCancel()

	log.Printf("chrome visit page%sn",url)

	var htmlContent string
	//first page
	err := chromedp.Run(timeoutCtx, VisitWeb(url, selector, &htmlContent))
	if err != nil {
		log.Printf("Run err : %v\n", err)
		return nil, err
	}

	items, err := GetDataFromContent(htmlContent)
	if err != nil {
		log.Fatalf("get data error %v\n", err)
	}
	log.Printf("Items: %v\n", len(items))

	//loop next page
	nextPageSel := "#__next > main > div > div:nth-child(3) > div.SearchComponent__StyledSearchComponent-sc-17x15nm-0.dEnjhg > div.container.paginationContainer > div > div > nav > ul > li:nth-child(4) > button"
	if CheckNextPage(htmlContent, nextPageSel) {
		var newContent string
		err := chromedp.Run(timeoutCtx, 
			chromedp.Click(nextPageSel, chromedp.NodeVisible),
			chromedp.Sleep(1 * time.Second),
			scrollAndRetrive(selector, &newContent),
		)
		if err != nil {
			log.Printf("Run err : %v\n", err)
			return nil, err
		}
		newItems, err := GetDataFromContent(newContent)
		if err != nil {
			log.Fatalf("get new data error %v\n", err)
		}
		items = append(items, newItems...)
	}


	return items, nil
}

func CheckNextPage(content string, nextPageSel string) bool {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatalln(err)
		return false
	}
	if s := dom.Find(nextPageSel); s != nil {
		return true
	} else {
		return false
	}
}

func VisitWeb(url string, selector string, htmlContent *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		scrollAndRetrive(selector, htmlContent),
	}
}

func scrollAndRetrive(selector string, htmlContent *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitVisible(selector),
		chromedp.ActionFunc(func(c context.Context) error {
			_, exp, err := runtime.Evaluate(`window.scrollTo({top: document.body.scrollHeight, behavior: 'smooth'});`).Do(c)
                if err != nil {
                    return err
                }
                if exp != nil {
                    return exp
                }
                return nil
		}),
		chromedp.Sleep(1 * time.Second),
		chromedp.ActionFunc(func(c context.Context) error {
			_, exp, err := runtime.Evaluate(`window.scrollTo({top: document.body.scrollHeight*0.75, behavior: 'smooth'});`).Do(c)
                if err != nil {
                    return err
                }
                if exp != nil {
                    return exp
                }
                return nil
		}),
		chromedp.Sleep(2 * time.Second),
		chromedp.ActionFunc(func(c context.Context) error {
			_, exp, err := runtime.Evaluate(`window.scrollTo({top: document.body.scrollHeight, behavior: 'smooth'});`).Do(c)
                if err != nil {
                    return err
                }
                if exp != nil {
                    return exp
                }
                return nil
		}),
		chromedp.Sleep(1 * time.Second),
		chromedp.OuterHTML(`document.querySelector("body")`, htmlContent, chromedp.ByJSPath),
	}
}

func GetDataFromContent(content string) ([]interface{}, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	items := []interface{}{}
	articleSelector :=     "div.SearchComponent__StyledSearchComponent-sc-17x15nm-0.dEnjhg > div.SearchComponentstyle__SearchComponentWrapper-sc-1l60lhw-11.iXiHCd > article"
	lazyArticleSelector := "div.SearchComponent__StyledSearchComponent-sc-17x15nm-0.dEnjhg > div.SearchComponentstyle__SearchComponentWrapper-sc-1l60lhw-11.iXiHCd > div > article"
	dom.Find(articleSelector).Add(lazyArticleSelector).Each(func(i int, s *goquery.Selection) {
		p1 := strings.TrimPrefix(s.Find("div.Pricestyle__PriceWrap-sc-kv48nd-0.dARVlW.search-price.price-medium-size > p").Text(), "$")
		price, err := strconv.ParseFloat(p1, 64)
		if err != nil {
			log.Fatalf("Error when parse price: %v\n", err)
		}
		img := s.Find("div.SearchProductTilestyle__ImageRatingContainer-sc-7jrh24-6.AmGqC > div.product-wrapper > div > a > figure > picture > img")
		imgUrl, _ := img.Attr("src")
		brandUrl, _ := s.Find("div.SearchProductTilestyle__ImageRatingContainer-sc-7jrh24-6.AmGqC > div.product-wrapper > div > div > figure > picture > img").Attr("src")
		log.Printf("price: %s", p1)
		items = append(items, Item{
			ID: primitive.NewObjectID(),
			Title: s.Find("div.SearchProductTilestyle__ImageRatingContainer-sc-7jrh24-6.AmGqC > div.text-rating-container > a > p").Text(),
			Price: price,
			ImgUrl: imgUrl,
			BrandImgUrl: brandUrl,
		})
	})
	return items, nil
}
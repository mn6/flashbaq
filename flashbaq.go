package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis"
	_ "github.com/joho/godotenv/autoload"
)

var symbolBase = "https://www.nasdaq.com/aspx/infoquotes.aspx"
var symbolErrMsg = "Quotes are not supported for this symbol."
var chartBase = "https://charting.nasdaq.com/ext/charts.dll?2-1-14-0-0-512-03NA000000"
var chartSuff = "-&SF:1|5-BG=FFFFFF-BT=0-HT=395--XTBL-"
var cacheTime = 300

type chartRet struct {
	Date   string `json:"date"`
	Last   string `json:"last"`
	Volume string `json:"volume"`
}

type symbolRet struct {
	Ticker        string     `json:"ticker"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	Market        string     `json:"market"`
	MarketStatus  string     `json:"marketStatus"`
	Price         string     `json:"price"`
	Change        string     `json:"change"`
	PercentChange string     `json:"percentChange"`
	ShareVolume   string     `json:"shareVolume"`
	TodaysHigh    string     `json:"todaysHigh"`
	TodaysLow     string     `json:"todaysLow"`
	BestBid       string     `json:"bestBid"`
	FiftyTwoHigh  string     `json:"fiftyTwoWeekHigh"`
	FiftyTwoLow   string     `json:"fiftyTwoWeekLow"`
	EPS           string     `json:"earningsPerShare"`
	OpenPrice     string     `json:"openPrice"`
	ClosePrice    string     `json:"closePrice"`
	ChartData     []chartRet `json:"chartData"`
	Website       string     `json:"website"`
}

type symbolJSON struct {
	Result []symbolRet `json:"result"`
}

var client = &http.Client{}
var rds = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func main() {
	port := os.Getenv("FLASHBAQ_PORT")
	if len(port) < 1 {
		port = ":7810"
	}
	origins := os.Getenv("FLASHBAQ_ALLOWED_ORIGINS")
	if len(origins) < 1 {
		origins = "*"
	}

	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{origins},
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge:         cacheTime,
	})
	r.Use(cors.Handler)

	r.Get("/symbol", symbol)

	log.Printf("Listening on port %s...", port)
	http.ListenAndServe(port, r)
}

func symbol(w http.ResponseWriter, r *http.Request) {
	tickers := r.URL.Query().Get("tickers")
	if len(tickers) < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{ error: \"NO TICKERS PROVIDED\" }"))
	} else {
		tickers = "&" + strings.Join(strings.Split(tickers, ","), "&")
		get := getDB(tickers)
		if len(get) > 1 {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(get))
			log.Println("RServed " + tickers)
			return
		}
		req, err := http.NewRequest("GET", symbolBase, nil)
		chk(err)
		req.Header.Add("cookie", "userSymbolList="+tickers)
		resp, err := client.Do(req)
		var results = []symbolRet{}
		body, err := goquery.NewDocumentFromReader(resp.Body)
		chk(err)
		symbolScrape(body, &results)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
		log.Println("Served " + tickers)
		out, err := json.Marshal(results)
		chk(err)

		setDB(tickers, string(out))
	}
}

func symbolScrape(doc *goquery.Document, res *[]symbolRet) {
	doc.Find("table#quotes_content_left_InfoQuotesResults > tbody > tr").Each(func(i int, s *goquery.Selection) {
		info := strings.Split(s.Find(".infoquote_qn > div").Eq(1).Text(), "|")
		ticker := cleanScrape(info[0])
		name := info[1]
		stocktype := info[2]
		market := strings.Replace(info[3], "Market : ", "", -1)
		marketstatus := s.Find("[id=\"" + ticker + "_MarketStatus\"]").Text()
		price := cleanMoney(s.Find(".lastsale_qn > label").Text())
		changeSpan := s.Find(".lastsale_qn > span")
		change := changeSpan.Find("label").Eq(0).Text()
		percentchange := changeSpan.Find("label").Eq(1).Text()
		if changeSpan.Is(".red") {
			change = "-" + change
			percentchange = "-" + percentchange
		}
		website := s.Find(".iqsumlinku").Text()
		tableInfo := s.Find(".infoquote_qn .genTable tr td")
		shareVolume := tableInfo.Eq(3).Text()
		todaysHigh := cleanMoney(tableInfo.Eq(5).Text())
		todaysLow := cleanMoney(tableInfo.Eq(5).Text())
		bestBid := cleanMoney(tableInfo.Eq(7).Text())
		fiftyTwoHigh := cleanMoney(tableInfo.Eq(9).Text())
		eps := cleanMoney(tableInfo.Eq(11).Text())
		openPrice := cleanMoney(tableInfo.Eq(13).Text())
		closePrice := cleanMoney(tableInfo.Eq(15).Text())
		fiftyTwoLow := cleanMoney(tableInfo.Eq(25).Text())
		chartData := chartScrape(ticker)

		cleanSymbolScrape(
			&name, &stocktype, &market, &marketstatus,
			&price, &change, &percentchange, &shareVolume,
			&todaysHigh, &bestBid, &fiftyTwoHigh, &fiftyTwoLow,
			&eps, &openPrice, &closePrice, &todaysLow,
		)
		*res = append(*res, symbolRet{
			Ticker:        ticker,
			Name:          name,
			Website:       website,
			Type:          stocktype,
			Market:        market,
			MarketStatus:  marketstatus,
			Price:         price,
			Change:        change,
			PercentChange: percentchange,
			ShareVolume:   shareVolume,
			TodaysHigh:    todaysHigh,
			TodaysLow:     todaysLow,
			BestBid:       bestBid,
			FiftyTwoHigh:  fiftyTwoHigh,
			FiftyTwoLow:   fiftyTwoLow,
			EPS:           eps,
			OpenPrice:     openPrice,
			ClosePrice:    closePrice,
			ChartData:     chartData,
		})
	})
}

func chartScrape(ticker string) []chartRet {
	var chart []chartRet
	req, err := http.NewRequest("GET", chartBase+ticker+chartSuff, nil)
	chk(err)
	resp, err := client.Do(req)
	body, err := goquery.NewDocumentFromReader(resp.Body)

	body.Find(".DrillDown > tbody").Find("tr").Not(":nth-child(1)").Each(func(i int, s *goquery.Selection) {
		data := s.Find("td")
		chart = append(chart, chartRet{
			Date:   data.Eq(0).Text(),
			Last:   data.Eq(1).Text(),
			Volume: data.Eq(2).Text(),
		})
	})
	for i := len(chart)/2 - 1; i >= 0; i-- {
		opp := len(chart) - 1 - i
		chart[i], chart[opp] = chart[opp], chart[i]
	}
	return chart
}

func cleanScrape(field string) string {
	return strings.Trim(strings.Replace(field, "\n", "", -1), " ")
}

func cleanSymbolScrape(fields ...*string) {
	for _, field := range fields {
		*field = strings.Trim(strings.Replace(*field, "\n", "", -1), " ")
	}
}

func remSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func cleanMoney(str string) string {
	return remSpace(strings.Replace(str, "$", "", -1))
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func getDB(ticker string) string {
	get, err := rds.Get("flashbaq:" + ticker).Result()
	if err != nil && err.Error() == "redis: nil" {
		get = ""
	} else {
		chk(err)
	}
	return get
}

func setDB(ticker string, data string) {
	err := rds.Set("flashbaq:"+ticker, data, time.Duration(cacheTime)*time.Second).Err()
	chk(err)
}

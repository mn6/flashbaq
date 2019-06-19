# flashbaq
Nasdaq scraper + API built in Go.

## Usage
```
http://localhost:7810
  /symbol                  Returns stock symbol information.
    ?tickers=goog,nvda     Comma-delimitered ticker list.
```

### Example data returned

```json
[
  {
    "ticker": "TSLA",
    "name": "Tesla, Inc.",
    "website": "https://www.tesla.com",
    "type": "Common Stock",
    "market": "NASDAQ-GS",
    "marketStatus": "Delayed Data - Jun. 18, 2019 Market Closed",
    "price": "224.74",
    "change": "-0.29",
    "percentChange": "-0.13%",
    "shareVolume": "12,715,788",
    "todaysHigh": "234.74",
    "todaysLow": "234.74",
    "bestBid": "N/A",
    "fiftyTwoWeekHigh": "387.46",
    "fiftyTwoWeekLow": "176.9919",
    "earningsPerShare": "-5.79",
    "openPrice": "228.72",
    "closePrice": "224.74",
    "chartData": [
      {
        "date": "06/18/2018",
        "last": "370.83",
        "volume": "12,025,450"
      },
      {
        "date": "06/19/2018",
        "last": "352.55",
        "volume": "12,734,840"
      },
      ...
    ],
    "news": [
      {
        "heading": "5 Top Stock Trades for Tuesday: FB, NFLX, GILD, TSLA",
        "url": "https://www.nasdaq.com/article/5-top-stock-trades-for-tuesday-fb-nflx-gild-tsla-cm1165025",
        "details": "6/17/2019 8:05:26 PM - InvestorPlace Media"
      },
      {
        "heading": "Why Tesla, C&J Energy Services, and CrowdStrike Holdings Jumped Today",
        "url": "https://www.nasdaq.com/article/why-tesla-cj-energy-services-and-crowdstrike-holdings-jumped-today-cm1164970",
        "details": "6/17/2019 8:29:00 PM - Motley Fool"
      },
      ...
    ]
  }
]
```

## Setup

1. Create a copy of `.env.example` and rename it `.env`
2. Edit your `.env` file variables as you wish
3. Ensure `redis` is installed and running.
4. Run built binary
5. Create systemd unit or equivalent

## Building

1. Ensure Go is installed
2. `go get` this repository
3. `cd` to the repository path (`$GOPATH/src/github.com/mn6/flashbaq`)
4. `go build` to build
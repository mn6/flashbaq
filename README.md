# flashbaq
Nasdaq scraper + API built in Go.

## Usage
```
http://localhost:7810
  /symbol                  Returns stock symbol information.
    ?tickers=goog,nvda     Comma-delimitered ticker list.
http://localhost:7810
  /chart                  Returns stock symbol information.
    ?ticker=amzn          One ticker per request.
```

### Example data returned

```/symbol?tickers=TSLA```
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

```/chart?ticker=TSLA```
```json
[
  {
    "date": "03/21/2019",
    "last": "1819.26",
    "volume": "5,740,075",
    "open": "1796.26",
    "high": "1823.75",
    "low": "1787.281"
  },
  {
    "date": "03/22/2019",
    "last": "1764.77",
    "volume": "6,347,858",
    "open": "1810.17",
    "high": "1818.98",
    "low": "1763.11"
  },
  {
    "date": "03/25/2019",
    "last": "1774.26",
    "volume": "5,098,182",
    "open": "1757.79",
    "high": "1782.6751",
    "low": "1747.5"
  },
  {
    "date": "03/26/2019",
    "last": "1783.76",
    "volume": "4,848,654",
    "open": "1793",
    "high": "1805.77",
    "low": "1773.3598"
  },
  {
    "date": "03/27/2019",
    "last": "1765.7",
    "volume": "4,316,646",
    "open": "1784.13",
    "high": "1787.5",
    "low": "1745.68"
  },
  {
    "date": "03/28/2019",
    "last": "1773.42",
    "volume": "3,030,860",
    "open": "1770",
    "high": "1777.93",
    "low": "1753.47"
  },
  ...
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
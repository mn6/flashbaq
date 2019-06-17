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
    "ticker": "GOOG",
    "name": "Alphabet Inc.",
    "type": "Class C Capital Stock",
    "market": "NASDAQ-GS",
    "marketStatus": "Delayed Data - Jun. 14, 2019 Market Closed",
    "price": "1,085.35",
    "change": "-3.42",
    "percentChange": "-0.31%",
    "shareVolume": "1,111,643",
    "todaysHigh": "1,092.69",
    "todaysLow": "1,092.69",
    "bestBid": "N/A",
    "fiftyTwoWeekHigh": "1,289.27",
    "fiftyTwoWeekLow": "970.11",
    "earningsPerShare": "39.87",
    "openPrice": "1,086.42",
    "closePrice": "1,085.35",
    "chartData": [
      {
        "date": "06/14/2018",
        "last": "1152.12",
        "volume": "1,350,085"
      },
      {
        "date": "06/15/2018",
        "last": "1152.26",
        "volume": "2,119,134"
      },
      ...
    ]
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
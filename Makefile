run:
	go run ./...

req:
	curl "http://localhost:5000/run?productID=PI_XBTUSD&period=candles_trade_1m&size=1&limitCoef=0.1"

stop:
	curl "http://localhost:5000/stop"

test:
	go test -v ./...

cover:
	go test -cover -coverprofile=coverage.out ./... && go tool cover -func=coverage.out && rm coverage.out

math.randomseed(os.time())

wrk.method = "GET"

local symbols = {"AAPL", "GOOG", "MSFT", "AMZN", "TSLA", "META", "NFLX", "NVDA", "IBM", "INTC"}

request = function()
  local symbol = symbols[math.random(1, #symbols)]
  local depth = math.random(5, 50)
  local path = string.format("/api/v1/orderbook/%s?depth=%d", symbol, depth)
  return wrk.format("GET", path)
end

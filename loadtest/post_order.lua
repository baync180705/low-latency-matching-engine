math.randomseed(os.time())

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

local symbols = {"AAPL", "GOOG", "MSFT", "AMZN", "TSLA", "META", "NFLX", "NVDA", "IBM", "INTC"}
local sides = {"BUY", "SELL"}
local orderTypes = {"LIMIT", "MARKET"}

request = function()
  local symbol = symbols[math.random(1, #symbols)]
  local side = sides[math.random(1, #sides)]
  local orderType = orderTypes[math.random(1, #orderTypes)]
  local price = math.random(10000, 10500)
  local qty = math.random(1, 500)

  local body = string.format(
    '{"symbol":"%s","side":"%s","type":"%s","price":%d,"quantity":%d}',
    symbol, side, orderType, price, qty
  )

  return wrk.format("POST", "/api/v1/orders", nil, body)
end

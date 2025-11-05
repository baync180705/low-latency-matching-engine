math.randomseed(os.time())

wrk.method = "GET"

request = function()
  local id = string.format("test-%d", math.random(1, 100000))
  return wrk.format("GET", "/api/v1/orders/" .. id)
end

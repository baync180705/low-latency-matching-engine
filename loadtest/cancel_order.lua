math.randomseed(os.time())

wrk.method = "DELETE"
wrk.headers["Content-Type"] = "application/json"

request = function()
  local id = string.format("test-%d", math.random(1, 100000)) -- the order id is wrong, I am just using this for loadtesting, it will eventually return 404
  return wrk.format("DELETE", "/api/v1/orders/" .. id)
end

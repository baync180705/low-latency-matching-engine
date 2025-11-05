#!/bin/bash

URL="http://localhost:8080"


echo "Starting load test for order submission..."
wrk -t3 -c50 -d15s -s ./loadtest/post_order.lua $URL #only testing the post order endpoint since it is the primary endpoint.

echo -e "\nFetching engine metrics..."
curl -s $URL/metrics | jq

echo -e "\nDone."

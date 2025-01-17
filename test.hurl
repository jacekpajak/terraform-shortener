# Test POST Request - Shorten URL
POST https://pu378ggj02.execute-api.eu-central-1.amazonaws.com/shorten
Content-Type: application/json

{
  "url": "https://example.com"
}

# HTTP/1.1 200
# [Asserts]
# jsonpath "$.short_url" exists
# jsonpath "$.original_url" == "https://example.com"



# # Test GET Request - Redirect URL
# GET https://pu378ggj02.execute-api.eu-central-1.amazonaws.com/<short-url>

# HTTP/1.1 301
# [Asserts]
# header "Location" == "https://example.com"

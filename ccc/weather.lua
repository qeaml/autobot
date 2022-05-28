-- Simple weather command that sends an HTTP request to https://wttr.in/
-- Usage: $weather <location>
-- Example: $weather Warsaw

if #args < 1 then
  say("Provide a location.")
  return
end

typing()

local strings = require ("strings")
local location = strings.join(args, "+")
local url = "https://wttr.in/" .. location .. "?format=4"
local http = require ("http")
local r = http.get(url)
if r == nil then
  say("HTTP error.")
  return
end
say(r.text)

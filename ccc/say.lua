-- Sends a message to a predefined webhook using the sender's name and avatar.
-- Usage: $say <anything>
-- Example: $say deez nuts

local http = require("http")
local str = require("strings")
-- replace this with your own discord webhook
local r = http.post("<discord webhook url>", "application/json", '{"content":"' .. str.join(args, " ") .. '","username":"' .. msg.author.name .. '","avatar_url":"' .. msg.author.avatar .. '"}')
if r.status ~= 204 then
  say("Error code: " .. r.status)
end

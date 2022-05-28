-- Extremely simple alternative to "Google it" actually sending a link to a
-- Google search page with the given term.
-- Usage: $lmgtfy <search term>
-- Example: $lmgtfy how to make a bomb

local strings = require("strings")
say("https://google.com/search?q=" .. strings.join(args, "+"))

# autobot

autobot is a somewhat advanced Discord bot written in Go using [discordgo](github.com/bwmarrin/discordgo).

## Features of autobot

* A command system split into three types of commands:
  * **Internal commands**, defined in the bot's source code.
  * **Simple commands**, defined using the bot itself. Simply replace all occuernces of `%s` in the command's string with arguments that the command was executed with. (e.g. a simple command called `gotcha` defined as `haha! i got your %s!` will turn into `haha! i got your nose!` when executed as `$gotcha nose`)
  * **Complex commands**, defined using Lua code. The bot has a somewhat extensive API for this Lua VM that allows it to create some pretty cool stuff without having to hardcore it into the bot. Examples of complex commands can be (eventually) found in the `ccc` directory.

* A simple text generation algorithm which learns on every message that the bot receives. It can use the data it gathered to generate text using the `!generate` internal command.

* A nice terminal user interface, using only one line of text for all the information you'll ever need.

## Getting started

1. Obtain a binary of the bot.
  This can be done by downloading a binary release of the bot or by cloning the repository and building yourself. The building process is as simple as running `go build`, the Go tool will take care of everything else. If you are building from source, I recommend using `-o bin/bot` to place the executable in the `bin` directory.

2. Find a good place for the binary.
  The bot will place all of it's files in the same directory as it's binary. So, if you get a bot binary `bot.exe` and place it in a directory called `autobot`, all of the bot's files -- the commands, the text generation model and everything else -- will also be placed in that directory. Same goes for the `config.yml` file. Speaking of which,

3. Create a `config.yml` file.
  Use the template below to create the config file for the bot.

4. Configure.
  Use the table below to change configuration settings in your template.

5. Run the bot.
  Run your bot's executable and use `!commands` (or whatever other internal command prefix you chose) to see if it's working. Press `Ctrl+C` at any moment in the window to shut the bot down. **IMPORTANT**: Do **not** close the bot by closing the window: this will cause all of it's data to **not** be saved -- that's your text generation data, your commands and other stuff -- remember to **always** press `Ctrl+C` to stop the bot before closing the window.

### `config.yml` template

```yml
token: "deez nuts"
admin: "396699211946655745"
nprefix: '!'
cprefix: '$'
```

### `config.yml` fields

* `token` - your bot's Discord token, obtained from the [Discord Developer Portal](discord.com/developers).
* `admin` - Discord ID of a user who will have access to certain admin-only command.
* `nprefix` - prefix used for internal commands. I use an exclamaition mark (`!`) for this, and as such it is used that way in the docs.
* `cprefix` - prefix used for simple and complex commands. I use a dollar sign (`$`) for this and it is used in the docs for simple and complex commands.

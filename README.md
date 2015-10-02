<h1>Twitch IRC Bot</h1>
Created in Go!


Original: https://stackoverflow.com/questions/13342128/simple-golang-irc-bot-keeps-timing-out

<h2>Features:</h2>
**Console Input** - Don't feel like going in your browser/IRC client? You can type your input into the program and it comes out if the bot were saying it! - although the text creates a newline from both os.Stdin and net.Conn (bufio.NewReader) 
**Nickname coloring**
This uses github.com/shiena/ansicolor
**Logging**
Functions for this are readfile and writefile (readfile meant to hook it up to a chatbot using python but project died - I recommend looking at https://github.com/aichaos/rivescript-go)

<h4>How to use</h4>
```
git clone github.com/pressure679/Twitch-IRC-Bot
or Download ZIP -> On the right sidebar

Create twitch_pass and insert the oauth from
http://twitchapps.com/tmi/

Then either "go run bot.go" or "go build bot.go"
```

Please add/remove/modify as you please!  :D

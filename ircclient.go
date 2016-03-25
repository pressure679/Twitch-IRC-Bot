//    Copyright (C) 2014  Vittus Peter Ove Maqe Mikiassen
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <http://www.gnu.org/licenses/>.

/*
My own Twitch bot
Twitch bot guide at https://github.com/justintv/Twitch-API/blob/master/IRC.md
Used as example:
http://github.com/Vaultpls/Twitch-IRC-Bot/blob/master/bot.go
- doc at https://gowalker.org/github.com/Vaultpls/Twitch-IRC-Bot
- - Added name coloring and removed automsg
http://github.com/paked/go-twitchbot/
- api for this is at https://github.com/thoj/go-ircevent
*/

package main
import (
	"fmt"
	"strings"
	"flag"
	"time"
	"io/ioutil"
	"bufio"
	"os"
	"net"
	"net/textproto"
	"github.com/shiena/ansicolor"
)

func main() {

	// commands available to configure irc client
	server := flag.String("server", "irc.twitch.tv", "The server to chat in.")
	port := flag.String("port", "6667", "The port to connect to server with.")
	channel := flag.String("channel", "#channel", "Sets the channel for the bot to go into.")
	nick := flag.String("nick", "naamik01", "The username of the account.")
	clientid := flag.String("clientid", "", "The client ID to use, see Settings - Connections in your account for yours.")
	maxmsg := flag.Int("spamtime", 5, "Set a minimum time between messages.")

	pass1, err := ioutil.ReadFile("twitch_pass.txt")
	if err != nil {
		fmt.Println("Error reading from twitch_pass.txt.  Maybe it isn't created?")
		os.Exit(1)
	}

	// configure irc client settings
	flag.Parse()
	ircbot := NewBot()
	ircbot.server = *server
	ircbot.port = *port
	ircbot.channel = *channel
	ircbot.nick = *nick
	ircbot.pass = = strings.Replace(string(pass1), "\n", "", 0)
	ircbot.clientid = *clientid
	ircbot.maxmsg = int64(*maxmsg)

	// start irc client
	go ircbot.ConsoleInput()
	ircbot.Connect()
	defer ircbot.socket.Close()
	fmt.Fprintf(ircbot.socket, "PASS %s\r\n", ircbot.pass)
	fmt.Fprintf(ircbot.socket, "NICK %s\r\n", ircbot.nick)
	fmt.Fprintf(ircbot.socket, "JOIN %s\r\n", ircbot.channel)
	reader := bufio.NewReader(ircbot.socket)
	tp := textproto.NewReader(reader)

	// For writing msg's from irc
	var offsetmsg, offsetmsgfromuser int
	var msg, msgfromuser string
	var mymsg1, mymsg2 string

	// For colouring msg's
	var colorname = make(map[string]int)
	var names []string
	var color int = 2
	colors := [...]string{"\x1b[31m", "\x1b[33m", "\x1b[35m", "\x1b[36m"}
	w := ansicolor.NewAnsiColorWriter(os.Stdout)

	// for chatbot
	var lines []string

	// start the bot
	for {
		line, err1 := tp.ReadLine()
		if err1 != nil {
			continue
		}
		
		// if irc msg contains the stream's' viewer's name's
		if strings.Contains(line, "353") {
			names = strings.Split(line[len(ircbot.channel) + len(ircbot.nick) * 2 + 25:], " ")
			fmt.Println("353: names")
			fmt.Println(names)
			fmt.Println()
			for x := 0; x < len(names); x++ {
				color = colorincr(color)
				colorname[names[x]] = color
			}
			continue

		// if a viewer joins the channel 
		} else if strings.Contains(line, "JOIN") {
			color = colorincr(color)
			colorname[line[1:strings.Index(line, "!")]] = color
			continue

		// if a ping msg is send from the server sned a pong
		} else if strings.Contains(line, "PING") {
			pongdata := strings.Split(line, "PING ")
			fmt.Fprintf(ircbot.socket, "PONG %s\r\n", pongdata[0])
			continue

		// To check if a command you send is invalid
		} else if strings.Contains(line, "421") {
			fmt.Println(line)

		// if a private msg is send, write the msg to console
		} else if strings.Contains(line, "PRIVMSG") {
			offsetmsg = strings.Index(line, "PRIVMSG") + len(ircbot.channel) + 10
			offsetmsgfromuser = strings.Index(line, "!")
			msg = line[offsetmsg:]
			msgfromuser = line[1:offsetmsgfromuser]
			txt := time.Now().String()[11:19] + ";" + "%s" + msgfromuser + "%s" + ": " + msg + "\n"

			// to decide which color to assign a name
			if colorname[msgfromuser] != 0 {
				fmt.Fprintf(w, txt, colors[colorname[msgfromuser] - 1], "\x1b[0m")
			} else { /* just o be sure msg is written*/
				names = append(names, msgfromuser)
				color = colorincr(color)
				colorname[msgfromuser] = color
				fmt.Fprintf(w, txt, colors[colorname[msgfromuser] - 1], "\x1b[0m")
			}

			// writefile("test.txt", msg)
			if (ircbot.msgcount % 2) == 0 {
				time.Sleep(15 * 1000000000)
			} else {
				time.Sleep(5 * 1000000000)
			}
			lines = readfile("test.txt")
			mymsg1 = lines[len(lines) - 1]
			mymsg2 = lines[len(lines) - 2]
			// loop through msg's and check latest msg to send a new one.
			for 
			ircbot.Msg(mymsg1)
			time.Sleep(time.Duration(len(lines[len(lines) - 2]) / 3))
			ircbot.Msg(mymsg2)
			ircbot.msgcount = ircbotmsgcount + 2
		}
	}
}
// Bot core
type Bot struct {
	server, port, channel string
	nick, pass, clientid string
	socket net.Conn
	command net.Conn
	lastmsg, maxmsg int64
	msgcount int
}
// New bot method
func NewBot() *Bot {
	return &Bot{
		server:       "irc.twitch.tv",
		port:         "6667",
		channel:      "#kinulii",
		nick:         "naamik01",
		pass:         "",
		clientid:     "",
		socket:       nil,
		command:      nil,
		lastmsg:      0,
		maxmsg:       5,
		msgcount:     0,
	}
}
// Connect to channel method
func (ircbot *Bot) Connect() {
	var err error
	fmt.Printf("Attempting to connect to server...\n")
	ircbot.socket, err = net.Dial("tcp", ircbot.server + ":" + ircbot.port)
	if err != nil {
		fmt.Printf("ircbot.socket Unable to connect to Twitch IRC server! Reconnecting in 10 seconds...\n")
		time.Sleep(10 * time.Second)
		ircbot.Connect()
	}
	fmt.Printf("Connected to IRC server %s\n", ircbot.server)
}
// Msg the channel method - to send msg's
func (ircbot *Bot) Msg(msg string) {
	if msg == "" {
		return
	}
	w := ansicolor.NewAnsiColorWriter(os.Stdout)
	fmt.Fprintf(ircbot.socket, "PRIVMSG " + ircbot.channel + " :"	+ msg	+ "\r\n")
	msg = "Me: %s" + msg + "%s" + "\n"
	fmt.Fprintf(w, msg, "\x1b[37m", "\x1b[0m")
	ircbot.lastmsg = time.Now().Unix()
	ircbot.msgcount++
}
// Console input method - to write msg's
func (ircbot *Bot) ConsoleInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if text == "" {
			return
		} else if text == "/quit" {
			fmt.Fprintf(ircbot.socket, "PART %s\r\n", ircbot.channel)			
			ircbot.socket.Close()
			fmt.Println("Quit")
			os.Exit(0)
		} else if strings.Contains(text, "/join ") {
			fmt.Fprintf(ircbot.socket, "PART %s\r\n", ircbot.channel)
			fmt.Fprintf(ircbot.socket, "JOIN %s\r\n", text[5:])
			ircbot.channel = text[5:]
			fmt.Println("Joined " + ircbot.channel)
		}
		if text != "" {
			ircbot.Msg(text)
			ircbot.msgcount++
		}
	}
}
// method to check for errors in a dumb-safe way
func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}
// Loop through colors to assign them somewhat randomly to names in IRC
func colorincr(n int) int {
	n++
	if n == 4 {
		n = 1
	}
	return n
}
// func read file
// optimize to read file and not create new file reader
// every method call
func readfile(filename string) []string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(file), "\n")
	return lines
}
func writefile(filename, txt string) {
	f, err := os.OpenFile(filename, os.O_APPEND, 0644)
	checkerr(err)
	_, err = f.WriteString("\n" + "go" + txt)
	checkerr(err)
	f.Close()
}

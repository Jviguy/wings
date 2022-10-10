# Wings
Wings is a library for [Dragonfly](https://github.com/df-mc/dragonfly), an asynchronous Minecraft server software written in Go. 

Wings implements auto completion for commands and displays them as they would be seen ingame. It also implements a console command sender/reader.

## Usage
Using wings is simple. Start it after registering your commands and calling server.Listen()

wings.New(server, log, wings.DefaultConfig()).Start()

## Preview
![IMAGE](https://media.discordapp.net/attachments/745099242904748153/1026840478693466223/unknown.png)

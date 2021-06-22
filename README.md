
# Minecraft release bot

Telegram bot that notifies when a new version of minecraft is released.
Check it out at https://t.me/minecraft_update

## Environment Variables

To run this project, you will the following environment variables. Can be passed as cli arguments

`MINE_TOKEN` Telegram bot token

`MINE_CHANNEL` Telegram chat id to send notifications

`MINE_LOG_CHANNEL` Telegram chat id to send log

  
## Usage

```bash
minecraft -token abcd -channel 123456789
```
  
## Run Locally

Clone the project

```bash
  git clone https://github.com/MrMarble/minecraft-update-go
```

Go to the project directory

```bash
  cd minecraft-update-go
```

Install dependencies

```bash
  go get
```

Run process

```bash
  task run -- -channel 123456789 -token abcdefg
```

  
## Handle Twitch Chat

Create in path main `.env` file

Example `.env` file:

```bash
DATABASE_URL=your_mongo_db_url
OAUTH_TOKEN=oauth:xxxxxxxxxxxxxxxxxxxxxxx
CHANNEL_NAME=your_channel_name
DATABASE_NAME=you_database_name
TWITCH_CLIENT_ID=your_client_id
TWITCH_CLIENT_SECRET=your_client_secret
```

OAuth token you receive from
[Twitch's OAuth Generator](http://twitchapps.com/tmi/)

```bash
$ go mod tidy
```

```bash
$ go run main.go
```
test

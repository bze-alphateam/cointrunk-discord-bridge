# CoinTrunk.io Discord bot

This repo contains a simple application written in Go that will listen to Tendermint websocket of BZE for new articles published in CoinTrunk module and post them to a discord webhook. 

It requires the following configuration file: 
```yaml
rest_url: "https://rest.getbze.com"
rpc_url: "rpc.getbze.com"
discord:
  webhook: DISCORD_WEBHOOK_HERE
history:
  filename: history.json
  filepath: ./
```

After compilation run:  
`ctbot start`  
If the configuration is valid it will start the program, fetch articles, check in local file `history.json` what's the last article it fetched.
If there's no history file it will publish last 50 articles found. 
If the file is present and contains the key `last_id` it will filter out all articles that have an id lower or equal with the provided value.
Once all articles are fetched and posted on discord it will save the latest posted article ID in `history.json`.
After everything is done it will connect to tendermint websocket using the `rpc_url` provided in the configuration file and when a new article is published it will post it on discord.

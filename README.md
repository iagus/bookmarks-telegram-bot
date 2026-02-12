### Bookmarks Telegram bot

A Telegram bot for saving bookmarks— great for personal self-hosting on a home server or VPS.

Minimal hassle. Keeps track of interesting links without dealing with specific browser bookmarks.

I built this to solve a particular issue of mine— I am good at accumulating information but bad at
keeping it tidy and organized. This is me trying to manage all the scattered links I have within
a single place.

##### Wishlist

- [ ] Restart the systemd services if npm tasks are successfull (I forgot oops)
- [ ] Fetch OpenGraph information thyself! I don't want to be rate limited by an external service!
- [ ] Better placeholder if the OpenGraph image for the link can't be fetched or whatever other reason

#### What does it do?

This bot listens for your messages in Telegram and saves the URLs in them in a text file you
must specify. It comes with a server that will process the text file, build an HTML, and serve it.

The cron is running frequently in my server. Checks for updates in the GitHub repository, and if
there are any, it runs the test, pass the linter (Prettier would be the only dependency, and
arguably could be just managed in the development phase but meh I was having a blast doing this part).

#### Installation

As mentioned, I host this myself by managing the Go and the Node services separately using `systemd`
in my VPS server. You do you! If you want to go down the same path as me:

##### Requirements

- A Telegram bot token (get one from @BotFather)
- An environment file with your secrets (example included)
- Go toolchain (if building from source)

##### Clone the repo

```
git clone https://github.com/iagus/bookmarks-telegram-bot.git
cd bookmarks-telegram-bot
```

##### Create your secrets.env and fill the info:

```
cp secrets.example.env secrets.env
```

##### Build Go

`go run main.go`

##### Start the Node server

```
cd server
npm run start
```

#### Deploy

Depends on what route you go. If, like me, you're in for `systemd`, I recommend
modifying the `cron.sh` file to your liking and adding it to your crontab.
Then you'll need two separate `systemd` services- one for the Go process, and
another for Node's. Remember to enable them once created. Logs are accessible
in (and followable) in nice format by running `journalctl -fu <service name> -o cat`.

There are examples of the service files. Watch out with the ownership and
file permissions- make sure your user can execute them:

###### Go

```
[Unit]
Description=Bookmarks Telegram Bot service
After=network.target

[Service]
EnvironmentFile=path/to/your/env/file
ExecStart=/path/to/compiled/go/binary
WorkingDirectory=/path/to/cloned/repo
Restart=no <- you choose
User=your user

[Install]
WantedBy=multi-user.target
```

###### Node

Remember that the `node` binary must be accessible for `systemd`.

```
[Unit]
Description=Bookmarks Telegram Bot writer Reader

[Service]
EnvironmentFile=path/to/your/env/file
ExecStart=/path/to/node /path/to/bookmarks.js
WorkingDirectory=/path/to/cloned/repo/server
Restart=on-failure <- you choose
User=your user

[Install]
WantedBy=multi-user.target
```

##### Serving styles

I did not implement serving the styles from Node, since I have an nginx vhost
up and running for this app server.

You can do the same, or modify the `bookmarks.js` script to handle the
styles request. But if you're using nginx/apache, I'd recommend serving
it from there.

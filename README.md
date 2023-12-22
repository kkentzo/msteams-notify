# msteams-notify

This programs published a message to an MS Teams webhook (Incoming
WebHook connector) as an adaptive card.

The card includes:

- a message title
- a message body
- a cc section with mentions (for notifying the corresponding users)

All of the above (except the mentions which are optional) must be
supplied from the command line.

The program's options (also available using the `-h` switch) are as follows:

```
Usage of msteams-notify:
  -body string
    	The message body (mandatory)
  -mentions string
    	An optional comma-separated link of mentions (in the form of emails) to be included in the message
  -title string
    	The message title (mandatory)
  -uri string
    	MS Teams Webhook URI (mandatory)
```

The program can be built locally by running:

```bash
$ go build -o msteams-notify
```

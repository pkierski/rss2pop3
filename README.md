# rss2pop3

Simple POP3 server with RSS news as emails.

This is an example for [pop3srv](https://github.com/pkierski/pop3srv) module.

## Usage

### Build
Just compile with 
```sh
go build .
```

### Start server
```sh
./rss2pop3
```

Server starts listening on `:pop3` address (which is usually port 110) by default.
You can specify one or more listening addres(es) with option `-p`, ex.:
```sh
./rss2pop3 -p 127.0.0.0:110
```

### Configuring POP3

Server uses user name as list of RSS/Atom addresses, separated by pipe (`|`). ex:
`https://medium.com/feed/tag/go|https://medium.com/feed/tag/pop3`.


Password is ignored, you can set any non-empty string.

# remarkable-sync

**remarkable-sync** is a Golang tool design to run on reMarkable paper tablet for syncing your document from different services.

## Supported services

* [Wallabag](https://www.wallabag.org)
* [Miniflux](https://miniflux.app)

## Build

Build for run on reMarkable tablet :

~~~
env GOOS=linux GOARCH=arm GOARM=7 go build -o remarkable-sync
~~~

## Connection

After configuring Wifi on your reMarkable paper tablet, you can access to it with SSH:

~~~
ssh root@192.168.X.X
~~~

*// Password and IP address are indicated in the end of **About** menu in reMarkable settings !*

## Install

For install, you just need to upload **remarkable-sync** to your reMarkable:

~~~
ssh root@192.168.X.X 'mkdir -p /usr/local/bin'
scp remarkable-sync root@192.168.X.X:/usr/local/bin
~~~

## Config

### Wallabag

Edit */etc/remarkable-sync/wallabag.json* on your reMarkable:

~~~
{
    "host": "app.wallabag.it",
    "client_id": "client_id",
    "client_secret": ""client_secret,
    "username": "login",
    "password": "password"
}
~~~

### Miniflux

Edit */etc/remarkable-sync/miniflux.json* on your reMarkable:

~~~
{
    "host": "app.miniflux.net",
    "token": "token"
}
~~~

## Usage

Connect to the tablet, run **remarkable-sync**, then restart **xochitl**:

~~~
ssh root@192.168.X.X
remarkable-sync
systemctl restart xochitl
~~~

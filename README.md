# reMarkable-sync

**reMarkable-sync** is a Golang tool design to run on reMarkable paper tablet for syncing your document from different services.

## Supported services

Currently the only supported service is [Wallabag](https://www.wallabag.org).

## Build

Build for run on reMarkable tablet :

~~~
env GOOS=linux GOARCH=arm GOARM=7 go build -o reMarkable-sync
~~~

## Connection

After configuring Wifi on your reMarkable paper tablet, you can access to it with SSH:

~~~
ssh root@192.168.X.X
~~~

*// Password and IP address are indicated in the end of **About** menu in reMarkable settings !*

## Install

For install, you just need to upload **reMarkable-sync** to your reMarkable:

~~~
scp reMarkable-sync root@192.168.X.X:
~~~

## Config

### Wallabag

Edit *~/.config/reMarkable-sync/wallabag.json* on your reMarkable:

~~~
{
    "host": "app.wallabag.it",
    "client_id": "client_id",
    "client_secret": ""client_secret,
    "username": "login",
    "password": "password"
}
~~~

## Usage

Connect to the tablet, run **reMarkable-sync**, then restart **xochitl**:

~~~
ssh root@192.168.X.X
./reMarkable-sync && systemctl restart xochitl
~~~

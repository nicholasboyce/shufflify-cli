# Shufflify CLI

A command line usable version of https://www.shufflify.app written in Go.
shufflify-cli uses OAuth PKCE to consume the Spotify API. It accepts as arguments the titles of the Spotify playlists you wish to shuffle together into one big playlist. It allows you to title the new shuffled mega-playlist and posts it to your Spotify account. The playlists are private by default. When inputting the title of the playlist, surround multi word titles with quotes, e.g. "My example playlist", and be mindful of escaping characters, e.g. "Playlist\!" instead of "Playlist!"

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

As of right now, the only official way to install this program is to have Go installed on your machine. 

For ease, I've linked the options to download from:

The [official Go website](https://go.dev/doc/install)

If you're using MacOS, you can use Homebrew to download Go as well.

```
brew install go
```

### Installing

Using shufflify-cli is easy! It can be installed by running:

```
go install github.com/nicholasboyce/shufflify-cli
```

### Usage Examples

#### Inputting playlist names

To shuffle together (at least 2) playlists, use the template below

```
shufflify-cli "Playlist 1" "Workout Bops" playlist2
```

As you can see, it is safe to use quotes around arguments that contain spaces alongside one word titles without quotes.

#### To get help with commandline flags

```
shufflify-cli --help
```

or 

```
shufflify-cli -h
```

#### Using flags

shufflify-cli creates a config file in order to help reduce the work it takes to login. It saves the information necessary to replicate the client, so all you need to do to log in is input your Client ID for the app and follow the instructions that appear in the terminal.

To logout (delete the config file at the current specified location), you can either manually delete the config JSON file shufflify-cli will look for OR use the -logout flag. 

```
shufflify-cli -logout
```

NOTE: This deletes the config file at the path specified by the env variable PATH_TO_CONFIG, which you should set before running, either by editing the .env file (this persists) or on the command line (sets the path only for one session):

```
PATH_TO_CONFIG=/path/to/config/here shufflify-cli
```

 If you do not set it before running, it defaults to ./config.json.


If you'd like to change the location where shufflify-cli looks for/saves the config file for ONE session, you can also use the -filepath flag as seen below. To persist this change, edit the .env directly.

If you change the location of the config file using -filepath, it's advisable to delete the current config file first, as it will not be deleted just by changing the config path - shufflify-cli will create a new file in the designated area. You can run -logout to delete the original config file before using -filepath.

```
shufflify-cli -filepath=/path/to/config/here.json
```

-logout and -filepath can be used in combination - logout will look for and delete the path specified in -filepath.

```
shufflify-cli -filepath=/path/to/config/here.json -logout
```


## Running the tests

If you've cloned this repo, you can simply run

```
go test
```

to run the tests provided.


## Author

[Nicholas Boyce](https://github.com/nicholasboyce)
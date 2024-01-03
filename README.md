# Shufflify CLI

A command line usable version of https://www.shufflify.app written in Go.
shufflify-cli uses OAuth PKCE to consume the Spotify API. It accepts as arguments the titles of the Spotify playlists you wish to shuffle together into one big playlist. It allows you to title the new shuffled mega-playlist and posts it to your Spotify account. The playlists are private by default. When inputting the title of the playlist, surround multi word titles with quotes, e.g. "My example playlist", and be mindful of escaping characters, e.g. "Playlist\\!" instead of "Playlist!"

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

This app works on free and premium Spotify accounts. If you'd like to queue the tracks shuffled by shufflify-cli from this app, you'll need a premium account. Otherwise, you can just press play on the playlist in Spotify.

As of right now, the only official way to install this program is to have Go installed on your machine. 

For ease, I've linked the options to download from:

The [official Go website](https://go.dev/doc/install)

If you're using MacOS, you can use Homebrew to download Go as well.

```
brew install go
```

You'll also need your own Client ID, which you can get by making an account at https://developer.spotify.com and then following their instructions in the [Create An App section](https://developer.spotify.com/documentation/web-api/tutorials/getting-started#create-an-app).

### Installing

Using shufflify-cli is easy! It can be installed by running:

```
go install github.com/nicholasboyce/shufflify-cli
```

If you'd like to edit the .env file however, it's advisable that you clone the repo instead so that you have access to the code.

```
git clone github.com/nicholasboyce/shufflify-cli <target-directory>
```

From there, you can build the executable with

```
go build
```

or install it so that it's globally accessible with 

```
go install
```

### Usage Examples

#### Logging In

shufflify-cli will prompt you to log in if your access token has expired, if you've never logged in before, or if you've changed where it should look for the app config file.

You log in simply by running 

```
shufflify-cli
```

or by running it with arguments and/or flags as seen below.

You can pass in your Client ID as an environment variable before running if you'd like:

```
CLIENT_ID=client_id_here shufflify-cli
```

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

When logged in shufflify-cli will by default save the JSON file it needs in a folder relative to its installation location. It really is the 'state' of the app (whether it's logged in or not). If that really upsets you, you have the freedom to change where this JSON file is saved. HOWEVER, The program may run in unintended ways if you aren't careful about where you save your files. Inputting your own path means relative paths will be relative to your current working directory. User beware!!!

If you'd like to change the location where shufflify-cli looks for/saves the config file for ONE session, you can use the env variable PATH_TO_CONFIG:

```
PATH_TO_CONFIG=/path/to/config/here shufflify-cli
```

You can also use the -filepath flag as seen below instead of the env variable.

```
shufflify-cli -filepath=/path/to/config/here.json
```

To logout (delete the config file at the current specified location), you can either manually delete the config JSON file shufflify-cli will look for OR use the -logout flag. 

```
shufflify-cli -logout
```

If you do not specify a path before running, it defaults to ./shufflify/config.json.


If you change the location of the config file, it's advisable to delete the current config file first, as it will not be deleted just by changing the config path - shufflify-cli will simply create a new file in the designated area. You can run -logout to delete the original config file before using -filepath.

-logout can be used in combination with -filepath or PATH_TO_CONFIG - logout will look for and delete the path specified.

```
shufflify-cli -filepath=/path/to/config/here.json -logout
```

OR 

```
PATH_TO_CONFIG=/path/to/config/here shufflify-cli -logout
```

## Running the tests

If you've cloned this repo, you can simply run

```
go test
```

to run the tests provided.


## Author

[Nicholas Boyce](https://github.com/nicholasboyce)
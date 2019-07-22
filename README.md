# Bitrise IO plugin for [Bitrise CLI]

A Terminal / Command Line interface for bitrise.io, to manage your apps on bitrise.io right from your terminal / command line.

## How to install

You can install this [Bitrise CLI](https://www.bitrise.io/cli) plugin using the bitrise cli:

```
bitrise plugin install https://github.com/bitrise-io/bitrise-plugins-io.git
```

## How to use

First you have to sign in / authenticate yourself.
Authentication requires a Bitrise.io Personal Access Token,
which you can generate at: [https://www.bitrise.io/me/profile#/security](https://www.bitrise.io/me/profile#/security).

Once you have the Personal Access Token register it for the plugin using the `auth` command:

```
bitrise :io auth --token 3c..NQ
```

You can check with which user you're signed into the plugin with the `auth whoami` command:

```
$ bitrise :io auth whoami

viktorbenei
```

List your apps with the `apps` command:

```
bitrise :io apps
```

List builds of a specific app (identified by its slug/ID you see in the `apps` list):

```
bitrise :io builds --app 72b...392
```

Limit it to just the last 5 builds:

```
bitrise :io builds --app 72b...392 --limit=5
```

For all the available commands use:

```
bitrise :io help
```

And to get all the available flags of a given command:

```
bitrise :io COMMAND --help
```

e.g. for the `builds` command to see all of its available options:

```
bitrise :io builds --help
```


## Development

### How to create a new release of this plugin

- bump `RELEASE_VERSION` in bitrise.yml
- commit the change(s)
- call `bitrise run create-release`
- check and update the generated `CHANGELOG.md`
- test the generated binaries in `_bin/` directory
- push these changes to the master branch
- once `deploy` workflow finishes on bitrise.io create a github release with the generated binaries

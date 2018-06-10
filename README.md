# Bitrise IO plugin for [Bitrise CLI]

A Terminal / Command Line interface for bitrise.io, to manage your apps on bitrise.io right from your terminal / command line.

## How to install

You can install this [Bitrise CLI](https://www.bitrise.io/cli) plugin using the bitrise cli:

```
bitrise plugin install https://github.com/bitrise-core/bitrise-plugins-io.git
```


## Development

### How to create a new release of this plugin

- bump `RELEASE_VERSION` in bitrise.yml
- comit these change
- call `bitrise run create-release`
- check and update the generated `CHANGELOG.md`
- test the generated binaries in `_bin/` directory
- push these changes to the master branch
- once `deploy` workflow finishes on bitrise.io create a github release with the generated binaries

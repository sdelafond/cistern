# Changelog

## Next version

### Features

* User interface: Column visibility, column order, sort order and all style attributes can now be set
in the configuration file (see the [sample configuration file included](https://github.com/nbedos/cistern/blob/master/cmd/cistern/cistern.toml) in the release archives for details)
* User interface: Add new columns: created, finished, xfail and url
* User interface: Support horizontal scrolling of table rows 
* User interface: Support sorting pipelines by any column 
* Core: Monitor all remotes of local repositories instead of just 'origin' ([issue #3](https://github.com/nbedos/cistern/issues/3))
* Configuration: Support setting API token from output of user-specified process ([issue #13](https://github.com/nbedos/cistern/issues/13))

### Bug fixes

* Lookup path if BROWSER is not a path itself ([issue #8](https://github.com/nbedos/cistern/issues/8))
* Add support for 'insteadOf' and 'pushInsteadOf' configuration for remote URLs  ([issue #7](https://github.com/nbedos/cistern/issues/7))
* Stages of Azure pipelines are now ordered

### Chores

* Refactor table widget for improved maintainability
* Rewrite build script in go for improved maintainability
* Rename repository


## Version 0.1.2 (2019-12-20)

* Fix: Binary releases now contain an executable built for the right system ([issue #4](https://github.com/nbedos/cistern/issues/4))
* Fix: Appveyor pipelines triggered for a tag incorrectly showed a branch as reference instead of a tag


## Version 0.1.1 (2019-12-18)

* Support for private GitLab instances ([issue #2](https://github.com/nbedos/cistern/issues/2))


## Version 0.1.0 (2019-12-18)
Initial release!

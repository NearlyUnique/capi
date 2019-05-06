# Command-Line API

Circle CI: [![CircleCI](https://circleci.com/gh/NearlyUnique/capi.svg?style=svg)](https://circleci.com/gh/NearlyUnique/capi)

## Overview

`capi` is a tool to simplify ad-hoc api testing where you might use `curl` with or without a wrapping shell script. Commonly used values can be collected from the command line, a `profile.json` file or input from some other `json` file. To help further bash ato complete is built in and is dynamic based on the api you have configured and are trying to call.

The `profile.json` file (or a file path held in env var `CAPI_PROFILE`) holds all api information. The profile lists apis, each api lists commands.

### A simple profile.json

```json
{
    "apis": [{
        "name":"httpbin",
        "baseUrl":"https://httpbin.org",
        "defaultHeader": {
            "any-header":"any-value {arg1}"
        },
        "commands": [{
            "name":"status",
            "path":"/status/{code}",
            "defaultHeader": {
                "x-a-header":"{arg2}"
            }
        }]
    }]
}
```

### Basic usage

```bash
export ARG1=value
capi httpbin status --code 418
```

For a configured api, `httpbin`, execute the `status` command replacing parameter `code` with `418`. `arg1` can be supplied as a cli arg but in this example is not, however is is loaded as an environment variables. CLI parameters override environment. variables. `arg2` does not get replaced because there is no CLI parameter or matching environment variable.

An api may have default header values, these are specified at the api level, they can be overridden for any command within that api.

#### Environment variables

- Environment variables is **case insensitive**.
- Where required a root `profile.json` value `envPrefix` can be added. Any env var with this prefix will be used in preference to one without. e.g. if `"envPrefix":"XYZ_"`, then `XYT_USER` would be used in preference to `USER` for parameter `{user}`
- cli arguments override environment variables.

## Installation (for auto complete on bash)

```bash
go get -u https://github.com/NearlyUnique/capi
```

### Option 1
```bash
# put capi in the path, then
complete -C capi capi
```

### Option 2
```bash
complete -C /path/to/capi capi
```

## Uninstall

```bash
complete -r capi
```


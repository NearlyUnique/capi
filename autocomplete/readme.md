# Complete (in bash)

## This package

This package performs 2 simple tasks.

1. Are all the expected inputs available to run as a `complete` command.
1. Expose those values as a go struct

## Basic inputs

The `complete`bash command is fairly straightforward but was not obvious to integrate with

### Installation

```bash
complete -C command name
```

where `-C` specifies a command (rather than a function or other type), `command` is the command buing 'completed' and `name` is the name of the command performing the completion suggestions. 

### CLI arguments

`os.Args` always returns at least one argument, the name of the command/program itself.

The next 3 arguments are from the `complete` command.

| # | Meaning  |
|---|----------|
| 1 | Command  |
| 2 | Word     |
| 3 | PrevWord |

### Environment Variables

| Variable     | Description                   |
|--------------|-------------------------------|
| `COMP_LINE`  | The whole line                |
| `COMP_POINT` | Position in the `COMP_LINE`   |
| `COMP_KEY`   | Key to start the completion   |
| `COMP_TYPE`  | Type of completion, see below |

#### `COMP_TYPE` values

| ASCII | Dec |                     Description                     |
| ----- | --- | --------------------------------------------------- |
| TAB   | 9   | for normal completion                               |
| ‘?’,  |     | for listing completions after successive tabs       |
| ‘!’,  |     | for listing alternatives on partial word completion |
| ‘@’,  |     | to list completions if the word is not unmodified   |
| ‘%’,  |     | for menu completion                                 |


## Outputs

1. Each line is an option
1. When a single result is returned it is used to **replace** the user entered value
1. Output is sorted by `complete`
1. Colours are not supported in bash
1. If all results start with the same _prefix_ then complete will automatically replace the users word with that prefix

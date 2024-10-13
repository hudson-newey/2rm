# TODO

## Additional command line arguments

- `--progress` Display a progress bar
- `--undo=n` Replace _n_ with the number of soft-deletes to recover
- `--defer` Delete files/directories that match the glob after a certain period of time (e.g. `2rm --expire-in=2D` *.bak)

### Smart File Filters

- `--older-than` Delete all files/directories created less than a period of time (e.g. `2rm --older-than=30d *.log`)
- `--larger-than` Delete all files/directories that match the specifier glob (e.g. `2rm --larger-than=1G *.iso`)
- `--smaller-than` Delete all files/directories that are smaller than a certain size (e.g. `2rm --smaller-than=1B *`)

## Features

### Improved error messages

2rm will display improved error messages as to why a directory or file could
not be deleted.

#### Example

```sh
$ rm file.txt
> 2rm: Permission denied: file.txt. Try 'sudo rm file.txt' or 'chmod +w file.txt'.
```

### Errors if the program is being used by another program

Throws an error if the file is in-use by another program.
This protection can be bypassed with the `-f` or `--force` flags.

### Parallel deletion

2rm will delete files in parallel (if possible)

### Dependency awareness

If you delete a symbol/linked library 2rm will warn that the deletion will break
your system, and provide a list of programs that depend on the `.so` file.

You can override this behavior with the `-f` or `--force` flags.

### Config-based deletion

```yml
# ~/.local/share/2rm/config.yml

# while the backup config option is already supported, we should support
# uploading backups to cloud locations such as s3
backups: s3://my-bucket/backups
```

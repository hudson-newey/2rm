# 2RM

"rm with guard rails"

A wrapper for the "rm" command with soft-deletes, config-based deletion, debug information, and saner defaults

## "GNU Like" command line arguments

- `-i` Interactively prompt before each deletion request
- `-I` Prompt if deleting more than the interactive threshold of files (default 3)
- `-r`, `-R`, `--recursive` Recursively delete a directory of files
- `-v`, `--verbose` Emit additional verbose information
- `-d`, `--dir` Only delete empty directories
- `--help` Display help information (without deletion)
- `--version` Display version information (without deletion)

## Additional command line arguments

- `--overwrite` Overwrite the disk location location with zeros
- `--hard` Do not soft-delete file
- `--soft` Soft delete a file and store a backup (default `/tmp/2rm`)
- `--silent` Do not print out additional information priduced by 2rm. This is useful for scripting situations
- `--dry-run` Perform a dry run and show all the files that would be deleted
- `--bypass-protected` Using this flag will allow you to delete a file protected by the 2rm config
- `--notify` Send a system notification once deletion is complete
- `--force` Bypass protections

## Unsupported command line arguments

- `--interactive[=WHEN]` Interactive with a custom threshold
- `--one-file-system` Do not allow cross-file-system deletes
- `-f`, `--force` Bypass protections

## Features

### Removes the ability to remove your root directory

I have done this so that you can't accidentally add a space and remove your root directory with a typo such as

```sh
$ rm -rf ./directory /
>
```

(yes I know that you have to use `--no-preserve-root` and I have removed that too)

### Delete directories without the `-r` flag

You no longer have to add the `-r` flag when deleting a directory

(although you still can if you want to)

### Soft-deletes by default

By default, the program will soft delete your files by adding a hard link to the file in the `/tmp/2rm` directory.

This means that the files underlying INode is not freed, and can be recovered from the `/tmp/2rm` directory if you deleted the wrong file by mistake.

By using the `/tmp` directory, the operating system will **automatically hard delete files upon restart**.

Sometimes you want to hard delete a file/directory every time that you run the `rm` command e.g. you probably want your `node_modules` hard deleted every time and never want to soft delete them.
In this case, you can modify your `~/.local/share/2rm/config.yml` file to always hard delete `node_modules`.

### Overwriting disk location with zeros

When deleting a file with the linux inbuilt `rm` command, the file is still available on disk.

Meaning that the file can still be recovered by any sufficiently technical user.

This can be problematic when dealing with sensitive files such as private keys that if leaked could lead to catastrophic consequences.

You can overwrite a files disk location (rendering it unrecoverable) by using the `--overwrite` flag.

2rm will still soft-delete the file by default, but the soft-deleted file will be completely filled with zeros.

I made the decision that overwritten files will still be soft deleted because it might be useful for timestamp logging/auditing purposes.
E.g. "when did I overwrite xyz"

If you want to fully delete a file from disk and the file system use both the `--overwrite` and `--hard` flags.

### Config-based deletion

You can specify what directories are soft-deleted anb hard-deleted by using the `~/.local/share/2rm/config.yml` file.

```yml
# ~/.local/share/2rm/config.yml

# defaults to /tmp/2rm/ if not specified
# in the config file
# any files that are soft deleted will be
# backed up in the `backups` directory
backups: /tmp/2rm/
# whenever files matching these paths are deleted
# the disk location will be overwritten with zeros
overwrite:
    # when deleting ssh keys, we always want to 
    # overwrite them with zeros to protect
    # against attackers recovering the production
    # ssh keys
    - ".ssh/*"
hard:
    - "node_modules/"
    - "target/"
    - ".angular/"
    - ".next/"
# always soft delete backup files, 
# regardless of it they are configured
# for a hard delete
soft:
    - "*.bak"
# do not allow deleting these files/directories
# without using the `--bypass-protected` flag this
# does not make the file protected at the system level
# through other tools, but it does protect against
# accidental deletion through 2rm
protected:
    - ".ssh/"
# when using the -I flag without any arguments, the user will be prompted
# for confirmation before deleting each file if the number of files is
# greater or equal to this threshold
# default is 3 files/directories
interactive: 10
```

### Error codes

#### Error code _1_

An error that would have been thrown by a traditional "rm" command such as
the GNU "rm" implementation.

#### Error code _2_

An error was thrown during 2rm functionality (e.g. deleting a protected file)

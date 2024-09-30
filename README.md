# 2RM

"rm with guard rails"

Wraps the rm command with a more secure, safer, and more private version

## Command line arguments

- `--hard` Do not soft-delete file
- `--soft` Soft delete a file. A backup will be stored in `/tmp/2rm`
- `--silent` Do not print out additional information priduced by 2rm. This is useful for scripting situations
- `--dry-run` Perform a dry run and show all the files that would be deleted
- `--bypass-protected` Using this flag will allow you to delete a file protected by the 2rm config

## Features

### Removes the ability to remove your root directory

I have done this so that you can't accidently add a space and remove your root directory with a typo such as

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

### Config-based deletion

You can specify what directories are soft-deleted anb hard-deleted by using the `~/.local/share/2rm/config.yml` file.

```yml
# ~/.local/share/2rm/config.yml

# defaults to /tmp/2rm/ if not specified
# in the config file
# any files that are soft deleted will be
# backed up in the `backups` directory
backups: /tmp/2rm/
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
# without using the `-f` or `--force` flags
# this does not make the file un-deletable
# through other tools, but it does protect
# against accidental deletion through 2rm
protected:
    - ".ssh/"
```

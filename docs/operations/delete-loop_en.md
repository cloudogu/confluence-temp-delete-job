# Command `tempdel delete-loop`

## Input parameters

This command accepts three basic input parameters:

1. Start directory - the directory that should be searched for files
1. File age - the maximum age in hours of the files that should not be deleted
1. Deletion run interval - the time between the search intervals in minutes


### Start directory

The start directory is a mandatory parameter. The start directory can be freely selected and is not set to a specific value. There is no check if `tempdel` is executed against system-relevant directories (**BETTER NOT,** because system damage may occur).

### File age

The `--age`/`-a` switch can be used to optionally specify the maximum age (counted in hours from `now`) that files can have without being deleted. Only a positive integer value is accepted. The default value is `12` hours.

### Deletion run interval

The `--age`/`-a` switch can be used to optionally specify the interval (counted in minutes) between each deletion execution. Only a positive integer value is accepted. The default value is `60` minutes.

## Manpage

```
NAME:
   tempdel delete-loop - Endless loop that recursively deletes files and directories according the given parameters

USAGE:
   tempdel delete-loop [command options] directory

DESCRIPTION:
   This command recursively walks the given start directory and deletes files older than the given `age`. Directories will only be deleted last and only if there are no files left to be contained. The loop will run eternally until it receives the following signals: SIGHUP, SIGINT (Strg+C), SIGTERM, SIGKILL.

OPTIONS:
   --age value, -a value       Sets the max. age of files and directories in hours that will be deleted. Must be larger than zero. (default: 12)
   --interval value, -i value  Sets the interval in minutes to run the deletion routine. Must be larger than zero. (default: 60)
   --help, -h                  show help (default: false)
```

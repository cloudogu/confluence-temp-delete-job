# develop confluence-temp-delete-job

This document describes developmen-centric aspects.

## Tools

Needed:

- Docker
- Go compiler (currently 1.15.x)

## Helpful scripts

For testing on a real file system, these functions and settings can help.

**Debug output:**

With a log-level of `debug` single file and directory views become transparent. However, depending on the amount of files and directories looked at, an obscene amount of data is generated:

```bash
tempdel --log-level debug delete-loop ...
```

**Time settings:**

For the sake of fast feedback cycles, it is worth specifying time settings more meaningfully than in a production environment. The switch `-i/--interval` can be reduced down to one minute. Lower test times are only possible in unit tests.

```bash
tempdel delete-loop -i 1 ...
```

**Test files with specific timestamps:**

Test files with certain timestamps help to show the delete behavior in the real filesystem:

```bash
# create directory
mkdir -p /tmp/conftemp
# create empty files with timestamp 2021-03-03 03:03
touch deleteMe -t 2021030303
touch deleteMe2 -t 2021030303
# create empty files with current timestamp
touch leaveMe
touch leaveMe2

tempdel delete-loop ...
```

## Architecture

`tempdel` consists of two main parts:

1. periodic executions
1. actual deletion of files and directories

### Job cycles

`tempdel` is to be executed in a Confluence container as a background job. Therefore, a cron mechanism is out of the question, since cron works as a system service and can, among other things, keep the container alive on termination. Also, outputting tempdel output to the Docker output stream would be rather inelegant.

Therefore, the `delete-loop` command generally does not terminate. Exceptions are `panics` and Unix system signals, which are intercepted and handled in a dedicated manner:
- `SIGINT` (Ctrl+C is pressed during execution)
- `SIGHUP`
- `SIGTERM`
- (`SIGKILL` cannot be intercepted by the program, because it terminates the whole process)

The cycle is enabled by [`time.Ticker`](https://golang.org/pkg/time/#Ticker). The spacing of the individual intervals is specified as minutes on the CLI side. However, internally this is calculated back and forth between seconds and minutes to allow for fast unit tests.

### Deletion in two phases

The deletion routine of `tempdel` is essentially based on [`filepath.Walk`](https://golang.org/pkg/path/filepath/#Walk). In it, a file tree is walked recursively (and alphabetically, to act deterministically).

The nature of Confluence temp files is still unclear. Therefore, a single delete operation consists of two phases:

1. recursively delete files older than desired.
1. recursively delete directories that are left empty.

The two-phase approach has the advantage of deleting old directories that may contain new files. A complicating factor is that a file deletion updates the file stamp of a directory. Therefore the 1st phase is limited only to files. In the 2nd phase then exclusively empty directories are deleted, since on their timestamp no more reliance is anyway because of the file delete update.

Translated with www.DeepL.com/Translator (free version)
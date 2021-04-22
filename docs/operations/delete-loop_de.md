# Kommando `tempdel delete-loop`

## Eingabeparameter

Dieses Kommando akzeptiert drei grundlegende Eingabeparameter:

1. Startverzeichnis - das Verzeichnis, das nach Dateien durchsucht werden soll
1. Dateialter - das maximale Alter in Stunden der Dateien, die nicht gelöscht werden sollen
1. Löschlaufintervall - die Zeit zwischen den Suchintervallen in Minuten

### Startverzeichnis

Das Startverzeichnis ist ein Pflichtparameter. Das Startverzeichnis kann freigewählt werden und ist nicht auf einen bestimmten Wert festgelegt. Es findet keine Überprüfung statt, ob `tempdel` gegenüber systemrelevante Verzeichnissen ausgeführt wird (**BESSER NICHT,** da Systemschaden auftreten kann). 

### Dateialter

Mit dem Schalter `--age`/`-a` lässt sich optional bestimmen, wie alt (in Stunden gezählt von `jetzt`) Dateien maximal sein können, ohne gelöscht werden. Es wird nur ein positiver Ganzzahlwert akzeptiert. Standardwert ist `12` Stunden.

### Löschlaufintervall

Mit dem Schalter `--age`/`-a` lässt sich optional bestimmen, welcher Abstand (in Minuten gezählt) zwischen den einzelnen Löschausführungen liegen soll. Es wird nur ein positiver Ganzzahlwert akzeptiert. Standardwert ist `60` Minuten.

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

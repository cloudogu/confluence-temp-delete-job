# confluence-temp-delete-job entwickeln

Dieses Dokument beschreibt Aspekte, die die Entwicklung betreffen.

## Werkzeuge

Benötigt wird:

- Docker
- Go-Compiler (aktuell 1.15.x)

## Hilfreiche Skripte

Für Tests auf einem echten Dateisystem können diese Funktionen und Einstellungen helfen.

**Debug-Ausgaben:**

Mit einem Log-level von `debug` werden einzelne Datei- und Verzeichnisbetrachtungen transparent. Abhängig von der Menge der betrachteten Dateien und Verzreichnisse wird allerdings eine obszöne Datenmenge erzeugt:

```bash
tempdel --log-level debug delete-loop ...
```

**Einstellungen von Zeiten:**

Um schneller Feedbackzyklen willen lohnt es sich, die Zeiteinstellungen aussagekräftiger als in einer Produktionsumgebung zu bestimmen. Der Schalter `-i/--interval` kann bis auf eine Minute reduziert werden. Geringere Testzeiten sind nur in Unit-Tests möglich.

```bash
tempdel delete-loop -i 1 ...
```

**Testdateien mit bestimmten Zeitstempeln:**

Testdateien mit bestimmten Zeitstempel können helfen, das Löschverhalten im echten Dateisystem darzustellen:

```bash
# Verzeichnis anlegen
mkdir -p /tmp/conftemp
# leere Dateien mit Zeitstempel 2021-03-03 03:03 anlegen
touch /tmp/conftemp/deleteMe -t 202103030303
touch /tmp/conftemp/deleteMe2 -t 202103030303
# leere Dateien mit aktuellem Zeitstempel anlegen
touch /tmp/conftemp/leaveMe
touch /tmp/conftemp/leaveMe2

tempdel delete-loop ...
```

## Architektur

`tempdel` besteht aus zwei wesentlichen Teilen:

1. periodische Ausführungen
1. tatsächliche Löschung von Dateien und Verzeichnissen

### Job-Zyklen

`tempdel` soll in einem Confluence-Container als Hintergrundjob ausgeführt werden. Daher kommt ein Cron-Mechanismus nicht infrage, da Cron als Systemdienst arbeitet und u. Um. den Container bei Beendigung am Leben halten kann. Auch die Ausgabe von tempdel-Ausgaben in den Docker-Ausgabestrom wäre eher unelegant.

Daher beendet sich das Command `delete-loop` grundsätzlich nicht. Ausnahmen bilden `Panics` und Unix-Systemsignale, die dediziert abgefangen und behandelt werden:
- `SIGINT` (Strg+C wird während der Ausführung gedrückt)
- `SIGHUP`
- `SIGTERM`
- (`SIGKILL` kann Programmseitig nicht abgefangen werden, da es den gesamten Prozess beendet)

Der Zyklus wird durch [`time.Ticker`](https://golang.org/pkg/time/#Ticker) ermöglicht. Der Abstand der einzelnen Intervalle wird CLI-seitig als Minuten angegeben. Intern wird dies jedoch zwischen Sekunden und Minuten hin- und hergerechnet, um schnelle Unit-Tests zu ermöglichen.

### Löschung in zwei Phasen

Die Löschroutine von `tempdel` beruht im wesentlich auf [`filepath.Walk`](https://golang.org/pkg/path/filepath/#Walk). Darin wird rekursiv (und alphabetisch, um deterministisch zu agieren) ein Dateibaum abgelaufen.

Die Natur von Confluence-Temp-Dateien ist noch ungeklärt. Daher besteht ein einzelner Löschvorgang aus zwei Phasen:

1. rekursiv Dateien löschen, die älter als gewünscht sind
1. rekursiv Verzeichnisse löschen, die leer übrig bleiben

Das Vorgehen in zwei Phasen hat den Vorteil, dass alte Verzeichnisse gelöscht werden, die evtl. neue Dateien enthalten. Erschwerend kommt hinzu, dass eine Dateilöschung den Dateistempel eines Verzeichnisses aktualisiert. Daher beschränkt sich die 1. Phase lediglich auf Dateien. In der 2. Phase werden dann ausschließlich leere Verzeichnisse gelöscht, da auf deren Zeitstempel wegen der Dateilöschaktualisierung ohnehin kein Verlass mehr ist.
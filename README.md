# Duphunter

## Usage

Vanilla call -
```bash
$ duphunter .
# Output -
/path/to/f1
/path/to/f1 -- /path/to/f1copy1
/path/to/f1 -- /path/to/f1copy2
/path/to/f2
/path/to/f2 -- /path/to/f2copy1
/path/to/f2 -- /path/to/f2copy2
```

Change output template -
```bash
$ duphunter . --basetmpl '' --outtmpl '{"source": "$1", "copy": "$2"}'
# Output -
{"source": "/path/to/f1" "copy": "/path/to/f1copy1"}
{"source": "/path/to/f1" "copy": "/path/to/f1copy2"}
{"source": "/path/to/f2" "copy": "/path/to/f2copy1"}
{"source": "/path/to/f2" "copy": "/path/to/f2copy2"}
```

Change output template to emulate a command. Note that the "commands" are not actually run, output template is always printed.
```bash
$ duphunter . --basetmpl '' --outtmpl 'cp -l $1 $2'
# Output -
cp -l /path/to/f1 /path/to/f1copy1
cp -l /path/to/f1 /path/to/f1copy2
cp -l /path/to/f2 /path/to/f2copy1
cp -l /path/to/f2 /path/to/f2copy2
```

## Help

Run `duphunter --help` to display available args.

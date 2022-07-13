# Duphunter

A command-line duplicate file finder with customizable output.

## Usage

### Call Examples
```bash
$ duphunter .
# Output -
/path/to/f1 -- /path/to/f1copy1
/path/to/f1 -- /path/to/f1copy2
/path/to/f2 -- /path/to/f2copy1
/path/to/f2 -- /path/to/f2copy2
```

Change output template to output json.
```bash
$ duphunter . --outtmpl '{"source": "$0", "copy": "$1"}'
# Output -
{"source": "/path/to/f1" "copy": "/path/to/f1copy1"}
{"source": "/path/to/f1" "copy": "/path/to/f1copy2"}
{"source": "/path/to/f2" "copy": "/path/to/f2copy1"}
{"source": "/path/to/f2" "copy": "/path/to/f2copy2"}
```

### Code Generation & Cleaning Up
Change output template to print commands to create hard links. Note that the
commands are not actually executed. Output is *always* just printed.

```bash
$ duphunter . --outtmpl 'cp -l $0 $1'
# Output -
cp -l /path/to/f1 /path/to/f1copy1
cp -l /path/to/f1 /path/to/f1copy2
cp -l /path/to/f2 /path/to/f2copy1
cp -l /path/to/f2 /path/to/f2copy2
```
If desired, you can copy the output to a bash script and execute it.

```bash
# WARNING: If you follow this example, existing attributes of replaced
# files such timestamps will be lost.

# Replace all duplicates with hard links.
duphunter . --outtmpl 'cp -l $0 $1' > cleanup.sh

# Review the commands carefully.
head cleanup.sh

# Run generated commands. After you do this, there is no going back!
source cleanup.sh
```

## Help

Run `duphunter --help` to display available args.

## Installation

```bash
git clone ...
cd ...
scripts/build_and_install.sh
```

# Duphunter

A command-line duplicate file finder with customizable output.

## Usage

### Call Examples
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

Change output template to output json.
```bash
$ duphunter . --basetmpl '' --outtmpl '{"source": "$1", "copy": "$2"}'
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
$ duphunter . --basetmpl '' --outtmpl 'cp -l $1 $2'
# Output -
cp -l /path/to/f1 /path/to/f1copy1
cp -l /path/to/f1 /path/to/f1copy2
cp -l /path/to/f2 /path/to/f2copy1
cp -l /path/to/f2 /path/to/f2copy2
```
If desired, you can copy the output to a bash script and execute it.

WARNING: Do it at your own risk!
```bash
# WARNING: If you follow this example, existing attributes of replaced
# files such timestamps will be lost.

# Replace all duplicates with hard links.
duphunter . --basetmpl '' --outtmpl 'cp -l $1 $2' > cleanup.sh

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
./build_and_install.sh
```
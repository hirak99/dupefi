A command-line duplicate file finder with customizable output.

# Installation (Linux)

Dependency: Install [golang](https://go.dev/doc/install) for your system.

Then run the following -

```bash
git clone https://github.com/hirak99/duphunter
cd duphunter
scripts/build_and_install.sh
```

Note: With minor modifications, it *should* also work on Windows. If anyone does
it, please send me a pull request.

# Usage

## Call Examples
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

# Advanced Usage
## Hard Link All Duplicates
To do this we will modify the template to generate code, and then run it.

Note that the commands are not actually executed. Output is *always* just
printed.

```bash
$ duphunter . --outtmpl 'cp -l $0 $1'
# Output -
cp -l /path/to/f1 /path/to/f1copy1
cp -l /path/to/f1 /path/to/f1copy2
cp -l /path/to/f2 /path/to/f2copy1
cp -l /path/to/f2 /path/to/f2copy2
```
If desired, you can copy the output to a bash script and execute it. Below is an
example that hard links all duplicate files.

```bash
# WARNING: If you follow this example, existing attributes of replaced
# files such timestamps will be lost.

# Replace all duplicates with hard links.
duphunter . --outtmpl 'cp -l $0 $1' > cleanup.sh

# Review the commands carefully.
head cleanup.sh

# Run generated commands.
# WARNING: Running this will make changes that you cannot undo!
source cleanup.sh
```

In similar way, you can create and run commands to delete, archive, symlink, query all duplicates too.

## Estimate Size of All Duplicates

Print the total size used in duplicates.

```bash
duphunter . --outtmpl '$1' | xargs du -sch
```

# Help

Run `duphunter --help` to display available args.

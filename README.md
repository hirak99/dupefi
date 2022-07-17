A command-line duplicate file finder designed with linux philosophy.

It allows customizable outputs, and thus chaining of commands (see Advanced
Usage below).

# Installation

## Linux

Dependency: Install [golang](https://go.dev/doc/install) for your system.

Then run the following -

```bash
git clone https://github.com/hirak99/dupefi
cd dupefi
scripts/build_and_install.sh
```

## Windows
With minor modifications, it *should* also work on Windows. If anyone does it,
please send me a pull request.

# Usage

## Call Examples
```bash
$ dupefi .
# Output -
"/path/to/f1" -- "/path/to/f1copy"
"/path/to/f2" -- "/path/to/f2copy"
```

List original (base) files followed by all duplicates -
```bash
$ dupefi . --basetmpl '"$1"' --outtmpl '  "$1"'
# Output -
"/path/to/f1"
  "/path/to/f1copy"
"/path/to/f2"
  "/path/to/f2copy"
```

Just list all duplicate files (not the originals) -
```bash
$ dupefi . --outtmpl '"$1"'
# Output -
"/path/to/f1copy"
"/path/to/f2copy"
```

Change output template to output json.
```bash
$ dupefi . --outtmpl '{"source": "$0", "copy": "$1"}'
# Output -
{"source": "/path/to/f1" "copy": "/path/to/f1copy"}
{"source": "/path/to/f2" "copy": "/path/to/f2copy"}
```

# Advanced Usage (Linux)
By itself, dupefi doesn't provide options for acting on the duplicate files -
e.g. to delete them. It is designed with the Linux philosophy in mind, so that
you can use it with other commands to do that and more.

Some examples are shown below.

## Hard Link All Duplicates

One liner -
> **Warning**
> This command will run immediately, and

```bash
# WARNING: This will run immediately and will be irreversible.
source <(dupefi . --outtmpl 'cp -lf "$0" "$1"')
```

### Explanation
The `dupefi` program doesn't run anything by itself.

So we change the output template so that what it prints are the bash commands
for hard linking -

```bash
$ dupefi . --outtmpl 'cp -l "$0" "$1"'
# Output -
cp -l "/path/to/f1" "/path/to/f1copy"
cp -l "/path/to/f2" "/path/to/f2copy"
```

Then we run it in the shell by sourcing it.

## Estimate Size of All Duplicates

Print the total size lost in duplicates -

```bash
dupefi . --outtmpl '$1' | xargs du -sch
```

## Delete Duplicates (Dangerous!)

> **Warning**
> This command will run immediately, and it will remove files irreversibly.

```bash
# WARNING: This will immediately, irreversibly, delete the duplicates.
source <(dupefi . --outtmpl 'rm "$1"')
```

For an explanation, see the example on hardlinking.

# Help

`dupefi [OPTIONS] Directories...`

Positional args -
* Directories... One or more directories to scan.

Options -

* **--minsize=**: Minimum file size to consider in bytes. The default value 1
  ignores all empty files.
* **--outtmpl=**: Template of output that's generated per each copy. Any `$0`
  will be replaced with base file, and `$1` with the copy.
* **--basetmpl=**: Template of output for the base file. Generated once per group
  of duplicate files. Any `$1` is replaced with the file name.
* **--regex=**: Regular expression for files to scan within directories
  specified. E.g. `'\.jpg$'` to consider only files with .jpg extension.
* **--regexnodup=**: After the duplicate scan, any files matching these will not
  be reported as duplicates. But they may still be reported as the original
  (base) file, i.e. `$0` in the --outtmpl arg.

Short options -
* **-c**: Use checksum to compare. In this mode, any files with equal sha256
  will be considered equal. The checksum will be computed only if there are two
  files with similar size. This can speed up comparison if there are multiple
  copies of large files, but can slow down comparison for slow CPU's.
* **-i**: Report a file as duplicate, even if it has the same inode as the base
  (i.e. is a hardlink) and doesn't take up additional space.
* **-v**: Verbose.

Please run `dupefi --help` to display all args.

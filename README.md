# vils

Like ls, list provided files or the files they contain if they are
directories, but send the list to your editor where you can manipulate the
list before files are renamed/removed according to your changes.

## Manifest format

Manifest lines start with an identifying number and a tab, followed by the
filename. To remove a file, simply remove the line. To rename a file, change
the file name without changing the identifying number.

## Caveats

Changes are applied in the order of the original manifest, and there can only
be one operation per file per manifest.  Therefore, you can't use `vils` to
e.g. swap filenames within the same manifest.

## Editor

`vils` prefers using your VISUAL editor, but will fall back to EDITOR if that
is empty or not set and finally just `vi` if neither is set.

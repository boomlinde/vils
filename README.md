# vils

Like ls, list provided files or the files they contain if they are
directories, but send the list to your editor where you can manipulate the
list before files are renamed/removed according to your changes.

## Manifest format

Manifest lines start with an identifying number and a tab, followed by the
filename. To remove a file, simply remove the line. To rename a file, change
the file name without changing the identifying number.

## Caveats

Each changed file in the manifest will be renamed, adding the suffix .vils.tmp
before the operation is applied. This is to ensure that operations don't
interfere with each other, for example so that you can swap filenames.

## Editor

`vils` prefers using your $VISUAL editor, but will fall back to $EDITOR if
that is empty or unset and finally just `vi` if neither is set.

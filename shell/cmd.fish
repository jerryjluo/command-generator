function cmd-generate --description "Generate a command with AI"
    set -l tmpfile (mktemp /tmp/cmd-output.XXXXXX)
    command stty sane </dev/tty 2>/dev/null
    command cmd --output $tmpfile </dev/tty
    if test $status -eq 0 -a -s $tmpfile
        commandline -r (cat $tmpfile)
    end
    rm -f $tmpfile
    commandline -f repaint
end

bind \cg cmd-generate

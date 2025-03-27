function prompt(term) {
    command = ""
    term.write("\r\n$ ")
}

function termInit() {
    var term = new Terminal({
        cursorBlink: true,
        fontFamily: "monospace",
    });
    term.open(document.getElementById('Commandline'));
    term.prompt = function() {
        term.write("\r\n$ ")
    }

    prompt(term)
    term.onData(function(evt) {
        switch (evt) {
            case '\u0003': // Ctrl+C
                term.write('^C');
                prompt(term);
                break;
            case '\r': // Enter
                // Post task
                prompt(term)
                command = '';
                break;
            case '\u007F': // Backspace (DEL)
                if (term._core.buffer.x > 2) {
                    term.write('\b \b');
                    if (command.length > 0) {
                        command = command.slice(0, command.length - 1);
                    }
                }
                break;
            default:
                if (evt >= String.fromCharCode(0x20) && evt <= String.fromCharCode(0x7E) || evt >= '\u00a0') {
                    command += evt;
                    term.write(evt);
                }
        }
    })
}

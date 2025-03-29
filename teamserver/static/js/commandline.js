var term = new Terminal({
    cursorBlink: true,
    fontFamily: "monospace",
});

dispose = term.onData()

function prompt(agentId) {
    command = ""
    term.write("\r\nAgent " + agentId + " $ ")
}

async function termInit(agentId) {
    dispose.dispose()
    term.clear()
    term.open(document.getElementById('Commandline'));
    const ws = new WebSocket(`/api/agents/${agentId}/socket`)
    await awaitSocketOpen(ws)

    prompt(agentId)
    dispose = term.onData(async function(evt) {
        switch (evt) {
            case '\u0003': // Ctrl+C
                term.write('^C');
                prompt(agentId);
                break;
            case '\r': // Enter
                term.write('\r\n')
                term.writeln(await runCommand(ws, command))
                prompt(agentId)
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

function awaitSocketOpen(ws) {
    return new Promise(function(resolve) {
        const listener = function() {
            ws.removeEventListener("open", listener)
            resolve()
        }
        ws.addEventListener("open", listener)
    })
}

function recieveOutput(ws) {
    return new Promise(function(resolve, reject) {
        const messageListener = function(message) {
            ws.removeEventListener("message", messageListener)
            resolve(message.data)

        }

        const errorListener = function(error) {
            ws.removeEventListener("error", errorListener)
            reject(error)
        }

        ws.addEventListener("message", messageListener)
        ws.addEventListener("error", errorListener)
    })
}

async function runCommand(ws, command) {
    ws.send(JSON.stringify({ Task: command }))
    const taskResultJson = await recieveOutput(ws)
    return JSON.parse(taskResultJson).Output
}

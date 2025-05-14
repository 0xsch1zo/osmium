import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';

var term = new Terminal({
    cursorBlink: true,
    fontFamily: "monospace",
    theme: {
        background: "#00000000", // transparent
    },
});
const fitAddon = new FitAddon();
term.loadAddon(fitAddon)

var onDataDispose = term.onData()
var command = ""
var promptLength = 0

function prompt(agentId) {
    command = ""
    const prompt = `Agent ${agentId} >> `
    promptLength = prompt.length
    return prompt
}
function promptNewLine(agentId) {
    term.write("\r\n" + prompt(agentId))
}

function promptNoNewLine(agentId) {
    term.write(prompt(agentId))
}

async function termInit(agentId) {
    var placeholder = document.getElementById("commandline-placeholder")
    if (placeholder != null) {
        placeholder.remove()
    }
    onDataDispose.dispose()
    term.reset()
    term.open(document.getElementById('commandline'))
    fitAddon.fit()
    const ws = new WebSocket(`/api/agents/${agentId}/socket`)
    await awaitSocketOpen(ws)

    await htmx.ajax('GET', `/api/agents/${agentId}/results`, '#task-results-body')
    let el = document.getElementById("task-results-body")
    el.addEventListener('htmx:oobBeforeSwap', function() {
        let taskResultsPlaceholder = document.getElementById('task-results-placeholder')
        if (taskResultsPlaceholder != null) {
            taskResultsPlaceholder.remove()
        }
    })

    promptNewLine(agentId)
    onDataDispose = term.onData(async (evt) => {
        switch (evt) {
            case '\u0003': // Ctrl+C
                term.write('^C');
                promptNewLine(agentId);
                break;
            case '\r': // Enter
                await handleCommand(ws, agentId, term)
                break;
            case '\u007F': // Backspace
                backspace(term)
                break;
            case '\x1b[C': // cursor right
                moveCursor('\x1b[C', term)
                break
            case '\x1b[D': // cursor left
                moveCursor('\x1b[D', term)
                break
            default:
                handleCharacter(evt, term)
        }
    })
}

window.termInit = termInit

async function handleCommand(ws, agentId, term) {
    term.write('\r\n')
    if (command) {
        term.writeln(formatOutput(await runCommand(ws, command)))
    }
    promptNoNewLine(agentId)
    command = '';
}

function handleCharacter(evt, term) {
    if (evt >= String.fromCharCode(0x20) && evt <= String.fromCharCode(0x7E) || evt >= '\u00a0') {
        if (term.buffer.normal.cursorX != promptLength + command.length) {
            const commandIndex = term.buffer.normal.cursorX - promptLength
            term.write(evt + command.slice(commandIndex))
            command = command.slice(0, commandIndex) + evt + command.slice(commandIndex)
            term.write(`\x1b[${command.length - commandIndex - 1}D`)
        } else {
            command += evt;
            term.write(evt);
        }
    }
}

function backspace(term) {
    if (term._core.buffer.x <= promptLength)
        return

    if (term.buffer.normal.cursorX != promptLength + command.length) {
        const originalPosition = term.buffer.normal.cursorX
        const commandIndex = originalPosition - promptLength
        term.write('\b')
        term.write(command.slice(commandIndex) + ' ') // overwrite the last character
        command = command.slice(0, commandIndex - 1) + command.slice(commandIndex)
        const onCursorMoveDispose = term.onCursorMove(() => {
            term.write(`\x1b[${term.buffer.normal.cursorX - commandIndex - promptLength + 1}D`)
            onCursorMoveDispose.dispose()
        })
    } else {
        term.write('\b \b');
        if (command.length > 0) {
            command = command.slice(0, command.length - 1);
        }
    }
}

function moveCursor(charcode, term) {
    if (term.buffer.normal.cursorX > promptLength && term.buffer.normal.cursorX < promptLength + command.length) {
        term.write(charcode)
    } else if (term.buffer.normal.cursorX == promptLength && charcode == '\x1b[C') {
        term.write(charcode)
    } else if (term.buffer.normal.cursorX == promptLength + command.length && charcode == '\x1b[D') {
        term.write(charcode)
    }
}

function formatOutput(output) {
    const length = output.length
    for (let i = 0; i < length; ++i) {
        if (output[i] === '\n') {
            output = output.slice(0, i) + '\r' + output.slice(i)
            i++ // skip return just added
        }
    }
    return output
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

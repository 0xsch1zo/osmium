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
    onDataDispose = term.onData(async function(evt) {
        switch (evt) {
            case '\u0003': // Ctrl+C
                term.write('^C');
                promptNewLine(agentId);
                break;
            case '\r': // Enter
                term.write('\r\n')
                if (command) {
                    term.writeln(formatOutput(await runCommand(ws, command)))
                }
                promptNoNewLine(agentId)
                command = '';
                break;
            case '\u007F': // Backspace (DEL)
                if (term._core.buffer.x > promptLength) {
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

window.termInit = termInit

function formatOutput(output) {
    const length = output.length
    for (let i = 0; i < length; ++i) {
        console.log(i)
        console.log(output.length)
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

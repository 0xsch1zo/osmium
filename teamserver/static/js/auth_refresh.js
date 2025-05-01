async function authRefresh() {
    const resp = await fetch("/api/auth/refreshTime")
    const json = await resp.json()
    const refTime = Date.parse(json.RefTime)
    const waitDuration = refTime - Date.now()
    setTimeout(async function() {
        await fetch("/api/auth/refresh", { method: "POST" })
        authRefresh()
    }, waitDuration)
}

authRefresh()

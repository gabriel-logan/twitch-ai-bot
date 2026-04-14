const resultEl = document.getElementById("result");

resultEl.style.display = "none";

async function setEnvironment() {
    const broadcasterId = document.getElementById("broadcasterId").value.trim();
    const botUserId = document.getElementById("botUserId").value.trim();
    const botUserName = document.getElementById("botUserName").value.trim();
    const keyword = document.getElementById("keyword").value.trim();
    const chatMessageMaxLength = document.getElementById("chatMessageMaxLength").value.trim();
    const groqMaxContextInput = document.getElementById("groqMaxContextInput").value.trim();
    const statusEl = document.getElementById("envStatus");

    const params = new URLSearchParams();

    if (broadcasterId) params.append("twitch_broadcaster_id", broadcasterId);
    if (botUserId) params.append("twitch_bot_user_id", botUserId);
    if (botUserName) params.append("twitch_bot_user_name", botUserName);
    if (keyword) params.append("twitch_key_word_to_call", keyword);
    if (chatMessageMaxLength) params.append("twitch_chat_message_max_length", chatMessageMaxLength);
    if (groqMaxContextInput) params.append("groq_max_context_input", groqMaxContextInput);

    if ([...params.keys()].length === 0) {
        alert("Fill at least one field");
        return;
    }

    try {
        statusEl.style.display = "block";
        statusEl.textContent = "Saving...";
        statusEl.className = "result mt-4";

        const res = await fetch(`/api/twitch/set-environment?${params.toString()}`, {
            method: "POST"
        });

        const data = await res.json();

        if (!res.ok) {
            throw new Error(data);
        }

        statusEl.textContent = "Environment updated successfully!";
        statusEl.className = "result mt-4 success";

        setTimeout(() => {
            location.reload();
        }, 600);
    } catch (err) {
        console.error(err);

        if (typeof err === "string") {
            statusEl.textContent = err;
        } else if (Array.isArray(err)) {
            statusEl.textContent = err.join(", ");
        } else if (err instanceof Error) {
            statusEl.textContent = err.message;
        } else {
            statusEl.textContent = "Failed to update environment";
        }

        statusEl.className = "result mt-4 error";
    }
}

async function startBot() {
    const btn = document.getElementById("botToggleBtn");
    btn.disabled = true;
    btn.innerHTML = '<span class="spinner"></span> Starting...';

    try {
        const response = await fetch("/api/twitch/start-bot", { method: "POST" });

        if (response.ok) {
            setTimeout(() => {
                location.reload();
            }, 1500);
        } else {
            btn.disabled = false;
            btn.innerHTML = 'Start Bot';
        }
    } catch (err) {
        console.error(err);
        btn.disabled = false;
        btn.innerHTML = 'Start Bot';
    }
}

async function stopBot() {
    const btn = document.getElementById("botToggleBtn");
    btn.disabled = true;
    btn.innerHTML = '<span class="spinner"></span> Stopping...';

    try {
        const response = await fetch("/api/twitch/stop-bot", { method: "POST" });

        if (response.ok) {
            setTimeout(() => {
                location.reload();
            }, 1500);
        } else {
            btn.disabled = false;
            btn.innerHTML = 'Stop Bot';
        }
    } catch (err) {
        console.error(err);
        btn.disabled = false;
        btn.innerHTML = 'Stop Bot';
    }
}

document.getElementById("botToggleBtn").addEventListener("click", function () {
    const isRunning = this.textContent.includes("Stop");
    if (isRunning) {
        stopBot();
    } else {
        startBot();
    }
});

async function fetchUser() {
    const inputEl = document.getElementById("loginInput");
    const input = inputEl.value.trim();

    if (!input) {
        alert("Please enter a login");
        return;
    }

    const logins = input
        .split(",")
        .map(l => l.trim())
        .filter(Boolean);

    const params = new URLSearchParams();

    for (const login of logins) {
        params.append("login", login);
    }

    try {
        resultEl.textContent = "Loading...";
        resultEl.style.display = "block";
        resultEl.className = "result mt-4 user-lookup-result";

        const res = await fetch(`/api/twitch/user-info?${params.toString()}`);

        if (!res.ok) {
            throw new Error(`HTTP ${res.status}`);
        }

        const data = await res.json();
        resultEl.textContent = JSON.stringify(data, null, 2);
        resultEl.className = "result mt-4 user-lookup-result";
    } catch (err) {
        console.error(err);
        resultEl.textContent = "Failed to fetch user(s)";
        resultEl.className = "result mt-4 error user-lookup-result";
    }
}

document.getElementById("loginInput").addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
        fetchUser();
    }
});

const BASE = "/pattern/api";
let tasks = [],
  filter = "all",
  callCount = 0,
  logEntries = [];
function addLog(msg) {
  callCount++;
  document.getElementById("callCount").textContent = callCount;
  const t = new Date().toLocaleTimeString("en", {
    hour12: false,
  });
  logEntries.unshift(`[${t}]  GetInstance() → ${msg}`);
  if (logEntries.length > 30) logEntries.pop();
  document.getElementById("activityLog").innerHTML = logEntries
    .map((e) => `<div class="log-entry">${e}</div>`)
    .join("");
}
async function loadTasks() {
  addLog("SELECT * FROM tasks");
  const r = await fetch(`${BASE}/tasks`);
  tasks = (await r.json()) || [];
  render();
  loadStats();
}
async function addTask() {
  const title = document.getElementById("taskInput").value.trim();
  if (!title) return;
  const priority = document.getElementById("prioritySelect").value;
  addLog("INSERT INTO tasks");
  await fetch(`${BASE}/tasks`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title, priority }),
  });
  document.getElementById("taskInput").value = "";
  loadTasks();
}
async function toggleTask(id) {
  addLog(`UPDATE tasks … id=${id}`);
  await fetch(`${BASE}/tasks/toggle?id=${id}`);
  loadTasks();
}
async function deleteTask(id) {
  addLog(`DELETE FROM tasks … id=${id}`);
  await fetch(`${BASE}/tasks/delete?id=${id}`);
  loadTasks();
}
async function loadStats() {
  const r = await fetch(`${BASE}/stats`);
  const s = await r.json();
  document.getElementById("instanceId").textContent = s.instance_id || 1;
  document.getElementById("statsContent").innerHTML = `
    <div class="ok-banner" style="margin-bottom:8px">
      <p style="color:#2e5c8a;font-family:'IBM Plex Mono',monospace;font-size:11px;line-height:1.7">${s.message}</p>
    </div>
    <div style="display:flex;justify-content:space-between;color:#aaa;margin-top:6px">
      <span>Вызовы GetInstance()</span><span style="color:#2e5c8a">${s.get_instance_calls}</span>
    </div>
    <div style="display:flex;justify-content:space-between;color:#aaa;margin-top:4px">
      <span>Открыть соединения с базой данных</span><span style="color:#4a7c59">${s.open_connections}</span>
    </div>`;
}
function setFilter(f, btn) {
  filter = f;
  document
    .querySelectorAll(".tab")
    .forEach((b) => b.classList.remove("active"));
  btn.classList.add("active");
  render();
}
function render() {
  const list = tasks.filter((t) =>
    filter === "active" ? !t.done : filter === "done" ? t.done : true,
  );
  const el = document.getElementById("taskList");
  if (!list.length) {
    el.innerHTML = `<div class="text-sm text-center py-6" style="color:#bbb">Нет заданий</div>`;
    return;
  }
  el.innerHTML = list
    .map(
      (t) =>
        `<div class="task-row ${t.done ? "done" : ""}">
          <div class="p-dot p-${t.priority}"></div>
          <div class="cb ${t.done ? "checked" : ""}" onclick="toggleTask(${t.id})"></div>
          <span class="task-title flex-1 text-sm">${esc(t.title)}</span>
          <span class="mono text-xs" style="color:#bbb">${t.priority}</span>
          <button class="del-btn" onclick="deleteTask(${t.id})">✕</button>
        </div>`,
    )
    .join("");
}
function esc(s) {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}
document.getElementById("taskInput").addEventListener("keydown", (e) => {
  if (e.key === "Enter") addTask();
});
loadTasks();
setInterval(loadStats, 8000);

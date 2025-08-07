const KEY = "selectedTeamIds"; // localStorage key

function getSelectedIds() {
  try { return JSON.parse(localStorage.getItem(KEY)) || []; }
  catch { return []; }
}
function setSelectedIds(ids) {
  const unique = Array.from(new Set(ids.map(Number).filter(Boolean)));
  localStorage.setItem(KEY, JSON.stringify(unique));
  renderSelected(unique);
}
function renderSelected(ids) {
  const box = document.getElementById("selected");
  if (!ids.length) { box.innerHTML = "<i>No teams selected</i>"; return; }
  box.innerHTML = ids.map(id => `<span class="chip" data-id="${id}" title="Click to remove">${id}</span>`).join("");
  box.querySelectorAll(".chip").forEach(el => {
    el.onclick = () => {
      const id = Number(el.dataset.id);
      setSelectedIds(getSelectedIds().filter(x => x !== id));
    };
  });
}

async function search() {
  const game = document.getElementById("game").value;
  const q = document.getElementById("q").value.trim();
  const out = document.getElementById("results");
  if (!q) { out.textContent = "Enter a search term."; return; }

  out.textContent = "Searching…";
  const url = `/api/teams/search?game=${encodeURIComponent(game)}&q=${encodeURIComponent(q)}`;
  const res = await fetch(url);
  if (!res.ok) { out.textContent = "Search failed."; return; }
  const teams = await res.json();

  if (!teams.length) { out.textContent = "No results."; return; }
  out.innerHTML = teams.map(t => `
    <div class="result" data-id="${t.id}" title="Click to add">
      <img src="${t.image_url || ""}" alt="" width="24" style="vertical-align:middle;margin-right:6px;">
      ${t.name} ${t.acronym ? `(${t.acronym})` : ""} — <small>ID ${t.id}</small>
    </div>`).join("");

  out.querySelectorAll(".result").forEach(el => {
    el.onclick = () => {
      const id = Number(el.dataset.id);
      setSelectedIds([...getSelectedIds(), id]);
    };
  });
}

window.addEventListener("DOMContentLoaded", () => {
  renderSelected(getSelectedIds());
  document.getElementById("searchBtn").onclick = search;

  // optional: debounce typing -> search
  let t; document.getElementById("q").addEventListener("input", () => {
    clearTimeout(t); t = setTimeout(search, 300);
  });
});

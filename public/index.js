const KEY = "followedTeams.v1";
let followed = new Set(JSON.parse(localStorage.getItem(KEY) || "[]"));

const els = {
  region: document.getElementById("region"),
  game: document.getElementById("game"),
  onlyFollowed: document.getElementById("onlyFollowed"),
  refresh: document.getElementById("refresh"),
  upcomingList: document.getElementById("upcomingList"),
  liveList: document.getElementById("liveList"),
  teamsSummary: document.getElementById("teamsSummary"),
  teamsMenu: document.getElementById("teamsMenu"),
  teamsSearch: document.getElementById("teamsSearch"),
  teamsList: document.getElementById("teamsList"),
  btnSelAll: document.getElementById("teamsSelectAll"),
  btnClear: document.getElementById("teamsClear"),
  btnSave: document.getElementById("teamsSave"),
};

// TODO: wire to your Go API
async function fetchUpcoming(params) {
  // return (await fetch(`/api/upcoming?${new URLSearchParams(params)}`)).json();
  return []; // placeholder
}
async function fetchLive(params) {
  // return (await fetch(`/api/live?${new URLSearchParams(params)}`)).json();
  return []; // placeholder
}

function renderMatches(list, mount) {
  if (!list.length) { mount.innerHTML = '<div class="card">No matches.</div>'; return; }
  mount.innerHTML = list.map(m => `
    <div class="card">
      <div>
        <strong>${m.tournament}</strong><br/>
        ${m.teamA} vs ${m.teamB}
      </div>
      <div style="text-align:right">
        <div>${new Date(m.startsAt).toLocaleString()}</div>
        ${m.status === "live" ? '<span>🔴 Live</span>' : ""}
      </div>
    </div>
  `).join("");
}

async function reload() {
  const params = {
    region: els.region.value,
    game: els.game.value,
    teams: [...followed].join(","),
    onlyFollowed: els.onlyFollowed.checked ? "1" : "",
  };
  const [upcoming, live] = await Promise.all([
    fetchUpcoming(params),
    fetchLive(params),
  ]);

  // Optional client-side filter by followed
  const filterByFollowed = (arr) =>
    els.onlyFollowed.checked && followed.size
      ? arr.filter(m => followed.has(m.teamAId) || followed.has(m.teamBId))
      : arr;

  renderMatches(filterByFollowed(upcoming), els.upcomingList);
  renderMatches(filterByFollowed(live), els.liveList);
}
// How to handle team search and selection
els.teamsSearch.addEventListener("input", () => {
  const searchTerm = els.teamsSearch.value.toLowerCase();
  const filteredTeams = TEAMS.filter(team => team.name.toLowerCase().includes(searchTerm));
  els.teamsList.innerHTML = filteredTeams.map(team => `
    <li>
      <label>
        <input type="checkbox" value="${team.id}" ${followed.has(team.id) ? "checked" : ""}>
        ${team.name}
      </label>
    </li>
  `).join("");
}); 

// --- Team picker bits (reuse your earlier version) ---
const TEAMS = []; // fetch from /api/teams on load if you want
function setFollowed(ids) {
  followed = new Set(ids);
  localStorage.setItem(KEY, JSON.stringify([...followed]));
  els.teamsSummary.textContent = `Followed teams (${followed.size})`;
}
els.btnSave.addEventListener("click", () => (els.teamsMenu.open = false));
els.onlyFollowed.addEventListener("change", reload);
els.region.addEventListener("change", reload);
els.game.addEventListener("change", reload);
els.refresh.addEventListener("click", reload);

// Polling for live updates
setInterval(reload, 30_000); // 30s; tune as needed

// Init
setFollowed([...followed]);
reload();

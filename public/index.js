const el = (id) => document.getElementById(id);

/* ----Filter Toolbar ----- */
const regionEl = el("region");
const gameEl = el("game");
const onlyFollowedEl = el("onlyFollowed");
const refreshBtn = el("refreshBtn");

/* ----Teams Menu ----- */
const teamsMenu = el("teamsMenu");
const menu = el("menu");
const teamsSummary = el("teamsSummary");
const teamSearchEl = el("teamSearch");
const teamsList = el("teamsList");
const teamsSelectAllBtn = el("teamsSelectAllBtn")
const teamsClearBtn = el("teamsClearBtn")
const teamsSaveBtn = el("teamsSaveBtn")

/* ----Matches ----- */
const upcomingMatchList = el("upcomingMatchList");
const liveMatchList = el("liveMatchList"); // TODO

const KEY = "followedTeams.v1";
let followed = new Set(JSON.parse(localStorage.getItem(KEY) || "[]"));
function saveFollowed() {
    localStorage.setItem(KEY, JSON.stringify([...followed]));
    teamsSummary.textContent = `Followed teams (${followed.size})`;
}

saveFollowed();

function renderEmpty() {
    upcomingMatchList.innerHTML = `<div class="card">Loading...</div>`;
    liveMatchList.innerHTML = `<div class="card" style="opacity:.7">workin on it...</div>`;
}
renderEmpty();

async function fetchUpcoming({ game, teamIDs }) {
    const u = new URL("/api/upcoming-matches", location.origin);
    if (game) {
        u.searchParams.set("game", game);
    }
    if (onlyFollowedEl.checked && teamIDs.length) {
        u.searchParams.set("teamIDs", teamIDs.join(","));
    }
    const res = await fetch(u);
    if (!res.ok) {
        throw new Error("Failed to fetch upcoming matches: " + res.status);
    }
    const data = await res.json();
    if (!data || !Array.isArray(data)) {
        throw new Error("Invalid data format received");
    }
    return data;
}

function renderUpcoming(list, mount) {
    console.log(list);
    if (!list || !list.length) {
        mount.innerHTML = `<div class="card">No matches...</div>`;
        return;
    }
    mount.innerHTML = list
        .map((m) => {
            const a = m.opponents?.[0]?.opponent ?? {};
            const b = m.opponents?.[1]?.opponent ?? {};
            const when = m.scheduled_at ? new Date(m.scheduled_at).toLocaleString() : "TBD";
            const lName = m.league?.name ?? "";

            const aLogo = a.image_url
                ? `<img class="team-logo" src="${a.image_url}" alt="${a.name || "Team"} logo" loading="lazy" decoding="async">`
                : `<span class="team-chip">${a.acronym || (a.name?.[0] ?? "??")}</span>`;

            const bLogo = b.image_url
                ? `<img class="team-logo" src="${b.image_url}" alt="${b.name || "Team"} logo" loading="lazy" decoding="async">`
                : `<span class="team-chip">${b.acronym || (b.name?.[0] ?? "??")}</span>`;

            const leagueLogo = m.league?.image_url
                ? `<img class="league-logo" src="${m.league.image_url}" alt="${lName}" loading="lazy" decoding="async">`
                : "";

            return `
      <article class="card"
               data-match-id="${m.id}"
               data-a-id="${a.id ?? ""}"
               data-b-id="${b.id ?? ""}">
        <div class="card-left">
          ${aLogo}
          <span class="team-name">${a.name ?? "TBD"}</span>
          <span class="vs">vs</span>
          ${bLogo}
          <span class="team-name">${b.name ?? "TBD"}</span>
        </div>
        <div class="card-right">
          <div class="league">${leagueLogo}<span>${m.league?.name ?? ""}</span></div>
          <time datetime="${m.begin_at || ""}">${when}</time>
        </div>
      </article>
    `;
        })
        .join("");
}

async function reload() {
    renderEmpty();
    try {
        const game = gameEl.value || "lol";
        const ids = [...followed]; // turn Set → array
        const upcoming = await fetchUpcoming({ game, teamIDs: ids });
        renderUpcoming(upcoming, upcomingMatchList);
    } catch (e) {
        console.error(e);
        upcomingMatchList.innerHTML = `<div class="card">Error loading matches</div>`;
    }
}

async function searchTeams(query, game, limit = 20, page = 1) {
    const u = new URL("/api/teams", location.origin);
    u.searchParams.set("q", query);
    u.searchParams.set("game", game);
    u.searchParams.set("limit", limit);
    u.searchParams.set("page", page);
    const res = await fetch(u);
    if (!res.ok) {
        throw new Error("Failed to fetch teams" + res.status);
    }
    return res.json();
}

const debounce = (fn, delay = 250) => {
    let timeout;
    return (...args) => {
        clearTimeout(timeout);
        timeout = setTimeout(() => {
            fn(...args);
        }, delay);
    };
};

const onTeamSearch = debounce(async () => {
    const query = teamSearchEl.value.trim();
    teamsList.innerHTML = "";
    if (query.length < 2) {
        return;
    }

    try {
        const items = await searchTeams(query, gameEl.value)
        if (!items.length) {
            teamsList.innerHTML = `<div style="opacity:.7;padding:.5rem 0">No teams found</div>`;
            return;
        }
        teamsList.innerHTML = items.map(team => `
            <label style="display:flex;gap:.5rem;align-items:center;padding:.15rem 0">
                <input type="checkbox" data-id="${team.id}" ${followed.has(team.id) ? "checked": ""}/>
                <span>${team.name}</span>
            </label>`).join("");
        teamsList.querySelectorAll('input[type="checkbox"]').forEach(cb =>{
            cb.addEventListener("change", (e) => {
                const id = Number(e.target.dataset.id);
                if (e.target.checked) { 
                    followed.add(id); 
                }
                else {
                    followed.delete(id);
                }
                saveFollowed();
            })
        })
    } catch (e) {
        console.error(e);
    }
}, 250);


refreshBtn.addEventListener("click", reload);
gameEl.addEventListener("change", reload);
//regionEl.addEventListener("change", reload);
onlyFollowedEl.addEventListener("change", reload);
teamSearchEl.addEventListener("input", onTeamSearch);

teamsSelectAllBtn.addEventListener("click", () => {
    teamsList.querySelectorAll('input[type="checkbox"]').forEach( cb => {
        cb.checked = true;
        followed.add(Number(cb.dataset.id));
    });
    saveFollowed();
})
teamsSaveBtn.addEventListener("click", () => {
    teamsMenu.open = false;
    reload();
})
teamsClearBtn.addEventListener("click", () => {
    followed.clear();
    saveFollowed();
    teamsList.querySelectorAll('input[type="checkbox"]').forEach( cb => 
        cb.checked = false
    );
})

reload();


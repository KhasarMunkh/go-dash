async function loadMatches() {
    const res = await fetch("/api/upcoming-matches");
    if (!res.ok) {
        document.getElementById("matches").innerHTML = "Failed to load data";
        return;
    }

    const data = await res.json();
    const container = document.getElementById("matches");
    container.innerHTML = "";

    if (data.length === 0) {
        container.innerHTML = "<p>No upcoming matches.</p>";
        return;
    }
    console.log(data);

    data.forEach((match) => {
        const card = document.createElement("div");
        card.className = "match-card";

        const date = new Date(match.begin_at);
        const timeString = date.toLocaleString();

        const teamHTML = match.opponents
            .map(
                (o) => `
      <div class="team">
        <img src="${o.opponent.image_url || ""}" alt="${o.opponent.name}" />
        <p>${o.opponent.name} (${o.opponent.acronym})</p>
      </div>
    `,
            )
            .join("");

        card.innerHTML = `
      <h2>${match.name}</h2>
      <p>Start Time: ${timeString}</p>
      <div class="teams">${teamHTML}</div>
    `;

        container.appendChild(card);
    });
}

window.addEventListener("DOMContentLoaded", loadMatches);

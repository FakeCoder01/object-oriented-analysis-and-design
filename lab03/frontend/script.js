const API_BASE = "http://localhost:8080/api";

const canvas = document.getElementById("canvas");
const modeSelect = document.getElementById("modeSelect");
const shapeType = document.getElementById("shapeType");
const shapeColor = document.getElementById("shapeColor");
const shapeSize = document.getElementById("shapeSize");
const currentModeLabel = document.getElementById("currentModeLabel");

async function fetchShapes() {
  const mode = modeSelect.value;
  try {
    const response = await fetch(`${API_BASE}/${mode}/shapes`);
    const shapes = await response.json();
    renderShapes(shapes);
  } catch (e) {
    console.error("Failed to fetch shapes:", e);
  }
}

function renderShapes(shapes) {
  canvas.innerHTML = "";
  shapes.forEach((shape) => {
    const el = document.createElement("div");
    el.className = "shape";
    el.style.left = `${shape.x}px`;
    el.style.top = `${shape.y}px`;
    el.style.width = `${shape.size}px`;
    el.style.height = `${shape.size}px`;
    el.style.backgroundColor = shape.color;

    if (shape.type === "circle") {
      el.style.borderRadius = "50%";
    }

    canvas.appendChild(el);
  });
}

canvas.addEventListener("click", async (e) => {
  const rect = canvas.getBoundingClientRect();
  const x = e.clientX - rect.left;
  const y = e.clientY - rect.top;

  const shape = {
    type: shapeType.value,
    x: Math.round(x),
    y: Math.round(y),
    size: parseInt(shapeSize.value),
    color: shapeColor.value,
  };

  const mode = modeSelect.value;
  try {
    await fetch(`${API_BASE}/${mode}/add`, {
      method: "POST",
      body: JSON.stringify(shape),
    });
    fetchShapes();
  } catch (e) {
    console.error("Error adding shape:", e);
  }
});

document.getElementById("undoBtn").addEventListener("click", async () => {
  const mode = modeSelect.value;
  await fetch(`${API_BASE}/${mode}/undo`, { method: "POST" });
  fetchShapes();
});

document.getElementById("redoBtn").addEventListener("click", async () => {
  const mode = modeSelect.value;
  await fetch(`${API_BASE}/${mode}/redo`, { method: "POST" });
  fetchShapes();
});

document.getElementById("clearBtn").addEventListener("click", async () => {
  const mode = modeSelect.value;
  await fetch(`${API_BASE}/${mode}/clear`, { method: "POST" });
  fetchShapes();
});

modeSelect.addEventListener("change", () => {
  currentModeLabel.textContent =
    modeSelect.options[modeSelect.selectedIndex].text;
  fetchShapes();
});

// init load
fetchShapes();

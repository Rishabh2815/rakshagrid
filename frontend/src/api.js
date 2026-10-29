// API client for RakshaGrid
// Handles authentication, request formatting, and WebSocket connections
// Uses /api/ paths for Vercel services deployment

const API_BASE = "/api";

export async function loginUser({ email, password, role }) {
  const response = await fetch(`${API_BASE}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password, role }),
  });

  if (!response.ok) {
    throw new Error("Login failed");
  }

  return response.json();
}

export async function raiseSOS(latitude, longitude, message) {
  const token = localStorage.getItem("jwt_token");
  const response = await fetch(`${API_BASE}/sos/raise`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ latitude, longitude, message }),
  });

  if (!response.ok) {
    throw new Error("Failed to raise SOS");
  }

  return response.json();
}

export async function fetchIncidents() {
  const token = localStorage.getItem("jwt_token");
  const response = await fetch(`${API_BASE}/incidents`, {
    headers: { Authorization: `Bearer ${token}` },
  });

  if (!response.ok) {
    throw new Error("Failed to fetch incidents");
  }

  return response.json();
}

export function setupWebSocket(onMessage) {
  const token = localStorage.getItem("jwt_token");
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  // Use same-origin WebSocket — Vercel rewrites /api/stream to backend
  const ws = new WebSocket(`${protocol}//${window.location.host}/api/stream?token=${token}`);

  ws.onopen = () => {
    console.log("WebSocket connected");
  };

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      onMessage(data);
    } catch (err) {
      console.error("Failed to parse WebSocket message:", err);
    }
  };

  ws.onerror = (error) => {
    console.error("WebSocket error:", error);
  };

  ws.onclose = () => {
    console.log("WebSocket disconnected. Attempting to reconnect in 5s...");
    setTimeout(() => setupWebSocket(onMessage), 5000);
  };

  return ws;
}

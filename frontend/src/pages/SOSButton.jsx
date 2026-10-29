import { useState, useEffect } from "react";
import "./SOSButton.css";

export default function SOSButton() {
  const [loading, setLoading] = useState(false);
  const [status, setStatus] = useState("ready");
  const [location, setLocation] = useState(null);
  const [error, setError] = useState("");

  useEffect(() => {
    // Get user's location on mount
    if ("geolocation" in navigator) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          setLocation({
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
          });
        },
        (err) => {
          setError("Location access denied. SOS may be less accurate.");
        }
      );
    }
  }, []);

  const handleSOSClick = async () => {
    setLoading(true);
    setStatus("sending");
    setError("");

    try {
      const token = localStorage.getItem("jwt_token");
      const response = await fetch("/v1/sos/raise", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          latitude: location?.latitude || 0,
          longitude: location?.longitude || 0,
          message: "Emergency assistance needed",
        }),
      });

      if (response.ok) {
        setStatus("sent");
        setTimeout(() => setStatus("ready"), 3000);
      } else {
        setError("Failed to send SOS");
        setStatus("ready");
      }
    } catch (err) {
      setError(err.message || "Network error");
      setStatus("ready");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="sos-container">
      <h2>Emergency SOS</h2>

      {location && (
        <div className="location-info">
          📍 Location: {location.latitude.toFixed(4)}, {location.longitude.toFixed(4)}
        </div>
      )}

      {error && <div className="error-banner">{error}</div>}

      <button
        className={`sos-button ${status}`}
        onClick={handleSOSClick}
        disabled={loading || status === "sent"}
      >
        {status === "ready" && "🆘 PRESS FOR HELP"}
        {status === "sending" && "Sending SOS..."}
        {status === "sent" && "✓ SOS Sent"}
      </button>

      <p className="disclaimer">
        Your location will be shared with the nearest responder. Help will arrive as soon as possible.
      </p>
    </div>
  );
}

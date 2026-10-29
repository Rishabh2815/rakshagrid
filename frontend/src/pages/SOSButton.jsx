import { useState, useEffect } from "react";
import "./SOSButton.css";

export default function SOSButton() {
  const [loading, setLoading] = useState(false);
  const [status, setStatus] = useState("ready");
  const [location, setLocation] = useState(null);
  const [error, setError] = useState("");

  useEffect(() => {
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
      const response = await fetch("/api/sos/raise", {
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
    <div className="sos-page">
      <div className="sos-container">
        <h2 className="sos-title">
          Emergency Alert System
        </h2>
        <p className="sos-subtitle">
          Get instant help from nearby responders when you need it most
        </p>

        {location && (
          <div className="location-info">
            <span className="location-badge">📍 Location detected</span>
            <p className="location-coords">
              {location.latitude.toFixed(4)}, {location.longitude.toFixed(4)}
            </p>
          </div>
        )}

        {error && <div className="error-message">{error}</div>}

        <div className="button-wrapper">
          <button
            className={`sos-button ${status}`}
            onClick={handleSOSClick}
            disabled={loading || status === "sent"}
          >
            {status === "ready" && (
              <>
                <span className="button-icon">🆘</span>
                <span className="button-label">Press for Help</span>
              </>
            )}
            {status === "sending" && (
              <>
                <span className="spinner"></span>
                <span className="button-label">Sending...</span>
              </>
            )}
            {status === "sent" && (
              <>
                <span className="button-icon">✓</span>
                <span className="button-label">Sent</span>
              </>
            )}
          </button>
        </div>

        <p className="sos-disclaimer">
          Your location will be shared with responders. Emergency assistance will arrive shortly.
        </p>
      </div>
    </div>
  );
}

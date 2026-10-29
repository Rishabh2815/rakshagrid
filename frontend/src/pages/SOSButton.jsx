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
    <div className="sos-container">
      <div className="sos-content">
        <h2 className="sos-title">Emergency Alert System</h2>
        <p className="sos-description">Press the button below to send an emergency SOS to nearby responders</p>

        {location && (
          <div className="location-info">
            <div className="location-icon">📍</div>
            <div className="location-text">
              <p className="location-label">Current Location</p>
              <p className="location-coords">{location.latitude.toFixed(4)}, {location.longitude.toFixed(4)}</p>
            </div>
          </div>
        )}

        {error && <div className="error-banner">{error}</div>}

        <div className="button-container">
          <button
            className={`sos-button ${status}`}
            onClick={handleSOSClick}
            disabled={loading || status === "sent"}
          >
            <span className="button-inner">
              {status === "ready" && (
                <>
                  <span className="button-icon">🆘</span>
                  <span className="button-text">Press for Help</span>
                </>
              )}
              {status === "sending" && (
                <>
                  <span className="spinner"></span>
                  <span className="button-text">Sending SOS...</span>
                </>
              )}
              {status === "sent" && (
                <>
                  <span className="button-icon">✓</span>
                  <span className="button-text">SOS Sent</span>
                </>
              )}
            </span>
          </button>
        </div>

        <p className="sos-disclaimer">
          Your location will be shared with the nearest responder. Emergency assistance will arrive as soon as possible.
        </p>

        <div className="info-cards">
          <div className="info-card">
            <div className="card-icon">⚡</div>
            <div className="card-content">
              <h3>Instant Response</h3>
              <p>Real-time alert to nearest responders</p>
            </div>
          </div>
          <div className="info-card">
            <div className="card-icon">📍</div>
            <div className="card-content">
              <h3>GPS Tracking</h3>
              <p>Precise location sharing</p>
            </div>
          </div>
          <div className="info-card">
            <div className="card-icon">🔒</div>
            <div className="card-content">
              <h3>Secure & Private</h3>
              <p>Your data is protected</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

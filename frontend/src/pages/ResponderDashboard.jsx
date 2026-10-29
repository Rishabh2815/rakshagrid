import "./ResponderDashboard.css";

export default function ResponderDashboard() {
  return (
    <div className="dashboard-container responder-dashboard">
      <div className="dashboard-header">
        <h2>🚨 Responder Dashboard</h2>
        <p>Active incidents waiting for response</p>
      </div>

      <div className="incidents-grid">
        <div className="incident-card">
          <div className="incident-status idle">Waiting</div>
          <h3>No Active Incidents</h3>
          <p>All systems operational. No active SOS requests.</p>
        </div>
      </div>

      <div className="stats-section">
        <div className="stat-card">
          <div className="stat-icon">✓</div>
          <div className="stat-content">
            <p className="stat-label">Response Time</p>
            <p className="stat-value">--:--</p>
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-icon">📍</div>
          <div className="stat-content">
            <p className="stat-label">Coverage Area</p>
            <p className="stat-value">--</p>
          </div>
        </div>
        <div className="stat-card">
          <div className="stat-icon">📊</div>
          <div className="stat-content">
            <p className="stat-label">Responses Today</p>
            <p className="stat-value">0</p>
          </div>
        </div>
      </div>
    </div>
  );
}

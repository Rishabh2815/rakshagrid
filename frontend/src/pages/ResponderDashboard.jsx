import "./ResponderDashboard.css";

export default function ResponderDashboard() {
  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <h2>Active Incidents</h2>
        <p>Respond to emergency calls from students</p>
      </div>

      <div className="incidents-list">
        <div className="empty-state">
          <div className="empty-icon">📋</div>
          <p>No active incidents</p>
          <span>All systems operational</span>
        </div>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-value">--</div>
          <div className="stat-label">Response Time</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">0</div>
          <div className="stat-label">Responses Today</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">--</div>
          <div className="stat-label">Coverage Area</div>
        </div>
      </div>
    </div>
  );
}

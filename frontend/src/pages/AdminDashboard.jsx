import "./ResponderDashboard.css";

export default function AdminDashboard() {
  return (
    <div className="dashboard-container admin-dashboard">
      <div className="dashboard-header">
        <h2>📊 Admin Dashboard</h2>
        <p>Campus safety compliance & analytics</p>
      </div>

      <div className="metrics-grid">
        <div className="metric-card">
          <div className="metric-icon">🆘</div>
          <div className="metric-info">
            <p className="metric-label">Total SOS Alerts</p>
            <p className="metric-value">0</p>
            <p className="metric-change">This month</p>
          </div>
        </div>

        <div className="metric-card">
          <div className="metric-icon">⏱️</div>
          <div className="metric-info">
            <p className="metric-label">Avg Response Time</p>
            <p className="metric-value">--</p>
            <p className="metric-change">seconds</p>
          </div>
        </div>

        <div className="metric-card">
          <div className="metric-icon">✓</div>
          <div className="metric-info">
            <p className="metric-label">Completion Rate</p>
            <p className="metric-value">--%</p>
            <p className="metric-change">resolved</p>
          </div>
        </div>

        <div className="metric-card">
          <div className="metric-icon">👥</div>
          <div className="metric-info">
            <p className="metric-label">Active Responders</p>
            <p className="metric-value">0</p>
            <p className="metric-change">online</p>
          </div>
        </div>
      </div>

      <div className="compliance-section">
        <h3>Compliance Report</h3>
        <div className="compliance-card">
          <p>No incidents recorded yet. Dashboard will update when SOS alerts are created.</p>
        </div>
      </div>
    </div>
  );
}

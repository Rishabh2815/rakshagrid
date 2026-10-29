import "./ResponderDashboard.css";

export default function AdminDashboard() {
  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <h2>Safety Dashboard</h2>
        <p>Campus compliance & incident analytics</p>
      </div>

      <div className="metrics-grid">
        <div className="metric-card">
          <div className="metric-value">0</div>
          <div className="metric-label">Total SOS Alerts</div>
          <div className="metric-detail">This month</div>
        </div>
        <div className="metric-card">
          <div className="metric-value">--%</div>
          <div className="metric-label">Response Rate</div>
          <div className="metric-detail">Average</div>
        </div>
        <div className="metric-card">
          <div className="metric-value">0</div>
          <div className="metric-label">Responders Online</div>
          <div className="metric-detail">Active now</div>
        </div>
      </div>

      <div className="compliance-section">
        <h3>Compliance Report</h3>
        <div className="compliance-empty">
          <p>No incidents recorded. Dashboard updates in real-time.</p>
        </div>
      </div>
    </div>
  );
}

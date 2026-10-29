import { useState } from "react";

export default function AdminDashboard() {
  return (
    <div style={styles.container}>
      <h2>Admin Dashboard</h2>
      <p>Campus administrators can view compliance reports and incident statistics.</p>
      <div style={styles.placeholder}>
        <p>📊 No data available</p>
      </div>
    </div>
  );
}

const styles = {
  container: {
    maxWidth: "800px",
    margin: "0 auto",
    padding: "2rem",
  },
  placeholder: {
    background: "#f0f0f0",
    padding: "3rem",
    borderRadius: "8px",
    textAlign: "center",
    color: "#666",
  },
};

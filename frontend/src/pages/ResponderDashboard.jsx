import { useState } from "react";

export default function ResponderDashboard() {
  return (
    <div style={styles.container}>
      <h2>Responder Dashboard</h2>
      <p>Active incidents will appear here. Responders can accept or decline SOS requests.</p>
      <div style={styles.placeholder}>
        <p>🚨 No active incidents</p>
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

import { useState } from "react";
import { loginUser } from "../api";
import "./LoginForm.css";

export default function LoginForm({ onLogin }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState("user");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const response = await loginUser({ email, password, role });
      if (response.token) {
        onLogin(response.token, role);
      } else {
        setError("Login failed. Please check your credentials.");
      }
    } catch (err) {
      setError(err.message || "An error occurred during login.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <div className="login-background">
        <div className="background-overlay"></div>
        <img 
          src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1200 800'%3E%3Cdefs%3E%3ClinearGradient id='grad1' x1='0%25' y1='0%25' x2='0%25' y2='100%25'%3E%3Cstop offset='0%25' style='stop-color:%23e8dcc8;stop-opacity:0.8' /%3E%3Cstop offset='50%25' style='stop-color:%23c9b8a8;stop-opacity:0.7' /%3E%3Cstop offset='100%25' style='stop-color:%236b8a5e;stop-opacity:0.9' /%3E%3C/linearGradient%3E%3C/defs%3E%3Crect width='1200' height='800' fill='url(%23grad1)'/%3E%3Cpath d='M 0,300 Q 300,250 600,300 T 1200,300 L 1200,800 L 0,800 Z' fill='%234a6a4a' opacity='0.6'/%3E%3C/svg%3E" 
          alt="Nature background" 
          className="background-image"
        />
      </div>

      <div className="login-container">
        <div className="login-header">
          <h1 className="brand-name">🛡️ RakshaGrid</h1>
          <p className="brand-tagline">Emergency Response System</p>
        </div>

        <div className="login-content">
          <h2 className="login-title">
            The best place to find your <span className="italic-text">Inner Peace</span>
          </h2>
          <p className="login-subtitle">
            Emergency assistance when you need it most. Find help from nearby responders instantly.
          </p>

          <form onSubmit={handleSubmit} className="login-form">
            <div className="form-group">
              <input
                type="email"
                placeholder="Email address"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="form-input"
              />
            </div>

            <div className="form-group">
              <input
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="form-input"
              />
            </div>

            <div className="form-group">
              <select 
                value={role} 
                onChange={(e) => setRole(e.target.value)}
                className="form-select"
              >
                <option value="user">I am a User (Student)</option>
                <option value="responder">I am a Responder (Security/Warden)</option>
                <option value="admin">I am an Administrator</option>
              </select>
            </div>

            {error && <div className="error-message">{error}</div>}

            <button 
              type="submit" 
              disabled={loading}
              className="login-button"
            >
              {loading ? "Logging in..." : "Login"}
            </button>
          </form>

          <div className="featured-section">
            <p className="featured-text">Featured as the safest place to go in</p>
            <div className="featured-brands">
              <span className="brand-logo">Forbes</span>
              <span className="brand-logo">MensHealth</span>
              <span className="brand-logo">Bloomberg</span>
              <span className="brand-logo">Washington Post</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

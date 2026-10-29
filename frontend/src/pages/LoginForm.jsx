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
      <div className="login-container">
        {/* Left Panel - Branding */}
        <div className="login-left">
          <div className="brand-card">
            <div className="brand-gradient"></div>
            <h2>Simplify emergency response with our platform</h2>
            <p>Manage your campus safety with our user-friendly emergency response system</p>
            <div className="brand-illustration">
              <div className="illustration-item item-1">📱</div>
              <div className="illustration-item item-2">🛡️</div>
              <div className="illustration-item item-3">⚡</div>
            </div>
          </div>
        </div>

        {/* Right Panel - Login Form */}
        <div className="login-right">
          <div className="login-card">
            <div className="logo-section">
              <div className="logo">🛡️</div>
              <h1>RakshaGrid</h1>
              <p className="tagline">Welcome Back</p>
            </div>

            <form onSubmit={handleSubmit} className="login-form">
              <div className="form-group">
                <input
                  type="email"
                  placeholder="Email"
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
                <a href="#" className="forgot-link">Forgot Password?</a>
              </div>

              <div className="form-group">
                <select 
                  value={role} 
                  onChange={(e) => setRole(e.target.value)}
                  className="form-select"
                >
                  <option value="user">User (Student)</option>
                  <option value="responder">Responder (Security/Warden)</option>
                  <option value="admin">Administrator</option>
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

            <div className="divider">or Login with</div>

            <div className="social-login">
              <button className="social-btn facebook">
                <span>f</span> Facebook
              </button>
              <button className="social-btn google">
                <span>G</span> Google
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

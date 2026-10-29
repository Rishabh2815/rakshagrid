import { useState } from "react";
import LoginForm from "./pages/LoginForm";
import SOSButton from "./pages/SOSButton";
import ResponderDashboard from "./pages/ResponderDashboard";
import AdminDashboard from "./pages/AdminDashboard";
import "./App.css";

export default function App() {
  const [user, setUser] = useState(null);
  const [role, setRole] = useState(null);

  const handleLogout = () => {
    localStorage.removeItem("jwt_token");
    localStorage.removeItem("user_role");
    setUser(null);
    setRole(null);
  };

  if (!user) {
    return <LoginForm onLogin={(token, userRole) => {
      setUser(token);
      setRole(userRole);
      localStorage.setItem("jwt_token", token);
      localStorage.setItem("user_role", userRole);
    }} />;
  }

  return (
    <div className="app">
      <header className="app-header">
        <div className="header-content">
          <h1 className="logo">🛡️ RakshaGrid</h1>
          <button className="logout-btn" onClick={handleLogout}>Logout</button>
        </div>
      </header>

      <main className="app-main">
        {role === "user" && <SOSButton />}
        {role === "responder" && <ResponderDashboard />}
        {role === "admin" && <AdminDashboard />}
      </main>
    </div>
  );
}

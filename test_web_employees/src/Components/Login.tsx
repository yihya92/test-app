import { useState } from "react";
import type React from "react";
import { useNavigate } from "react-router-dom";
import type { IState } from "../App";

interface IProps {
  people: IState["people"];
  setLoggedInEmail: React.Dispatch<React.SetStateAction<string | null>>;
}

const Login = ({ people, setLoggedInEmail }: IProps) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    const emailLower = email.toLowerCase();
    console.log("[Login] entered email:", JSON.stringify(emailLower), "password:", JSON.stringify(password));
    console.log("[Login] people:", people.map(p => ({ Email: p.Email, NewPassword: p.NewPassword })));
    const match = people.some(p => p.Email?.toLowerCase() === emailLower && p.NewPassword === password);
    if (match) {
      setLoggedInEmail(emailLower);
      navigate("/app/list");
    } else {
      setError("Invalid email or password.");
    }
  };

  return (
    <div style={{ display: "flex", justifyContent: "center", alignItems: "top", minHeight: "100vh" }}>
      <div style={{ width: "450px" }}>
        <h1>Employees Login</h1>
        <form onSubmit={handleSubmit}>
          <div>
            <input
              type="email"
              className="AddToList-input"
              placeholder="Email"
              style={{ width: "100%" }}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>
          <div style={{ position: "relative" }}>
            <input
              type={showPassword ? "text" : "password"}
              className="AddToList-input"
              placeholder="Password"
              style={{ width: "100%" }}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <button
              type="button"
              onClick={() => setShowPassword(prev => !prev)}
              style={{ position: "absolute", right: "8px", top: "50%", transform: "translateY(-50%)", background: "none", border: "none", cursor: "pointer", fontSize: "18px" }}
            >
              {showPassword ? "🙈" : "👁️"}
            </button>
          </div>
          {error && <p style={{ color: "red", fontSize: "14px" }}>{error}</p>}
          <button type="submit" className="AddToList-btn">Login</button>
          <a onClick={() => navigate("/app/signup")}>
            Not registered yet, click here
          </a>
        </form>
      </div>
    </div>
  );
};

export default Login;

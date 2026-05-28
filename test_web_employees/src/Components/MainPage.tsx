import { Link } from 'react-router-dom'

function MainPage() {
  return (
    <div style={{ textAlign: "center" }}>
      <h1>Welcome to our Company</h1>
      <Link to="/app/login">
        <button className="AddToList-btn" style={{ width: "40%" }}>Enter</button>
      </Link>
    </div>
  );
}
export default MainPage;
import { useState, useEffect, useCallback } from 'react'
import { Routes, Route } from 'react-router-dom'
import './App.css'
import AddToList from './Components/Signup'
import MainPage from './Components/MainPage'
import Login from './Components/Login'
import List from './Components/List'

export const API_BASE = "http://localhost:9900"
export const API_EMPLOYEES = `${API_BASE}/EMP/V1/HTTP_Employees/`

export interface IState {
  people: {
    Id?: number
    Key?: string
    Login?: string
    Name: string
    Email: string
    Age: number
    Position: string
    PhoneNumber: string
    Department: string
    Unit: string
    NewPassword: string
    ConfirmPassword: string
    EmployeeType: "Admin" | "Employee"
  }[]
}

interface BackendEmployee {
  Id: number
  Key: string
  Login: string
  Name: string
  name?: string
  Email: string
  email?: string
  Age: number
  age?: number
  Position: string
  PhoneNumber: string
  Department: string
  Unit: string
  NewPassword: string
  ConfirmPassword: string
}

function toFrontend(e: BackendEmployee): IState["people"][0] {
  return {
    Id: e.Id,
    Key: e.Key,
    Login: e.Login,
    Name: e.Name ?? e.name ?? "",
    Email: e.Email ?? e.email ?? "",
    Age: e.Age ?? e.age ?? 0,
    Position: e.Position,
    PhoneNumber: e.PhoneNumber,
    Department: e.Department,
    Unit: e.Unit,
    NewPassword: e.NewPassword,
    ConfirmPassword: e.ConfirmPassword,
    EmployeeType: e.Id === 1 ? "Admin" : "Employee",
  }
}

function App() {
  const [people, setPeople] = useState<IState["people"]>([])
  const [loggedInEmail, setLoggedInEmail] = useState<string | null>(() =>
    localStorage.getItem('loggedInEmail')
  )
  const [loading, setLoading] = useState(true)
  const [fetchError, setFetchError] = useState<string | null>(null)

  useEffect(() => {
    if (loggedInEmail) {
      localStorage.setItem('loggedInEmail', loggedInEmail)
    } else {
      localStorage.removeItem('loggedInEmail')
    }
  }, [loggedInEmail])

  const fetchPeople = useCallback(() => {
    setLoading(true)
    setFetchError(null)
    fetch(API_EMPLOYEES)
      .then(r => r.json())
      .then(data => {
        if (data.Status === "successful" && Array.isArray(data.Data)) {
          const sorted = (data.Data as BackendEmployee[]).sort((a, b) => a.Id - b.Id)
          setPeople(sorted.map(toFrontend))
        } else if (data.ErrorDescription) {
          setFetchError(data.ErrorDescription)
        }
      })
      .catch(err => setFetchError(String(err)))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => { fetchPeople() }, [])

  if (loading) {
    return <div style={{ display: "flex", justifyContent: "center", alignItems: "center", minHeight: "100vh", fontSize: "1.2rem" }}>Loading employees...</div>
  }

  if (fetchError) {
    return (
      <div style={{ display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center", minHeight: "100vh", gap: "1rem" }}>
        <p style={{ color: "red" }}>Failed to load employees: {fetchError}</p>
        <button className="AddToList-btn" onClick={fetchPeople}>Retry</button>
      </div>
    )
  }

  return (
    <Routes>
      <Route path="/" element={<MainPage />} />
      <Route path="/app/list" element={
        <List
          key={loggedInEmail ?? ""}
          people={people}
          setPeople={setPeople}
          loggedInEmail={loggedInEmail}
          setLoggedInEmail={setLoggedInEmail}
          fetchPeople={fetchPeople}
        />
      } />
      <Route path="/app/*" element={
        <div className="App">
          <Routes>
            <Route path="login" element={<Login people={people} setLoggedInEmail={setLoggedInEmail} />} />
            <Route path="signup" element={<AddToList people={people} fetchPeople={fetchPeople} />} />
          </Routes>
        </div>
      } />
    </Routes>
  );
}

export default App

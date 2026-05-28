import React, { useState, type JSX } from "react";
import { useNavigate } from "react-router-dom";
import { type IState, API_EMPLOYEES } from "../App";

interface IProps {
    people: IState["people"];
    setPeople: React.Dispatch<React.SetStateAction<IState["people"]>>;
    loggedInEmail: string | null;
    setLoggedInEmail: React.Dispatch<React.SetStateAction<string | null>>;
    fetchPeople?: () => void;
}

type Person = IState["people"][0];

function toBackendPayload(person: Person, edits: Person) {
    return {
        Id: person.Id,
        Login: person.Login,
        name: edits.Name,
        email: edits.Email,
        age: edits.Age,
        Position: edits.Position,
        PhoneNumber: edits.PhoneNumber,
        Department: edits.Department,
        Unit: edits.Unit,
    };
}

const List: React.FC<IProps> = ({ people, setPeople, loggedInEmail, setLoggedInEmail, fetchPeople }) => {
    const currentUser = people.find(p => p.Email.toLowerCase() === loggedInEmail);
    const isAdmin = currentUser?.EmployeeType === "Admin";
    const navigate = useNavigate();
    const [selected, setSelected] = useState<Set<number>>(new Set());
    const [editingIndex, setEditingIndex] = useState<number | null>(null);
    const [editValues, setEditValues] = useState<Person | null>(null);
    const [actionError, setActionError] = useState<string | null>(null);

    const refresh = () => { if (fetchPeople) fetchPeople(); };

    const toggleSelect = (i: number) => {
        setSelected(prev => {
            const next = new Set(prev);
            next.has(i) ? next.delete(i) : next.add(i);
            return next;
        });
    };

    const toggleSelectAll = () => {
        setSelected(
            selected.size === people.length
                ? new Set()
                : new Set(people.map((_, i) => i))
        );
    };

    const handleDelete = (i: number) => {
        if (!confirm(`Delete ${people[i].Name}?`)) return;
        const login = people[i].Login;
        if (!login) { setActionError("Cannot delete: employee has no Login"); return; }
        setActionError(null);
        fetch(`${API_EMPLOYEES}${encodeURIComponent(login)}`, { method: "DELETE" })
            .then(async r => {
                if (!r.ok) {
                    const data = await r.json().catch(() => ({}));
                    setActionError(data.ErrorDescription || "Delete failed");
                    return;
                }
                refresh();
                setSelected(prev => {
                    const next = new Set<number>();
                    prev.forEach(s => { if (s < i) next.add(s); else if (s > i) next.add(s - 1); });
                    return next;
                });
                if (editingIndex === i) setEditingIndex(null);
            })
            .catch(err => setActionError(String(err)));
    };

    const handleDeleteSelected = () => {
        if (!confirm(`Delete ${selected.size} employee(s)?`)) return;
        setActionError(null);
        const toDelete = Array.from(selected).map(i => people[i]).filter(p => p.Login);
        Promise.all(
            toDelete.map(p =>
                fetch(`${API_EMPLOYEES}${encodeURIComponent(p.Login!)}`, { method: "DELETE" })
            )
        )
            .then(results => {
                const failed = results.filter(r => !r.ok);
                if (failed.length > 0) setActionError(`${failed.length} deletion(s) failed`);
                refresh();
                if (editingIndex !== null && selected.has(editingIndex)) setEditingIndex(null);
                setSelected(new Set());
            })
            .catch(err => setActionError(String(err)));
    };

    const handleEdit = (i: number) => {
        setEditingIndex(i);
        setEditValues({ ...people[i] });
        setActionError(null);
    };

    const handleSave = (i: number) => {
        if (!editValues) return;
        const person = people[i];
        const login = person.Login;
        if (!login) { setActionError("Cannot save: employee has no Login"); return; }
        setActionError(null);
        fetch(`${API_EMPLOYEES}${encodeURIComponent(login)}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(toBackendPayload(person, editValues)),
        })
            .then(async r => {
                if (!r.ok) {
                    const data = await r.json().catch(() => ({}));
                    setActionError(data.ErrorDescription || "Save failed");
                    return;
                }
                setPeople(prev => prev.map((p, idx) => (idx === i ? editValues : p)));
                setEditingIndex(null);
                setEditValues(null);
                refresh();
            })
            .catch(err => setActionError(String(err)));
    };

    const handleCancel = () => {
        setEditingIndex(null);
        setEditValues(null);
        setActionError(null);
    };

    const handleEditChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setEditValues(prev => (prev ? { ...prev, [name]: value } : prev));
    };

    const renderEmpty = (): JSX.Element => (
        <div className="List-empty">
            <p>No employees found.</p>
            {isAdmin && (
                <button className="AddToList-btn logout-btn" onClick={() => navigate("/app/signup")}>
                    Add Employee
                </button>
            )}
        </div>
    );

    const renderRows = (): JSX.Element[] =>
        people.map((person, i) => {
            const isEditing = editingIndex === i;
            const isSelf = person.Email.toLowerCase() === loggedInEmail;
            return (
                <tr key={i} className="List-row">
                    {isAdmin && (
                        <td>
                            <input
                                type="checkbox"
                                checked={selected.has(i)}
                                onChange={() => toggleSelect(i)}
                                className="List-checkbox"
                            />
                        </td>
                    )}
                    {isAdmin && isEditing && editValues ? (
                        <>
                            <td><input name="Name" value={editValues.Name} onChange={handleEditChange} placeholder="Name" className="AddToList-input" /></td>
                            <td><input name="Email" value={editValues.Email} onChange={handleEditChange} placeholder="Email" className="AddToList-input" /></td>
                            <td><input name="Position" value={editValues.Position} onChange={handleEditChange} placeholder="Position" className="AddToList-input" /></td>
                            <td><input name="Department" value={editValues.Department} onChange={handleEditChange} placeholder="Department" className="AddToList-input" /></td>
                            <td><input name="Unit" value={editValues.Unit} onChange={handleEditChange} placeholder="Unit" className="AddToList-input" /></td>
                            <td><input name="PhoneNumber" value={editValues.PhoneNumber} onChange={handleEditChange} placeholder="Phone" className="AddToList-input" /></td>
                            <td><input name="EmployeeType" value={editValues.EmployeeType} onChange={handleEditChange} placeholder="Employee Type" className="AddToList-input" /></td>
                        </>
                    ) : (
                        <>
                            <td>{person.Name}</td>
                            <td>{person.Email}</td>
                            <td>{person.Position}</td>
                            <td>{person.Department}</td>
                            <td>{person.Unit}</td>
                            <td>{person.PhoneNumber}</td>
                            <td>
                                <span style={{
                                    fontSize: "12px",
                                    background: person.EmployeeType === "Admin" ? "#1976d2" : "#4caf50",
                                    color: "#fff",
                                    borderRadius: "4px",
                                    padding: "2px 8px",
                                }}>
                                    {person.EmployeeType}
                                </span>
                            </td>
                        </>
                    )}
                    {isAdmin && (
                        <td className="List-actions">
                            {isEditing ? (
                                <>
                                    <button className="List-btn save-btn" onClick={() => handleSave(i)}>Save</button>
                                    <button className="List-btn cancel-btn" onClick={handleCancel}>Cancel</button>
                                </>
                            ) : (
                                <>
                                    <button className="List-btn edit-btn" onClick={() => handleEdit(i)}>Edit</button>
                                    {!isSelf && (
                                        <button className="List-btn delete-btn" onClick={() => handleDelete(i)}>Delete</button>
                                    )}
                                </>
                            )}
                        </td>
                    )}
                </tr>
            );
        });

    return (
        <div>
            <div className="List-topbar">
                <div style={{ display: "flex", alignItems: "center", gap: "0.75rem" }}>
                    <h1>List of Employees</h1>
                </div>
                <div style={{ display: "flex", alignItems: "center", gap: "0.75rem" }}>
                    {currentUser && (
                        <span style={{ fontSize: "14px", color: "#555" }}>
                            {currentUser.Name}&nbsp;
                            <span style={{
                                fontSize: "11px",
                                background: isAdmin ? "#1976d2" : "#4caf50",
                                color: "#fff",
                                borderRadius: "4px",
                                padding: "2px 7px",
                            }}>
                                {currentUser.EmployeeType}
                            </span>
                        </span>
                    )}
                    {isAdmin && selected.size > 0 && (
                        <button className="AddToList-btn delete-selected-btn logout-btn" onClick={handleDeleteSelected}>
                            Delete Selected ({selected.size})
                        </button>
                    )}
                    <button className="AddToList-btn logout-btn" onClick={() => { setLoggedInEmail(null); navigate("/"); }}>Logout</button>
                </div>
            </div>
            {!isAdmin && (
                <p style={{ margin: "0.5rem 1rem", fontSize: "13px", color: "#888" }}>
                    You are viewing this list as an Employee. Contact your admin to make any changes.
                </p>
            )}
            {actionError && (
                <p style={{ margin: "0.5rem 1rem", fontSize: "13px", color: "red" }}>{actionError}</p>
            )}
            {people.length === 0 ? renderEmpty() : (
                <table className="List-table" style={{ width: "100%", borderCollapse: "collapse" }}>
                    <thead>
                        <tr className="List-header-row">
                            {isAdmin && (
                                <th>
                                    <input
                                        type="checkbox"
                                        checked={selected.size === people.length}
                                        onChange={toggleSelectAll}
                                        className="List-checkbox"
                                        title="Select all"
                                    />
                                </th>
                            )}
                            <th>Name</th>
                            <th>Email</th>
                            <th>Position</th>
                            <th>Department</th>
                            <th>Unit</th>
                            <th>Phone</th>
                            <th>Employee Type</th>
                            {isAdmin && <th>Actions</th>}
                        </tr>
                    </thead>
                    <tbody>
                        {renderRows()}
                    </tbody>
                </table>
            )}
        </div>
    );
};

export default List;

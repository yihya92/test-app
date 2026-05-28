import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { type IState as Props, API_EMPLOYEES } from "../App";
import { DatePicker } from "@mui/x-date-pickers/DatePicker";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import dayjs, { type Dayjs } from "dayjs";
import {
    TextField,
    InputAdornment,
    Select,
    MenuItem,
    FormHelperText,
    FormControl,
    Box,
} from "@mui/material";

interface IProps {
    people: Props["people"]
    fetchPeople: () => void
}

const COUNTRY = { code: "+961", flag: "🇱🇧", label: "Lebanon" };

const isValidLebanesPhone = (phone: string): boolean =>
    /^\d{7,8}$/.test(phone);

const deriveEmail = (name: string): string => {
    const parts = name.trim().split(/\s+/).filter(Boolean);
    if (parts.length < 2) return "";
    const first = parts[0][0].toLowerCase();
    const last = parts[parts.length - 1].toLowerCase();
    return `${first}${last}@africell.com`;
};

const isValidAfricellEmail = (email: string): boolean =>
    /^[^\s@]+@africell\.com$/.test(email);

const passwordRules = [
    { test: (p: string) => p.length >= 8,           msg: "At least 8 characters" },
    { test: (p: string) => /[A-Z]/.test(p),         msg: "At least one uppercase letter" },
    { test: (p: string) => /[a-z]/.test(p),         msg: "At least one lowercase letter" },
    { test: (p: string) => /[0-9]/.test(p),         msg: "At least one number" },
    { test: (p: string) => /[^A-Za-z0-9]/.test(p), msg: "At least one special character" },
];

const getPasswordErrors = (p: string): string[] =>
    passwordRules.filter(r => !r.test(p)).map(r => r.msg);

const err = (msg: string) => (
    <p style={{ color: "red", fontSize: "13px", margin: "0 0 8px 0" }}>{msg}</p>
);

const AddToList: React.FC<IProps> = ({ people, fetchPeople }) => {
    const navigate = useNavigate();
    const [submitError, setSubmitError] = useState("");
    const [input, setInput] = useState({ name: "", age: 0, position: "", department: "", unit: "", newPassword: "", confirmPassword: "" });
    const [birthDate, setBirthDate] = useState<Dayjs | null>(null);
    const [localPhone, setLocalPhone] = useState("");

    const [nameError, setNameError] = useState("");
    const [emailError, setEmailError] = useState("");
    const [ageError, setAgeError] = useState("");
    const [positionError, setPositionError] = useState("");
    const [departmentError, setDepartmentError] = useState("");
    const [unitError, setUnitError] = useState("");
    const [phoneError, setPhoneError] = useState("");
    const [passwordErrors, setPasswordErrors] = useState<string[]>([]);
    const [confirmPasswordError, setConfirmPasswordError] = useState("");
    const [duplicateError, setDuplicateError] = useState("");

    const [email, setEmail] = useState("");
    const [showNewPassword, setShowNewPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
        const { name, value } = e.target;
        setInput(prev => ({ ...prev, [name]: value }));
        if (name === "name") {
            if (value.trim()) setNameError("");
            setEmail(deriveEmail(value));
            setEmailError("");
            setDuplicateError("");
        }
        if (name === "position" && value.trim()) setPositionError("");
        if (name === "department" && value.trim()) setDepartmentError("");
        if (name === "unit" && value.trim()) setUnitError("");
    };

    const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
        setEmail(e.target.value);
        if (emailError) setEmailError("");
    };

    const handleEmailBlur = (): void => {
        if (!email) { setEmailError("Email is required"); return; }
        if (!isValidAfricellEmail(email))
            setEmailError("Email must end with @africell.com");
        else
            setEmailError("");
    };

    const handleDateChange = (date: Dayjs | null): void => {
        setBirthDate(date);
        const age = date ? dayjs().diff(date, "year") : 0;
        setInput(prev => ({ ...prev, age }));
        if (age > 0) setAgeError("");
    };

    const handlePhoneChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
        const digits = e.target.value.replace(/\D/g, "").slice(0, 8);
        setLocalPhone(digits);
        if (phoneError) setPhoneError("");
    };

    const handlePhoneBlur = (): void => {
        if (!localPhone)
            setPhoneError("Phone number is required");
        else if (!isValidLebanesPhone(localPhone))
            setPhoneError("Enter a valid Lebanese number (7–8 digits)");
        else
            setPhoneError("");
    };

    const handlePasswordBlur = (): void => {
        setPasswordErrors(getPasswordErrors(input.newPassword));
    };

    const handleConfirmPasswordBlur = (): void => {
        setConfirmPasswordError(
            input.confirmPassword && input.confirmPassword !== input.newPassword
                ? "Passwords do not match"
                : ""
        );
    };

    const handleClick = (): void => {
        let valid = true;

        if (!input.name.trim()) { setNameError("Name is required"); valid = false; }
        if (!email || !isValidAfricellEmail(email)) { setEmailError("Email must end with @africell.com"); valid = false; }
        if (input.age < 1) { setAgeError("Date of birth is required"); valid = false; }
        if (!input.position.trim()) { setPositionError("Position is required"); valid = false; }
        if (!input.department.trim()) { setDepartmentError("Department is required"); valid = false; }
        if (!input.unit.trim()) { setUnitError("Unit is required"); valid = false; }
        if (!localPhone) { setPhoneError("Phone number is required"); valid = false; }
        else if (!isValidLebanesPhone(localPhone)) { setPhoneError("Enter a valid Lebanese number (7–8 digits)"); valid = false; }

        const pwErrors = getPasswordErrors(input.newPassword);
        if (pwErrors.length > 0) { setPasswordErrors(pwErrors); valid = false; }

        if (input.newPassword !== input.confirmPassword) { setConfirmPasswordError("Passwords do not match"); valid = false; }

        if (!valid) return;

        if (people.some(p => p.Name.toLowerCase() === input.name.toLowerCase())) {
            setDuplicateError("This name is already registered");
            return;
        }

        const login = email.replace("@africell.com", "");
        const payload = {
            Login: login,
            Name: input.name,
            Email: email,
            Age: input.age,
            Position: input.position,
            PhoneNumber: localPhone ? `${COUNTRY.code}${localPhone}` : "",
            Department: input.department,
            Unit: input.unit,
            NewPassword: input.newPassword,
            ConfirmPassword: input.confirmPassword,
        };

        fetch(API_EMPLOYEES, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
        })
            .then(async r => {
                if (!r.ok) {
                    const data = await r.json().catch(() => ({}));
                    setSubmitError(data.ErrorDescription || "Signup failed");
                    return;
                }
                fetchPeople();
                navigate("/app/login");
            })
            .catch(err => setSubmitError(String(err)));
    };

    return (
        <div className="AddToList">
            <h1 style={{ whiteSpace: "nowrap", fontSize: "1.6rem" }}>Employees Signup</h1>

            <input
                type="text"
                placeholder="Name *"
                className="AddToList-input"
                value={input.name}
                onChange={handleChange}
                onBlur={() => { if (!input.name.trim()) setNameError("Name is required"); }}
                name="name"
            />
            {nameError && err(nameError)}
            {duplicateError && err(duplicateError)}

            <TextField
                label="Email *"
                type="email"
                fullWidth
                value={email}
                onChange={handleEmailChange}
                onBlur={handleEmailBlur}
                error={!!emailError}
                helperText={emailError}
                placeholder="jdoe@africell.com"
                sx={{ mb: "0.3rem" }}
            />

            <LocalizationProvider dateAdapter={AdapterDayjs}>
                <DatePicker
                    label="Date of Birth *"
                    value={birthDate}
                    onChange={handleDateChange}
                    disableFuture
                    sx={{ width: "100%", mb: "0.3rem" }}
                    slotProps={{ textField: { error: !!ageError, fullWidth: true } }}
                />
            </LocalizationProvider>
            {ageError && err(ageError)}

            <input
                type="text"
                placeholder="Position *"
                className="AddToList-input"
                value={input.position}
                onChange={handleChange}
                onBlur={() => { if (!input.position.trim()) setPositionError("Position is required"); }}
                name="position"
            />
            {positionError && err(positionError)}

            <FormControl error={!!phoneError} fullWidth sx={{ mb: "0.3rem" }}>
                <TextField
                    label="Phone Number *"
                    fullWidth
                    value={localPhone}
                    onChange={handlePhoneChange}
                    onBlur={handlePhoneBlur}
                    error={!!phoneError}
                    placeholder="e.g. 70123456"
                    slotProps={{
                        input: {
                            startAdornment: (
                                <InputAdornment position="start">
                                    <Select
                                        value={COUNTRY.code}
                                        variant="standard"
                                        disableUnderline
                                        sx={{ minWidth: 95, mr: 0.5, fontSize: "0.95rem" }}
                                        renderValue={() => (
                                            <Box sx={{ display: "flex", alignItems: "center", gap: 0.5 }}>
                                                <span>{COUNTRY.flag}</span>
                                                <span>{COUNTRY.code}</span>
                                            </Box>
                                        )}
                                    >
                                        <MenuItem value={COUNTRY.code}>
                                            {COUNTRY.flag}&nbsp;{COUNTRY.label}&nbsp;({COUNTRY.code})
                                        </MenuItem>
                                    </Select>
                                </InputAdornment>
                            ),
                        },
                    }}
                />
                {phoneError && <FormHelperText>{phoneError}</FormHelperText>}
            </FormControl>

            <input
                type="text"
                placeholder="Department *"
                className="AddToList-input"
                value={input.department}
                onChange={handleChange}
                onBlur={() => { if (!input.department.trim()) setDepartmentError("Department is required"); }}
                name="department"
            />
            {departmentError && err(departmentError)}

            <input
                type="text"
                placeholder="Unit *"
                className="AddToList-input"
                value={input.unit}
                onChange={handleChange}
                onBlur={() => { if (!input.unit.trim()) setUnitError("Unit is required"); }}
                name="unit"
            />
            {unitError && err(unitError)}

            <div style={{ position: "relative" }}>
                <input
                    type={showNewPassword ? "text" : "password"}
                    placeholder="New Password *"
                    className="AddToList-input"
                    value={input.newPassword}
                    onChange={handleChange}
                    onBlur={handlePasswordBlur}
                    name="newPassword"
                    style={{ width: "100%" }}
                />
                <button
                    type="button"
                    onClick={() => setShowNewPassword(prev => !prev)}
                    style={{ position: "absolute", right: "8px", top: "50%", transform: "translateY(-50%)", background: "none", border: "none", cursor: "pointer", fontSize: "18px" }}
                >
                    {showNewPassword ? "🙈" : "👁️"}
                </button>
            </div>
            {passwordErrors.length > 0 && (
                <ul style={{ color: "red", fontSize: "13px", margin: "0 0 8px 0", paddingLeft: "1.2rem" }}>
                    {passwordErrors.map(e => <li key={e}>{e}</li>)}
                </ul>
            )}

            <div style={{ position: "relative" }}>
                <input
                    type={showConfirmPassword ? "text" : "password"}
                    placeholder="Confirm Password *"
                    className="AddToList-input"
                    value={input.confirmPassword}
                    onChange={handleChange}
                    onBlur={handleConfirmPasswordBlur}
                    name="confirmPassword"
                    style={{ width: "100%" }}
                />
                <button
                    type="button"
                    onClick={() => setShowConfirmPassword(prev => !prev)}
                    style={{ position: "absolute", right: "8px", top: "50%", transform: "translateY(-50%)", background: "none", border: "none", cursor: "pointer", fontSize: "18px" }}
                >
                    {showConfirmPassword ? "🙈" : "👁️"}
                </button>
            </div>
            {confirmPasswordError && err(confirmPasswordError)}

            {submitError && err(submitError)}
            <button className="AddToList-btn" onClick={handleClick}>Signup</button>
        </div>
    );
};

export default AddToList;

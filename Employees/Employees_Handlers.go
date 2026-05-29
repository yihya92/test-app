package Employees

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (Uc *UserControl) HTTP_Employees(w http.ResponseWriter, r *http.Request) {
	var sr API_Standard_response
	SourceIp, _ := GetRequestIP(r)
	sr.SourceIP = SourceIp
	sr.Login = r.Header.Get("Login")
	sr.SourceApp = r.Header.Get("SourceApp")
	sr.AccessKey = r.URL.Path
	sr.AccessMethod = r.Method
	sr.HostId = Configuration.HostId
	sr.ReceiveDate = time.Now()

	method := r.Method
	switch method {
	case "GET":
		sr.TransactionType = "Employees - Read"
		params := make(map[string]string)
		Login := r.URL.Query().Get("Login")
		if Login != "" {
			params["Login"] = Login
		}
		employees, err := Uc.Employee_Get(params)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to get employee(s)"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		sr.Data = employees

	case "POST":
		sr.TransactionType = "Employees - Create"
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var employee Employee
		err = json.Unmarshal(body, &employee)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		err = Uc.Employee_Add(sr.Login, employee)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to add employee"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		sr.Data = employee
	case "PUT":
		sr.TransactionType = "Employees - Update"
		vars := mux.Vars(r)
		LoginToEdit := vars["Login"]
		//parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to read request body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		var employee Employee
		err = json.Unmarshal(body, &employee)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to Unmarshal body"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		err = Uc.Employee_Edit(sr.Login, employee, LoginToEdit)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to edit employee"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, true)
			return
		}
		sr.Data = employee

			case "DELETE":
		sr.TransactionType = "Employee - Delete"
		vars := mux.Vars(r)
		LoginToDelete := vars["Login"]
		if LoginToDelete == "" {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": Login is empty"
			sr.ErrorDescription = "Login is empty"
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
		err := Uc.Employee_Delete(sr.Login, LoginToDelete)
		if err != nil {
			sr.Status = "failed"
			sr.StatusCode = http.StatusBadRequest
			sr.StatusDescription = http.StatusText(http.StatusBadRequest) + ": failed to delete employee"
			sr.ErrorDescription = err.Error()
			Uc.HTTP_API_Standard_response(w, r, sr, false)
			return
		}
	}

	sr.Status = "successful"
	sr.StatusCode = http.StatusOK
	sr.StatusDescription = http.StatusText(http.StatusOK)
	Uc.HTTP_API_Standard_response(w, r, sr, false)
}

func (Uc *UserControl) HTTP_API_Standard_response(w http.ResponseWriter, r *http.Request, transaction API_Standard_response, KeepDataInDB bool) {
	transaction.StatusDate = time.Now()
	transaction.Elapsedtime = (time.Since(transaction.ReceiveDate).Nanoseconds()) / 1000000
	Uc.Write_StandardResponse_log(transaction, "", KeepDataInDB)
	w.Header().Set("Content-Type", "application/json")
	if transaction.Status == "successful" {
		w.WriteHeader(transaction.StatusCode)
	} else if transaction.Status == "failed" {
		w.WriteHeader(transaction.StatusCode)
	}
	json.NewEncoder(w).Encode(transaction)
}


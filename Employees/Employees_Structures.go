package Employees

import "time"

type Employee struct {
	Key             string `bson:"Key" json:"Key"`
	Id              int64  `bson:"Id" json:"Id"`
	Login           string `bson:"Login" json:"Login"`
	Name            string `bson:"name" json:"name"`
	Email           string `bson:"email" json:"email"`
	Age             int    `bson:"age" json:"age"`
	Position        string `bson:"Position" json:"Position"`
	PhoneNumber     string `bson:"PhoneNumber" json:"PhoneNumber"`
	Department      string `bson:"Department" json:"Department"`
	Unit            string `bson:"Unit" json:"Unit"`
	NewPassword     string `bson:"NewPassword" json:"NewPassword"`
	ConfirmPassword string `bson:"ConfirmPassword" json:"ConfirmPassword"`
}
type API_Standard_response struct {
	//response source detail
	SourceIP        string    `bson:"SourceIP" json:"-"`
	Login           string    `bson:"Login" json:"Login"`
	SourceApp       string    `bson:"SourceApp" json:"-"`
	Language        string    `bson:"Language" json:"-"`
	AccessKey       string    `bson:"AccessKey" json:"-"`
	AccessMethod    string    `bson:"AccessMethod" json:"-"`
	HostId          string    `bson:"HostId" json:"-"`
	DeviceId        string    `bson:"DeviceId" json:"DeviceId"`
	DeviceName      string    `bson:"DeviceName" json:"DeviceName"`
	DeviceType      string    `bson:"DeviceType" json:"DeviceType"`
	OSType          string    `bson:"OSType" json:"OSType"`
	ReceiveDate     time.Time `bson:"ReceiveDate" json:"-"`
	TransactionType string    `bson:"TransactionType" json:"-"`
	TokenType       string    `bson:"TokenType" json:"TokenType"`
	//response detail
	Data interface{}
	//response result
	Status            string    `bson:"Status" json:"Status"` //successful, failed
	StatusCode        int       `bson:"StatusCode" json:"StatusCode"`
	StatusDescription string    `bson:"StatusDescription" json:"StatusDescription"` //error description if there is an error
	ErrorDescription  string    `bson:"ErrorDescription" json:"ErrorDescription"`
	StatusDate        time.Time `bson:"StatusDate" json:"-"`
	Elapsedtime       int64     `bson:"Elapsedtime" json:"-"`
}

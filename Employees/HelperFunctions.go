package Employees

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func GetRequestIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}
	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("no valid ip found")
}

func CheckMandatoryFields(st interface{}) (err error) {
	ptype := reflect.TypeOf(st)
	pvalue := reflect.ValueOf(st)

	ptr := reflect.New(ptype)
	temp := ptr.Elem()
	temp.Set(pvalue)
	for i := 0; i < ptype.NumField(); i++ {
		field := ptype.Field(i)
		//log.Println("field: ", field.Name)
		//check for mandatory fields
		if alias, ok := field.Tag.Lookup("V"); ok {
			if alias != "" {
				//log.Println("field: " + field.Name + " V: " + alias)
				existsMandatory := strings.Index(alias, "M")
				if existsMandatory == -1 {
					continue
				}
				switch field.Type.String() {
				case "string":
					fieldValue := temp.FieldByName(field.Name).String()
					if fieldValue == "" {
						err = errors.New(field.Name + " cannot be empty")
						return err
					}
					//log.Println(field.Name + ": " + alias + " - Value: " + fieldValue + " - type: " + field.Type.String())
				case "time.Time":
					fieldValue := pvalue.FieldByName(field.Name).Interface().(time.Time)
					if fieldValue.IsZero() {
						err = errors.New(field.Name + " cannot be empty")
						return err

					}
				case "int64":
					continue
				case "int32":
					continue
				case "int":
					continue
				case "float64":
					continue
				case "float32":
					continue
				case "bool":
					continue
				default:
					//log.Println(field.Type.String())
					ptrSub := reflect.New(field.Type)
					tempSub := ptrSub.Elem()
					pSubvalue := temp.FieldByName(field.Name)
					tempSub.Set(pSubvalue)
					err = CheckMandatoryFields(tempSub.Interface())
					if err != nil {
						return err
					}
				}

			}
		}
		// else {
		// 	log.Println("(not specified)")
		// }
	}
	return nil
}

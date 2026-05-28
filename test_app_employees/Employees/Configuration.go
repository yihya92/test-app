package Employees

import (
	"mongox"
	"redisx"
	"time"
)

var Configuration ConfigType

type ConfigType struct {
	HttpServicePort string
	HostId          string
	DB_Name         string
	Version         string
	Module          string

	IsProduction bool
	IsPrimary    bool //set to true for primary instance, for secondary should be false

	MongoDB struct {
		ReplicaSet string
		UserName   string
		Password   string
		HostIP_1   string
		HostPort_1 string
		HostIP_2   string
		HostPort_2 string
		HostIP_3   string
		HostPort_3 string
		HostIP_4   string
		HostPort_4 string
	}

	Redis redisx.Config
	Mongo mongox.Config
}

func GetDefaultConfiguration() (err error) {

	Configuration = setDefaultConfiguration_Test()

	return nil
}

func setDefaultConfiguration_Test() (Configuration ConfigType) {
	Configuration.HttpServicePort = "9900"
	Configuration.HostId = "EMP-01"
	Configuration.DB_Name = "Employees_DB"

	Configuration.Version = "V1"
	Configuration.Module = "EMP"

	Configuration.IsProduction = false
	Configuration.IsPrimary = true

	Configuration.Redis.Mode = redisx.ModeSingle
	Configuration.Redis.Addr = "127.0.0.1:6379"
	Configuration.Redis.Username = "admin"
	Configuration.Redis.Password = "ADMIN_PASSWORD"
	Configuration.Redis.DB = 0
	Configuration.Redis.DB = 0
	Configuration.Redis.KeyPrefix = "pdc:test:"
	Configuration.Redis.DefaultTTL = -1

	//New Mongodb config
	Configuration.Mongo.URI = "mongodb://localhost:9500"
	// --- Authentication ---
	Configuration.Mongo.Username = "db_root"
	Configuration.Mongo.Password = "dbrootpassword"
	Configuration.Mongo.AuthSource = "admin" // or your DB name if user is DB-scoped
	// --- Client identification ---
	Configuration.Mongo.AppName = "emp" //Shows in Mongo logs & profiler
	// --- Timeouts ---
	Configuration.Mongo.ConnectTimeout = 10 * time.Second         //Prevents startup hangs
	Configuration.Mongo.ServerSelectionTimeout = 10 * time.Second //Prevents blocking when primary unavailable
	Configuration.Mongo.SocketTimeout = 30 * time.Second          //Protects against stuck reads/writes
	// --- Pooling (tune per service load) ---
	Configuration.Mongo.MinPoolSize = 20 //Prevents overload under traffic
	Configuration.Mongo.MaxPoolSize = 400
	Configuration.Mongo.MaxConnIdleTime = 5 * time.Minute
	// --- Reliability ---
	Configuration.Mongo.RetryReads = true  //Handles transient network issues
	Configuration.Mongo.RetryWrites = true //Handles transient network issues
	// --- Topology ---
	Configuration.Mongo.Direct = false // true only for standalone when you KNOW it’s standalone

	//MongoDB SL Production
	// Configuration.MongoDB.UserName = "db_root"
	// Configuration.MongoDB.Password = "I@s54D0Grdara_r@23R"
	// Configuration.MongoDB.HostIP_1 = "Mega_db" //"host.docker.internal"
	// Configuration.MongoDB.HostPort_1 = "27017" //"27017"

	return
}

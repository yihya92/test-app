package Employees

import (
	"context"
	"errors"
	"log"
	"mongox"
	"redisx"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ---------------------------------------------------------------------
// Mongo collections helper (raw) – to reliably check Matched/DeletedCount
// ---------------------------------------------------------------------
func (Uc *UserControl) col(db, col string) *mongo.Collection {
	return Uc.MongoClient.Mongo.Database(db).Collection(col)
}

// ---------------------------------------------------------------------
// Initialize Mongox (kept), but writes below use raw collections for counts
// ---------------------------------------------------------------------
var (
	Mdb_AutoIncrement *mongox.Repository
	Mdb_Employees     *mongox.Repository
	Mdb_LogTrail      *mongox.Repository
)

func (UC *UserControl) InitializeMongoxRepositories() error {
	//Create mongox DB wrapper
	db, err := mongox.NewDB(UC.MongoClient.Mongo, Configuration.DB_Name, 5*time.Second)
	if err != nil {
		return err
	}
	//Create a repository bound to a collection
	if Mdb_Employees, err = mongox.NewRepository(db, "Col_Employees"); err != nil {
		return err
	}
	if Mdb_LogTrail, err = mongox.NewRepository(db, "Col_LogTrail"); err != nil {
		return err
	}
	return nil
}

func ensureIndex(repo *mongox.Repository, keys bson.D, opts *options.IndexOptionsBuilder) error {
	_, err := mongox.CreateIndex(
		context.Background(),
		repo.Coll,
		keys,
		opts,
	)
	return err
}

// ---------------------------------------------------------------------
// Redis keys
// ---------------------------------------------------------------------
func (d Employee) RedisKey() string { return "Employees:{" + d.Key + "}" }

// ---------------------------------------------------------------------
// Redis loader – FIXED: timeout + returns error (no log.Fatal)
// ---------------------------------------------------------------------

func (Uc *UserControl) RedisDataLoader() error {
	FlushBeforeLoad := Configuration.IsPrimary

	// A single timeout for the whole load operation (tune as needed)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Helper to reduce repetition
	load := func(name string, col *mongo.Collection, pattern string, fn func() (int64, int64, error)) error {
		loaded, total, err := fn()
		if err != nil {
			return err
		}
		log.Printf("%s Redis load: loaded=%d total=%d", name, loaded, total)
		return nil
	}
	// -------------------------
	// Employees
	// -------------------------
	{
		mcol := Mdb_Employees.Coll
		if err := load("Employees", mcol, "Employees:*", func() (int64, int64, error) {
			return redisx.LoadMongoToRedis[Employee](
				ctx,
				Uc.Redis,
				mcol,
				redisx.MongoLoadOptions{
					BatchSize:       2000,
					TTL:             0, // persistent
					FlushBeforeLoad: FlushBeforeLoad,
					FlushPattern:    "Employees:*",
					UseUnlink:       true,
				},
			)
		}); err != nil {
			return err
		}
	}

	return nil
}

// ---------------------------------------------------------------------
// Auto increment
// ---------------------------------------------------------------------

func (Uc *UserControl) GetNewId(Identifier string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return redisx.NextAutoIncrementID(
		ctx,
		Uc.Redis,
		Uc.col(Configuration.DB_Name, "Col_AutoIncrement"),
		Identifier,
		redisx.NextIDOptions{
			RedisBase:          "AutoIncrement",
			EmitReconcileEvent: true,
			ReconcileStream:    "AutoIncrement:reconcile",
			MongoRetries:       3,
			RetryBackoff:       500 * time.Millisecond,
		},
	)
}

func (Uc *UserControl) Write_StandardResponse_log(record API_Standard_response, Collection string, keepDatainDB bool) {
	if !keepDatainDB {
		record.Data = nil
	}
	Db := Mdb_LogTrail.Coll.Database().Name() + "_Logs"
	Col := Mdb_LogTrail.Coll.Name()
	if Collection != "" {
		Col = Collection
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := Uc.MongoClient.Mongo.Database(Db)
	collection := db.Collection(Col)
	_, err := mongox.InsertOne(ctx, collection, record)

	if err != nil {
		log.Println("Error in Write_StandardResponse_log:", err, " (", record, ")")
		return
	}
}

// Employee Functions //

func (Uc *UserControl) Employee_Get(filters map[string]string) (employees []Employee, err error) {
	var login string
	if len(filters) > 0 {
		login = filters["Login"]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if login != "" {
		employees_na := Employee{Key: login}
		employee, err := redisx.GetJSON[Employee](ctx, Uc.Redis, employees_na.RedisKey())
		if err != nil {
			if err == redis.Nil {
				err = errors.New("login does not exist")
				return employees, err
			}
			return employees, err
		}
		return []Employee{employee}, nil
	} else {
		employees, err := redisx.GetAllJSONByPattern[Employee](
			ctx,
			Uc.Redis,
			redisx.ScanJSONOptions{
				Pattern:      "Employees:*",
				ScanCount:    500,
				PipelineSize: 250,
				Limit:        10000,
			},
		)
		if err != nil {
			return nil, err
		}

		return employees, nil
	}
}

func (Uc *UserControl) Employee_Add(Login string, employee Employee) (err error) {
	//1- check for Mandatory fields
	err = CheckMandatoryFields(employee)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	existing_emp_i := Employee{Key: employee.Login}

	// 2- Check for duplication in Redis
	existing_emp, err := redisx.GetJSON[Employee](ctx, Uc.Redis, existing_emp_i.RedisKey())
	if err == nil && existing_emp.Id != 0 {
		employee.Id = existing_emp.Id
	} else if err != nil && err == redis.Nil {
		//3- fill auto increment fields
		newId, err := Uc.GetNewId("Employee-Id")
		if err != nil {
			return err
		}
		employee.Id = newId
	} else {
		return err
	}

	//4- fill default fields
	employee.Key = employee.Login
	existing_emp.Id = employee.Id
	existing_emp.Key = employee.Email
	existing_emp.Login = employee.Email
	existing_emp.Name = employee.Name
	existing_emp.Email = employee.Email
	existing_emp.Age = employee.Age
	existing_emp.Position = employee.Position
	existing_emp.PhoneNumber = employee.PhoneNumber
	existing_emp.Department = employee.Department
	existing_emp.Unit = employee.Unit
	existing_emp.NewPassword = employee.NewPassword
	existing_emp.ConfirmPassword = employee.ConfirmPassword

	// 5- Insert into MongoDB
	_, err = Mdb_Employees.InsertOne(ctx, existing_emp)
	if err != nil {
		if mongox.IsDuplicateKey(err) {
			return errors.New("already exists: " + err.Error())
		}
		return errors.New("insert failed: " + err.Error())
	}

	// 6- Populate Redis cache
	if err = redisx.SetJSON(ctx, Uc.Redis, existing_emp.RedisKey(), existing_emp); err != nil {
		return err
	}

	return nil
}

func (Uc *UserControl) Employee_Edit(Login string, employee Employee, LoginToEdit string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	employee_Na := Employee{Key: LoginToEdit}
	// 1- Load existing user from Redis
	employee_before, err := redisx.GetJSON[Employee](ctx, Uc.Redis, employee_Na.RedisKey())
	if err != nil {
		if err == redis.Nil {
			return errors.New("employee does not exit")
		}
		return err
	}
	//2- check if Id is filled
	if employee.Id != employee_before.Id {
		err = errors.New("employee Id does not match")
		return err
	}
	//3- preserve immutable fields; editable fields come from the request body
	employee.Key = employee_before.Key
	employee.Id = employee_before.Id
	employee.Login = employee_before.Login
	employee.NewPassword = employee_before.NewPassword
	employee.ConfirmPassword = employee_before.ConfirmPassword

	//4- check for Mandatory fields
	err = CheckMandatoryFields(employee)
	if err != nil {
		return err
	}
	// 5- Update MongoDB (upsert preserved)
	_, err = Mdb_Employees.UpdateOne(
		ctx,
		bson.M{"Key": employee.Key},
		bson.M{"$set": employee},
		options.UpdateOne().SetUpsert(false),
	)
	if err != nil {
		return errors.New("update failed: " + err.Error())
	}
	// 6- Update Redis cache
	if err = redisx.SetJSON(ctx, Uc.Redis, employee.RedisKey(), employee); err != nil {
		return err
	}

	return nil
}

func (Uc *UserControl) Employee_Delete(Login string, employeeLoginToDelete string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if employeeLoginToDelete != "" {
		employee_na := Employee{Key: employeeLoginToDelete}
		// load from redis
		employee, err := redisx.GetJSON[Employee](ctx, Uc.Redis, employee_na.RedisKey())
		if err != nil {
			if err == redis.Nil {
				return errors.New("login does not exist")
			}
			return err
		}
		// delete from MongoDB
		if _, err = Mdb_Employees.DeleteOne(ctx, bson.M{"Key": employeeLoginToDelete}); err != nil {
			return err
		}

		// delete from Redis
		if _, err = redisx.DelJSON(ctx, Uc.Redis, employee.RedisKey()); err != nil {
			return err
		}
		return nil

	} else {
		err = errors.New("login cannot be empty")
		return err
	}
}

// Go connection Sample Code:
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

var db *sql.DB
var server = "ec2736-db-server.database.windows.net"
var port = 1433
var user = "ec2736"
var password = "E@4JtDWBkepmCXS"
var database = "db"

func main() {
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")
	createDeleteTest(25, "25x1s_2xcreate_2xdelete.txt", time.Second)
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	_, err = CreateEmployee("Jake", "United States")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	// count, err := ReadEmployees()
	// 	// if err != nil {
	// 	// 	log.Fatal("Error reading Employees: ", err.Error())
	// 	// }
	// 	_, err = DeleteEmployee("Nikita")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Fprintf(w, "added and deleted successfully\n")
	// })
	// fmt.Println("Server is running on port 8080...")
	// http.ListenAndServe(":8080", nil)

	// fmt.Printf("Read %d row(s) successfully.\n", count)

}

func createDeleteTest(reps int, fileName string, pause time.Duration) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < reps; i++ {
		time.Sleep(pause)
		beforeTime := time.Now()
		_, err := CreateEmployee("Jake", "United States")
		if err != nil {
			log.Fatal(err)
		}
		_, err = CreateEmployee("Jake", "Germany")
		if err != nil {
			log.Fatal(err)
		}
		_, err = DeleteEmployee("Nikita")
		if err != nil {
			log.Fatal(err)
		}
		_, err = DeleteEmployee("Nikita")
		if err != nil {
			log.Fatal(err)
		}
		afterTime := time.Now()
		executionTime := afterTime.Sub(beforeTime)
		str := fmt.Sprintf("%v\n", executionTime)
		file.WriteString(str)
	}
}

func concurrentCreateDeleteTest(reps int, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	times := make(chan string, 300)
	for i := 0; i < reps; i++ {
		time.Sleep(3 * time.Second)
		for j := 0; j < 2; j++ {

			go func() {
				// time.Sleep(interval)
				beforeTime := time.Now()
				_, err := CreateEmployee("Jake", "United States")
				if err != nil {
					log.Fatal(err)
				}
				_, err = CreateEmployee("Jake", "Germany")
				if err != nil {
					log.Fatal(err)
				}
				_, err = DeleteEmployee("Nikita")
				if err != nil {
					log.Fatal(err)
				}
				_, err = DeleteEmployee("Nikita")
				if err != nil {
					log.Fatal(err)
				}
				afterTime := time.Now()
				executionTime := afterTime.Sub(beforeTime)
				str := fmt.Sprintf("%v\n", executionTime)
				times <- str
			}()
		}
		for i := 0; i < 2; i++ {
			str := <-times
			file.WriteString(str)
		}
	}

}

func ConnectToDB() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")
}

func CreateEmployee(name string, location string) (int64, error) {
	ctx := context.Background()
	var err error

	if db == nil {
		err = errors.New("CreateEmployee: db is null")
		return -1, err
	}

	// Check if database is alive.
	err = db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := `
      INSERT INTO TestSchema.Employees (Name, Location) VALUES (@Name, @Location);
      select isNull(SCOPE_IDENTITY(), -1);
    `

	stmt, err := db.Prepare(tsql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("Name", name),
		sql.Named("Location", location))
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, err
	}

	return newID, nil
}

func DeleteEmployee(name string) (int64, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := fmt.Sprintf("DELETE FROM TestSchema.Employees WHERE Name = @Name;")

	// Execute non-query with named parameters
	result, err := db.ExecContext(ctx, tsql, sql.Named("Name", name))
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

// ReadEmployees reads all employee records
func ReadEmployees() (int, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := fmt.Sprintf("SELECT Id, Name, Location FROM TestSchema.Employees;")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return -1, err
	}

	defer rows.Close()

	var count int

	// Iterate through the result set.
	for rows.Next() {
		var name, location string
		var id int

		// Get values from row.
		err := rows.Scan(&id, &name, &location)
		if err != nil {
			return -1, err
		}

		// fmt.Printf("ID: %d, Name: %s, Location: %s\n", id, name, location)
		count++
	}

	return count, nil
}

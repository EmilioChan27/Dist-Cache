// Go connection Sample Code:
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unsafe"

	c "github.com/EmilioChan27/Dist-Cache/common"
	_ "github.com/microsoft/go-mssqldb"
)

type DB struct {
	InnerDB *sql.DB
	Ctx     context.Context
}

func NewDB() *DB {
	var db *sql.DB
	var server = "ec2736-db-server.database.windows.net"
	var port = 1433
	var user = "ec2736"
	var password = "E@4JtDWBkepmCXS"
	var database = "db"
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
	return &DB{InnerDB: db, Ctx: context.Background()}
}

func getContentFromFile(file *os.File) []string {
	fileInfo, err := file.Stat()
	c.CheckErr(err)
	length := fileInfo.Size()
	byteOutput := make([]byte, length)
	_, err = file.Read(byteOutput)
	c.CheckErr(err)
	strOutput := string(byteOutput)
	strArrOutput := strings.Split(strOutput, "\\")
	return strArrOutput[:len(strArrOutput)-1]
}

func getArticleTitle(content string) string {
	titleUnTruncated := strings.Split(strings.Split(content, ":")[1], "\n")[0]
	title := titleUnTruncated[:len(titleUnTruncated)-4]
	return title
}
func getArticleCategory(content string) string {
	return strings.Split(content, ":")[0][4:]
}
func (db *DB) InsertTestArticles(filename string, maxAuthorId int) {
	file, err := os.Open(filename)
	c.CheckErr(err)
	contentArr := getContentFromFile(file)
	for i, content := range contentArr {
		title := getArticleTitle(content)
		category := getArticleCategory(content)
		a := &c.Article{Title: title, Content: content, ImageUrl: "Some random imageurl", Category: category, AuthorId: i % maxAuthorId}
		newId, err := db.AddArticle(a)
		c.CheckErr(err)
		fmt.Printf("Newid: %d\n", newId)
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

func (db *DB) checkDb() error {
	var err error
	if db == nil {
		err = errors.New("CreateEmployee: db is null")
		return err
	}
	err = db.InnerDB.PingContext(db.Ctx)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CreateAuthor(a *c.Author) (int64, error) {
	err := db.checkDb()
	c.CheckErr(err)
	tsql := `
      INSERT INTO IWSchema.Authors (Name, Bio, Email, ImageUrl) VALUES (@Name, @Bio, @Email, @Imageurl);
      select isNull(SCOPE_IDENTITY(), -1);
    `
	stmt, err := db.InnerDB.Prepare(tsql)
	c.CheckErr(err)
	defer stmt.Close()
	row := stmt.QueryRowContext(
		context.Background(),
		sql.Named("Name", a.Name),
		sql.Named("Bio", a.Bio),
		sql.Named("Email", a.Email),
		sql.Named("ImageUrl", a.ImageUrl),
	)
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, err
	}

	return newID, nil

}

func (db *DB) GetArticleById(id int) (*c.Article, error) {
	err := db.checkDb()
	c.CheckErr(err)
	tsql := fmt.Sprintf(`DECLARE @Id = %d
	SELECT * FROM IWSchema.Articles WHERE Id = @Id`, id)
	beforeTime := time.Now()
	row, err := db.InnerDB.QueryContext(db.Ctx, tsql)
	fmt.Printf("True DB exec Time: %v\n", time.Since(beforeTime))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer row.Close()
	// beforeTime = time.Now()
	var a *c.Article = &c.Article{}
	for row.Next() {
		err = row.Scan(&a.Id, &a.AuthorId, &a.Category, &a.Title, &a.ImageUrl, &a.Content, &a.CreatedAt, &a.UpdatedAt, &a.Likes, &a.Size)
		if err != nil {
			return nil, err
		}
	}
	// fmt.Printf("Scan time: %v\n", time.Since(beforeTime))
	return a, nil
}

func (db *DB) GetArticlesByCategory(category string, limit int) ([]*c.Article, error) {
	err := db.checkDb()
	c.CheckErr(err)
	limitStr := ""
	if limit != -1 {
		limitStr = fmt.Sprintf("TOP %d", limit)
	}
	pattern := "%" + category + "%"
	tsql := fmt.Sprintf(`DECLARE @CategoryPattern NVARCHAR(50) = '%s'
	SELECT %s* FROM IWSchema.Articles WHERE Category LIKE @CategoryPattern ORDER BY UpdatedAt DESC;`, pattern, limitStr)
	beforeTime := time.Now()
	rows, err := db.InnerDB.QueryContext(db.Ctx, tsql)
	fmt.Printf("True DB exec Time: %v\n", time.Since(beforeTime))
	articleList := make([]*c.Article, limit)
	if err != nil {
		fmt.Println(err)
		return articleList, err
	}
	defer rows.Close()
	var count int
	// beforeTime = time.Now()
	index := 0
	for rows.Next() {
		var a *c.Article = &c.Article{}
		err := rows.Scan(&a.Id, &a.AuthorId, &a.Category, &a.Title, &a.ImageUrl, &a.Content, &a.CreatedAt, &a.UpdatedAt, &a.Likes, &a.Size)
		if err != nil {
			return make([]*c.Article, 0), err
		}
		articleList[index] = a
		index++
		count++
	}
	// fmt.Printf("Scanning time: %v\n", time.Since(beforeTime))
	// fmt.Printf("Returning %d articles\n", count)
	return articleList[:index], nil
}

func (db *DB) AddArticle(a *c.Article) (int64, error) {
	var err error
	if db == nil {
		err = errors.New("AddArticle: db is null")
		return -1, err
	}
	err = db.InnerDB.PingContext(db.Ctx)
	if err != nil {
		return -1, err
	}
	tsql := `INSERT INTO IWSchema.Articles (AuthorId, Category, Title, ImageUrl, Content, CreatedAt, UpdatedAt, Likes, Size) VALUES(@AuthorId, @Category, @Title, @ImageUrl, @Content, @CreatedAt, @UpdatedAt, @Likes, @Size);
	      select isNull(SCOPE_IDENTITY(), -1);
	`
	stmt, err := db.InnerDB.Prepare(tsql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(
		db.Ctx,
		sql.Named("AuthorId", a.AuthorId),
		sql.Named("Category", a.Category),
		sql.Named("Title", a.Title),
		sql.Named("ImageUrl", a.ImageUrl),
		sql.Named("Content", a.Content),
		sql.Named("CreatedAt", time.Now()),
		sql.Named("UpdatedAt", time.Now()),
		sql.Named("Likes", 0),
		sql.Named("Size", int(unsafe.Sizeof(a))),
	)
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, err
	}
	return newID, nil
}
// func DeleteEmployee(name string) (int64, error) {
// 	ctx := context.Background()

// 	// Check if database is alive.
// 	err := db.PingContext(ctx)
// 	if err != nil {
// 		return -1, err
// 	}

// 	tsql := fmt.Sprintf("DELETE FROM TestSchema.Employees WHERE Name = @Name;")

// 	// Execute non-query with named parameters
// 	result, err := db.ExecContext(ctx, tsql, sql.Named("Name", name))
// 	if err != nil {
// 		return -1, err
// 	}

// 	return result.RowsAffected()
// }

// ReadEmployees reads all employee records
// func ReadEmployees() (int, error) {
// 	ctx := context.Background()

// 	// Check if database is alive.
// 	err := db.PingContext(ctx)
// 	if err != nil {
// 		return -1, err
// 	}

// 	tsql := fmt.Sprintf("SELECT Id, Name, Location FROM TestSchema.Employees;")

// 	// Execute query
// 	rows, err := db.QueryContext(ctx, tsql)
// 	if err != nil {
// 		return -1, err
// 	}

// 	defer rows.Close()

// 	var count int

// 	// Iterate through the result set.
// 	for rows.Next() {
// 		var name, location string
// 		var id int

// 		// Get values from row.
// 		err := rows.Scan(&id, &name, &location)
// 		if err != nil {
// 			return -1, err
// 		}

// 		// fmt.Printf("ID: %d, Name: %s, Location: %s\n", id, name, location)
// 		count++
// 	}

// 	return count, nil
// }

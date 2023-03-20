package cloud_sql

import (
	"compute_engine/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" //This is needed to import postgres driver.
)

type CloudSQL struct {
	db        *sql.DB
	tableName string
}

func (db *CloudSQL) init() error {
	_, err := db.db.Query(
		fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS %v (
                timestamp timestamp PRIMARY KEY default current_timestamp,
                content TEXT NOT NULL
            );
        `, db.tableName))

	return err
}

func (db *CloudSQL) Insert(content string) error {
	_, err := db.db.Exec(fmt.Sprintf("INSERT INTO %v (content) VALUES ($1)", db.tableName), content)
	return err

}

func (db *CloudSQL) selectByContent(content string) ([][]string, error) {
	rows, err := db.db.Query(fmt.Sprintf("SELECT * FROM %v WHERE content = $1", db.tableName), content)
	// rows, err := db.db.Query(fmt.Sprintf("SELECT * FROM %v", db.tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := [][]string{}
	for rows.Next() {
		var s1 string
		var s2 string
		rows.Scan(&s1, &s2)
		ret = append(ret, []string{s1, s2})
	}

	return ret, nil

}

func New(config config.PostgresConfig) (*CloudSQL, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=require", config.User, config.Password, config.Host, config.Port, config.DatabaseName)
	connection, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to connect to the database")
	}
	ret := &CloudSQL{db: connection, tableName: config.TableName}
	err = ret.init()
	if err != nil {
		return nil, err
	}
	return ret, nil

}

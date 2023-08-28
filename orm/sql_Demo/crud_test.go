package sql_Demo

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestDb(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	// 这里你就可可以用你的de了

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 处了Select语句都是ExecContext
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
	require.NoError(t, err)
	//	使用？作为查询的阐述的占位符
	res, err := db.ExecContext(ctx, `INSERT INTO test_model VALUES (?,?,?,?)`, 1, 2, 3, 4)
	require.NoError(t, err)

	res.RowsAffected()
	affected, err := res.RowsAffected()
	require.NoError(t, err)
	log.Println("受影响行数", affected)
	lastId, err := res.LastInsertId()
	require.NoError(t, err)
	log.Println("最后执行的", lastId)

	row := db.QueryRowContext(ctx,
		"SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 1)

	require.NoError(t, row.Err())
	tm := TestModelDB{}

	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)

	require.NoError(t, err)

	row = db.QueryRowContext(ctx,
		"SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 2)

	require.Error(t, sql.ErrNoRows, row.Err())
	tm = TestModelDB{}
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	//sql: no rows in result set
	fmt.Println(err)

	rows, err := db.QueryContext(ctx,
		"SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 2)

	require.Error(t, sql.ErrNoRows, rows.Err())
	for rows.Next() {
		tm = TestModelDB{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		require.NoError(t, err)
		log.Println(err)

	}
	tm = TestModelDB{}

	cancel()

}

type TestModelDB struct {
	Id int64
	// ""
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

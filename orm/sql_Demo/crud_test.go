package sql_Demo

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	// 这里你就可以用 DB 了
	// sql.OpenDB()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 除了 SELECT 语句，都是使用 ExecContext
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
	require.NoError(t, err)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)

	// 使用 ？ 作为查询的参数的占位符
	res, err := tx.ExecContext(ctx, "INSERT INTO test_model(`id`, `first_name`, `age`, `last_name`) VALUES(?, ?, ?, ?)",
		1, "Tom", 18, "Jerry")
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			log.Println(err)
		}
		return
	}
	require.NoError(t, err)
	affected, err := res.RowsAffected()
	require.NoError(t, err)
	log.Println("受影响行数", affected)
	lastId, err := res.LastInsertId()
	require.NoError(t, err)
	log.Println("最后插入 ID", lastId)

	// 提交事务
	err = tx.Commit()

	cancel()

}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

package sql_Demo

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestSQLMock(t *testing.T) {
	db, mock, err := sqlmock.New() // (*sql.DB, Sqlmock, error)
	defer db.Close()
	require.NoError(t, err)

	mockRows := sqlmock.NewRows([]string{"id", "first_name"})
	mockRows.AddRow(1, "Tom")
	// 正则表达式
	mock.ExpectQuery("SELECT id,first_name FROM `user`.*").WillReturnRows(mockRows)
	mock.ExpectQuery("SELECT id FROM `user`.*").WillReturnError(errors.New("mock err"))

	rows, err := db.QueryContext(context.Background(), "SELECT id,first_name FROM `user`WHERE id =1")
	fmt.Println(err)
	//require.NoError(t, err)

	for rows.Next() {
		tm := TestModelDB{}
		err = rows.Scan(&tm.Id, &tm.FirstName)
		require.NoError(t, err)
		log.Println(tm)
	}
	_, err = db.QueryContext(context.Background(), "SELECT id FROM `user`WHERE id=1")
	assert.Equal(t, errors.New("mock err"), err)
}

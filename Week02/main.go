package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func main() {
	fmt.Println("hello world")
}

var (
	ctx context.Context
	db  *sql.DB
)

// DAO层中，如果遇到 sqlErrNoRows 应该把返回的数据值写成空，然后根据需求选择是否Log这个Error
func dao() (string, error) {

	// copied from go example

	// In normal use, create one Stmt when your process starts.
	stmt, err := db.PrepareContext(ctx, "SELECT username FROM users WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Then reuse it each time you need to issue the query.
	id := 43
	var username string
	err = stmt.QueryRowContext(ctx, id).Scan(&username)
	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", err
	default:
		return username, nil
	}
}

func biz() error {
	userID, err := dao()
	if err != nil {
		return err
	}
	// 如果DAO层获取的数据是空，Service层根据业务可以选择返回404 Error，或者不处理。
	if userID == "" {
		return fmt.Errorf("404 error")
	}
	// biz logic
	// 处理UserID
	return nil
}

package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nudelfabrik/GoRadius"
)

type Database struct {
	sqlite *sql.DB
}

func NewDatabase() *Database {
	db := &Database{}
	err := db.init()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return db
}

func (db *Database) init() error {

	var err error
	db.sqlite, err = sql.Open("sqlite3", "/Users/bene/Downloads/freeradius.db")

	return err
}

func (db *Database) AddUser(user GoRadius.User) (err error) {
	tx, err := db.sqlite.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Adduser: ", r)
			err = errors.New("AddUser Failed")
		}
	}()

	execute(tx, "INSERT INTO radcheck (username,attribute,op,value) VALUES (?,'NT-Password',':=',?)", user.Name, user.PwHash)

	execute(tx, "INSERT INTO radusergroup (username,groupname) VALUES (?, ?)", user.Name, user.Name)

	execute(tx, "INSERT INTO radgroupreply (groupname,attribute,op,value) VALUES (?, 'Tunnel-Type', ':=', 13)", user.Name)
	execute(tx, "INSERT INTO radgroupreply (groupname,attribute,op,value) VALUES (?, 'Tunnel-Medium-Type', ':=', 6)", user.Name)
	execute(tx, "INSERT INTO radgroupreply (groupname,attribute,op,value) VALUES (?, 'Tunnel-Private-Group-Id', ':=', ?)", user.Name, user.VLAN)

	return tx.Commit()

}

func (db *Database) GetUser(name string) (user *GoRadius.User) {
	tx, err := db.sqlite.Begin()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("GetUser: ", r)
			user = nil
		}
	}()

	user = &GoRadius.User{}
	user.Name = name
	rows := query(tx, "SELECT value FROM radcheck WHERE attribute='NT-Password' AND username=?", name)
	if rows.Next() {
		err = rows.Scan(&user.PwHash)
		if err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("User does not exist"))
	}
	rows = query(tx, "SELECT value FROM radgroupreply WHERE attribute='Tunnel-Private-Group-Id' AND groupname=?", name)
	if rows.Next() {
		err = rows.Scan(&user.VLAN)
		if err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("User does not exist"))
	}

	return user

}

func (db *Database) DeleteUser(name string) error {
	tx, err := db.sqlite.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("DeleteUser: ", r)
			err = errors.New("DeleteUser Failed")
		}
	}()

	execute(tx, "DELETE FROM radcheck where username=?", name)
	execute(tx, "DELETE FROM radusergroup where username=?", name)
	execute(tx, "DELETE FROM radgroupreply where groupname=?", name)
	return tx.Commit()
}

func execute(tx *sql.Tx, statement string, values ...interface{}) {
	stmt, err := tx.Prepare(statement)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = stmt.Exec(values...)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
}

func query(tx *sql.Tx, statement string, values ...interface{}) *sql.Rows {
	stmt, err := tx.Prepare(statement)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	rows, err := stmt.Query(values...)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	return rows
}

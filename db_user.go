package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// NOTE user management

// Count users from database.
func sqlUserCount() (error, int) {
	var n int
	if err := db.QueryRow("SELECT COUNT(*) from user").Scan(&n); err != nil {
		log.Println(err)
		return err, 0
	}
	return nil, n
}

// Get a user identifyed by email and password.
func sqlAuthenticateUser(u *User) error {
	if err := db.QueryRow("SELECT * FROM user WHERE email=? AND password=?", u.Email, u.Password).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.isActive); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Get the user identifyed by email.
func sqlGetUser(u *User) error {
	if err := db.QueryRow("SELECT * FROM user WHERE email=?", u.Email).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.isActive); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Create new user.
func sqlInsert(u *User) error {
	if _, err := db.Exec("INSERT user SET name=?,email=?,password=?", &u.Name, &u.Email, &u.Password); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Update a user identifyed by email and password.
func sqlUpdateUser(u *User) error {
	if _, err := db.Exec("UPDATE user SET name=?, password=? WHERE email=?", u.Name, u.Password, u.Email); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Delete a user identifyed by email and password.
func sqlDeleteUser(u *User) error {
	if _, err := db.Exec("DELETE FROM user WHERE email=? AND password=?", u.Email, u.Password); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

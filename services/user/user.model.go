package user

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID         int       `json:"id"`
	First_Name string    `json:"first_name"`
	Last_Name  string    `json:"last_name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}

type UserModel struct {
	dB *sql.DB
}

func (u *UserModel) InsertOne(user *User) (int, error) {
	sqlStatement := `INSERT INTO users (first_name , last_name , email , password) VALUES ($1 ,$2 ,$3 ,$4) RETURNING id;`
	var insertId int
	err := u.dB.QueryRow(sqlStatement, user.First_Name, user.Last_Name, user.Email, user.Password).Scan(&insertId)

	if err != nil {
		return 0, err
	}
	return insertId, nil
}
func (u *UserModel) GetById(id int) (*User, error) {
	query := `SELECT * FROM users WHERE id=$1;`
	user := &User{}
	// fmt.Println(id)
	if err := u.dB.QueryRow(query, id).Scan(&user.ID, &user.First_Name, &user.Last_Name, &user.Email, &user.Password, &user.Created_At, &user.Updated_At); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no rows found where id : %v", id)
		}
		return nil, err
	}
	return user, nil
}
func (u *UserModel) GetAll() ([]*User, error) {
	sqlStmt := `SELECT * FROM users limit 20;`
	rows, err := u.dB.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []*User{}
	for rows.Next() {
		var user User = User{}
		rows.Scan(&user.ID, &user.First_Name, &user.Last_Name, &user.Email, &user.Password, &user.Created_At, &user.Updated_At)
		// fmt.Println(user)
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil

}
func (u *UserModel) UpdateOne(id int, user User) {}
func (u *UserModel) DeleteById(id int)           {}
func (u *UserModel) FindByEmail(email string) ([]*User, error) {
	query := `SELECT * FROM users WHERE email=$1`
	rows, err := u.dB.Query(query, email)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()
	users := []*User{}
	for rows.Next() {
		var user User = User{}
		rows.Scan(&user.ID, &user.First_Name, &user.Last_Name, &user.Email, &user.Password, &user.Created_At, &user.Updated_At)
		// fmt.Println(user)
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return users, nil

}

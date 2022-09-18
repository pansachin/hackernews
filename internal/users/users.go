package users

import (
	"database/sql"
	"log"

	database "github.com/pansachin/hackernews/internal/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"
)

type WrongUserOrPasswordError struct{}

func (*WrongUserOrPasswordError) Error() string {
	return "Wrong user name or password"
}

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (user *User) Create() {
	stmt, err := database.Db.Prepare("INSERT INTO Users(Username,Password) Values(?,?)")
	print(stmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	hashPwd, err := HashPwd(user.Password)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(user.UserName, hashPwd)
	if err != nil {
		log.Fatal(err)
	}

}

func HashPwd(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func ChekPwdHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func GetUserIDByUsername(username string) (int, error) {
	stmt, err := database.Db.Prepare("SELECT ID FROM Users where Username = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(username)
	var ID int
	err = row.Scan(&ID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return ID, nil
}

func (u *User) Authenticate() bool {
	stmt, err := database.Db.Prepare("SELECT Password FROM USERS WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res := stmt.QueryRow(u.UserName)

	var hashPwd string
	err = res.Scan(&hashPwd)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatal(err)
	}

	return ChekPwdHash(u.Password, hashPwd)
}

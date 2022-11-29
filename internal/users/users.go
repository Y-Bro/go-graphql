package users

import (
	"database/sql"
	"log"

	database "github.com/Y-bro/go-graphql/internal/pkg/db/migrations/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

func (user *User) Create() {
	stmt, err := database.Db.Prepare("INSERT INTO Users(username,password) VALUES(?,?)")

	if err != nil {
		log.Fatalln("Error Create user: ", err)
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		if err != nil {
			log.Fatalln("Error hash creation user: ", err)
		}
	}

	_, err = stmt.Exec(user.Username, hashedPassword)

	if err != nil {
		log.Fatalln(err)
	}
}

func (user *User) Authenticate() bool {
	stmt, err := database.Db.Prepare("SELECT Password FROM Users where Username = ?")

	if err != nil {
		log.Fatalln(err)
	}

	row := stmt.QueryRow(user.Username)

	var hashedPass string

	err = row.Scan(&hashedPass)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	return CheckPasswordHash(user.Password, hashedPass)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUserIdByUsername check if a user exists in database by given username
func GetUserIdByUsername(username string) (int, error) {
	statement, err := database.Db.Prepare("select ID from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return Id, nil
}

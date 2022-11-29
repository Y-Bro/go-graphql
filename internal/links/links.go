package links

import (
	"fmt"
	"log"

	database "github.com/Y-bro/go-graphql/internal/pkg/db/migrations/mysql"
	"github.com/Y-bro/go-graphql/internal/users"
)

type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

func (link Link) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title,Address, UserId) VALUES(?,?,?)")

	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(link.Title, link.Address, link.User.ID)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}

	log.Print("Row inserted!")
	return id
}

func GetAll() []Link {
	stmt, err := database.Db.Prepare("select L.id, L.title, L.address, L.UserID, U.Username from Links L inner join Users U on L.UserID = U.ID")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	links := []Link{}
	var username string
	var id string

	for rows.Next() {

		link := Link{}

		err := rows.Scan(&link.ID, &link.Title, &link.Address, &id, &username)
		if err != nil {
			log.Fatal(err)
		}

		link.User = &users.User{
			ID:       id,
			Username: username,
		}
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)

	}

	return links

}

func GetLink(id string) *Link {
	stmt, err := database.Db.Prepare("select L.id, L.title, L.address, L.UserID, U.Username from Links L inner join Users U on L.UserID = U.ID WHERE L.id = ?")

	if err != nil {
		log.Fatal(err)
	}

	link := &Link{}
	var username string
	var userId string

	err = stmt.QueryRow(id).Scan(&link.ID, &link.Title, &link.Address, &userId, &username)

	if err != nil {
		return nil
	}

	user := &users.User{
		ID:       userId,
		Username: username,
	}

	link.User = user

	fmt.Println("after exec")

	if err != nil {
		log.Fatal(err)
	}

	return link

}

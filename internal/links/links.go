package links

import (
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
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title,Address) VALUES(?,?)")

	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(link.Title, link.Address)
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
	stmt, err := database.Db.Prepare("Select id, title, address FROM Links")

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

	for rows.Next() {

		link := Link{}

		err := rows.Scan(&link.ID, &link.Title, &link.Address)
		if err != nil {
			log.Fatal(err)
		}
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)

	}

	return links

}

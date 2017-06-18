package models


import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

type (
	PostgresStore struct{
		session *sql.DB
	}
)


/*
	Postgres Methods

	http://www.alexedwards.net/blog/organising-database-access
	https://godoc.org/github.com/lib/pq
	d, _ := json.Marshal(p)
	Postgres().get(d)
*/
func Postgres() *PostgresStore {
	c := Config()
	n := c.session.Postgres

	dbinfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", n.User, n.Pass, n.Server, n.Port, n.Database)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}

	return &PostgresStore{db}
}

func (pg PostgresStore) get(data []byte) {

	db := pg.session
	defer db.Close()

	fmt.Println("# Querying")

	var uuid string
	var name string
	var first_name string
	var last_name string
	var role string
	var enabled string
	var locked interface{} // TODO: Change to bool

	stmt := `SELECT "UUID", name, first_name, last_name, role, enabled, locked FROM "user" WHERE id=$1;`
	row := db.QueryRow(stmt, 1)
	if err := row.Scan(&uuid, &name, &first_name, &last_name, &role, &enabled, &locked); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", name)

	//fmt.Printf("%+v\n", pg.session)

}

func (pg PostgresStore) getMany(data []byte) {

	db := pg.session
	defer db.Close()

	fmt.Println("# Querying")

	// TODO: UUID to lowercase
	stmt := `SELECT "UUID", name, first_name, last_name, role, enabled, locked FROM "user"`
	rows, err := db.Query(stmt)
	if err != nil {
		panic(err)
	}

	//orgs := make([]models.Organization,0)
	var names []string

	defer rows.Close()
	for rows.Next() {
		var uuid string
		var name string
		var first_name string
		var last_name string
		var role string
		var enabled string
		var locked interface{} // TODO: Change to bool

		//org := models.Organization{}
		//&org.uuid, &org.name, &org.first_name, &org.last_name, &org.role, &org.enabled, &org.locked
		//orgs = append(orgs, org)

		if err := rows.Scan(&uuid, &name, &first_name, &last_name, &role, &enabled, &locked); err != nil {
			fmt.Printf("%+v\n", err)
		}
		names = append(names, name)
	}

	fmt.Printf("%+v\n", names)

}


func (pg PostgresStore) create(data []byte) {
	db := pg.session
	defer db.Close()
}


func (pg PostgresStore) update(data []byte) {
	db := pg.session
	defer db.Close()
}


func (pg PostgresStore) remove(data []byte) {
	db := pg.session
	defer db.Close()
}
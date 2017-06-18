package models

import (
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"encoding/json"
	"fmt"
)

type (
	Neo4jStore struct{
		session bolt.Conn
	}
)


/*
	Neo4j Methods
	Bolt Protocol used for all Neo4j transactions

	https://github.com/johnnadratowski/golang-neo4j-bolt-driver
	https://godoc.org/github.com/johnnadratowski/golang-neo4j-bolt-driver

	d, _ := json.Marshal(p)
	Neo4j().get(d)
*/
func Neo4j() *Neo4jStore {

	c := Config()
	n := c.session.Neo4j

	dbinfo := fmt.Sprintf("bolt://%s:%s@%s:%s", n.User, n.Pass, n.Server, n.Port)
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(dbinfo)
	if err != nil {
		panic(err)
	}

	return &Neo4jStore{conn}
}


func (neo Neo4jStore) get(data []byte) {

}


func (neo Neo4jStore) create(data []byte) {

	conn := neo.session
	defer conn.Close()

	query := `CREATE (n:Organization {
		username:{username},
		age:{age},
		postalCode:{postalCode},
		longitude:{longitude},
		image:{image},
		phone:{phone},
		contact:{contact},
		streetAddress2:{streetAddress2},
		role:{role},
		firstName:{firstName},
		gender:{gender},
		bio:{bio},
		email:{email},
		country:{country},
		timestamp:{timestamp},
		streetAddress1:{streetAddress1},
		object:{object},
		uuid:{uuid},
		status:{status},
		lastName:{lastName},
		privacy:{privacy},
		city:{city},
		stateProvince:{stateProvince},
		displayAddress:{displayAddress},
		latitude:{latitude}
		})`

	// Here we prepare a new statement. This gives us the flexibility to
	// cancel that statement without any request sent to Neo
	//stmt, err := conn.PrepareNeo(query)
	//if err != nil {
	//	panic(err)
	//}

	var objmap map[string]interface{}
	err := json.Unmarshal(data, &objmap)
	fmt.Println(err)
	fmt.Println(objmap)

	// Executing a statement just returns summary information
	result, err := conn.ExecNeo(query, objmap)
	if err != nil {
		panic(err)
	}

	numResult, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1

	// Closing the statment will also close the rows
	//stmt.Close()

	//fmt.Printf("%+v\n", neo.session)

}


func (neo Neo4jStore) update(data []byte) {

}


func (neo Neo4jStore) remove(data []byte) {

}

package models

import (
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"encoding/json"
	"net/http"
	"fmt"
	"log"
)

type (

	MessagePayload struct {
		Auid 		string `json:"auid"`	// UUID of object making request (authorized UUID)
		Uuid 		string `json:"uuid"`	// UUID of object were referring to
		Object 		string `json:"object"` 	// object type were referring to. Usually mapped to the service
		Key 		string `json:"key"`		// key used for search
		Keyword 	string `json:"keyword"`	// value used for search
		Body 		string `json:"body"`	// Message Body key / value map {key:value, key:value}
		Perspective 	string `json:"perspective"`	// perspective were making the query from. ex: admin, superadmin, etc
		Results 	string `json:"results"`	// qty of objects returned
		Page 		string `json:"page"`	// page were results start
		Http_method 	string `json:"http_method"`  // GET, POST, PUT, DELETE - tell service what request method was used
		Method 		string `json:"method"` 	// Methods are the function that will process the request
		Version 	string `json:"version"` // Version of service requested
	}

	ResponseObject struct {
		Error 	bool `json:"error"`		// if error is present
		Message string `json:"message"`	// error message
		Code	int `json:"code"`		// error code, returned by api
		Stats 	map[string]string `json:"stats"`	// stats map
		Response	interface{} `json:"response"`	// main response, map or array
	}

	Organization struct {
		Object 		string `json:"object"`	// describes what object is being represented
		Uuid 		bson.ObjectId `json:"uuid"`
		Username	string `json:"username"`
		Name 		string `json:"name"`
		Bio			string `json:"bio"`
		Email 		string `json:"email"`
		Phone 		string `json:"phone"`
		Contact		string `json:"contact"`  // contact settings json string
		Privacy     string `json:"privacy"`  // privacy settings json string
		StreetAddress1 string	`json:"streetAddress1"`
		StreetAddress2 string	`json:"streetAddress2"`
		City		string `json:"city"`
		StateProvince	string `json:"stateProvince"`
		PostalCode	string `json:"postalCode"`
		Country		string `json:"country"`
		DisplayAddress	string `json:"displayAddress"`
		Longitude	string `json:"longitude"`
		Latitude	string `json:"latitude"`
		Image 		string `json:"image"`
		Timestamp	int `json:"timestamp"`	// time object was created
	}

	DataStore struct{
		session *mgo.Session
	}

	jsondata map[string]string

)

func ModelManager() *DataStore {

	c := Config()
	n := c.session.Mongo
	dbinfo := fmt.Sprintf("mongodb://%s:%s@%s", n.User, n.Pass, n.Server) // n.Port

	// connect to mongo server
	s, err := mgo.Dial(dbinfo)

	// Check if connection error, is mongo running?
	if err != nil {
		log.Println("Mongo connection error - ", dbinfo)
		panic(err)
	}

	return &DataStore{s}
}

/*
	Main Object Methods
*/
func (ds DataStore) GetObjectProfile(data MessagePayload) (ResponseObject) {

	s := ds.session.Copy()  // copy session
	defer s.Close()  // close session
	c := s.DB("csupplier").C("organization") // db : client, collection : object

	p := Organization{}
	id := data.Auid // default to first person profile query <I want to see my profile>
	if len(data.Uuid) != 0 {  // 3rd person query if uuid is provided <I want to see your profile>
		id = data.Uuid
	}

	// validate the uuid
	if !bson.IsObjectIdHex(id) {  // Invalid uuid
		return ResponseObject{ // Response Object
			Error: true,
			Message: "uuid missing or invalid",
			Code: http.StatusBadRequest, // update error
		}
	}

	// run query
	if err := c.Find(bson.M{"uuid": bson.ObjectIdHex(id)}).One(&p); err != nil {
		return ResponseObject{ // Response Object
			Error: true,
			Message: "no love.. sorry",
			Code: http.StatusNotFound,
		}
	}

	// Response Object
	return ResponseObject{
		Code: http.StatusOK,
		Response: p,
	}
}


func (ds *DataStore) GetObjects(data MessagePayload) (ResponseObject) {

	var organization []Organization
	var qry string

	s := ds.session.Copy()  // copy session
	defer s.Close()  // close session
	c := s.DB("csupplier").C("organization") // db : client, collection : object

	switch data.Key {

	case "all":  // all documents

		// run query
		if err := c.Find(nil).All(&organization); err != nil {
			return ResponseObject{ // Response Object
				Error: true,
				Message: "no love.. sorry",
				Code: http.StatusNotFound,
			}
		}

	case "email":  // email search
		qry = string(data.Keyword)
		if len(qry) == 0 {
			return ResponseObject{ // Response Object
				Error: true,
				Message: "missing key = keyword",
				Code: http.StatusBadRequest,
			}
		}

		// run query
		if err := c.Find(bson.M{"email": bson.M{"$regex": bson.RegEx{qry, "i"}}}).All(&organization); err != nil {
			return ResponseObject{ // Response Object
				Error: true,
				Message: "no love.. sorry",
				Code: http.StatusNotFound,
			}

		}

	default:  // firstname, lastname or username search
		qry = string(data.Keyword)
		if len(qry) == 0 {
			return ResponseObject{ // Response Object
				Error: true,
				Message: "missing key = keyword",
				Code: http.StatusBadRequest,
			}
		}

		// run query
		if err := c.Find( bson.M{ "$or": []bson.M{
			bson.M{"lastname": bson.M{"$regex": bson.RegEx{qry, "i"}}},
			bson.M{"firstname": bson.M{"$regex": bson.RegEx{qry, "i"}}},
			bson.M{"username": bson.M{"$regex": bson.RegEx{qry, "i"}}},
		}} ).All(&organization); err != nil {

			return ResponseObject{ // Response Object
				Error: true,
				Message: "no love.. sorry",
				Code: http.StatusNotFound,
			}
		}
	}

	if organization == nil {
		// Response Object
		return ResponseObject{
			Message: "nothing found",
			Code: http.StatusNotFound,
		}
	}

	// Response Object
	return ResponseObject{
		Code: http.StatusOK,
		Response: organization,
	}

}


func (ds *DataStore) CreateObject(data MessagePayload) (ResponseObject) {

	s := ds.session.Copy()
	defer s.Close()

	p := Organization{
		Uuid: bson.NewObjectId(),
		Object: "organization",
	}
	err := json.Unmarshal([]byte(data.Body), &p)
	if err != nil {
		return ResponseObject{ // Response Object
			Error: true,
			Message: "unmarshal error",
			Code: http.StatusBadRequest,
		}
	}

	c := s.DB("csupplier").C("organization")

	// run query
	if err := c.Insert(p); err != nil {
		return ResponseObject{ // Response Object
			Error: true,
			Message: "unique record conflict",
			Code: http.StatusConflict,
		}
	}

	// Response Object
	return ResponseObject{
		Code: http.StatusCreated,
		Response: p,
	}
}


func (ds *DataStore) CreateObjectFile(data MessagePayload) (ResponseObject) {

	// Response Object
	return ResponseObject{
		Code: http.StatusCreated,
		Response: Organization{
			Object: "organization",
			Name: "Test",
		},
	}
}


func (ds *DataStore) UpdateObjectProfile(data MessagePayload) (ResponseObject) {

	s := ds.session.Copy() // mgo session
	defer s.Close() // make sure we eventually close it
	c := s.DB("csupplier").C("organization") // object db & collection <client> <object>

	p := Organization{} // object struct
	id := data.Auid // default to first person profile query <I want to see my profile>
	if len(data.Uuid) != 0 {  // 3rd person query if uuid is provided <I want to see your profile>
		id = data.Uuid
	}

	// validate the uuid
	if !bson.IsObjectIdHex(id) {  // Invalid uuid
		return ResponseObject{ // Response Object
			Error: true,
			Message: "uuid missing or invalid",
			Code: http.StatusBadRequest, // update error
		}
	}

	var bm bson.M // BSON Map

	// unmarshal json string and map to bson map
	err := json.Unmarshal([]byte(data.Body), &bm)
	if err != nil {
		return ResponseObject{ // Response Object
			Error: true,
			Message: "unmarshal error",
			Code: http.StatusBadRequest, // unmarshal error
		}
	}

	// add bson map to mgo change struct
	change := mgo.Change{
		Update: bson.M{"$set": bm},
		ReturnNew: true,
	}

	// update object using bson map and change struct. update object struct with response
	// NOTE: if body KEY doesn't match object KEY, update will do nothing on the unmatched key but will update others
	if _, uerr := c.Find(bson.M{"uuid": bson.ObjectIdHex(id)}).Apply(change, &p); uerr != nil {
		return ResponseObject{ // Response Object
			Error: true,
			Message: "invalid update request",
			Code: http.StatusBadRequest, // update error - check for valid uuid
		}
	}

	// Response Object
	return ResponseObject{
		Code: http.StatusCreated,
		Response: p,
	}
}


func (ds *DataStore) UpdateObjects(data MessagePayload) (ResponseObject) {

	// Response Object
	return ResponseObject{
		Code: http.StatusOK,
		Response: Organization{
			Object: "organization",
			Name: "Test",
		},
	}
}


func (ds *DataStore) RemoveObject(data MessagePayload) (ResponseObject) {

	s := ds.session.Copy() // mgo session
	defer s.Close() // make sure we eventually close it
	c := s.DB("csupplier").C("organization") // object db & collection <client> <object>

	// TODO: ensure authorized to remove
	//auid := data.Auid // verify request is valid - person authorized to delete

	if len(data.Uuid) == 0 || !bson.IsObjectIdHex(data.Uuid) {  // Invalid uuid
		return ResponseObject{ // Response Object
			Error: true,
			Message: "uuid missing or invalid",
			Code: http.StatusBadRequest, // update error
		}
	}

	// Remove user
	if err := c.Remove(bson.M{"uuid": bson.ObjectIdHex(data.Uuid)}); err != nil {
		return ResponseObject{ // Response Object
			Error: true,
			Message: "uuid missing or invalid",
			Code: http.StatusNotFound, // not found error. id was valid but not found
		}
	}

	// Response Object
	return ResponseObject{
		Code: http.StatusNoContent, // 204
	}

}


func (ds *DataStore) RemoveObjectFile(data MessagePayload) (ResponseObject) {

	// Response Object
	return ResponseObject{
		Code: http.StatusNoContent,  // 204
	}
}


/*
	Redis Methods

	https://github.com/garyburd/redigo
	https://godoc.org/github.com/garyburd/redigo/redis
*/
func Redis(data []byte)  {

	c, err := redis.Dial("tcp", "#####.publb.lalalala.com:6379")
	if err != nil {
		panic(err)
	}

	n, err := c.Do("APPEND", "key", "value")
	fmt.Println(n)
	fmt.Println(err)

	c.Send("SET", "foo", "bar")
	c.Send("GET", "foo")
	c.Flush()
	c.Receive() // reply from SET
	v, err := c.Receive() // reply from GET
	fmt.Println(v)

	//////////

	response, err := c.Do("AUTH", "YOUR_PASSWORD")
	fmt.Println(response)

	if err != nil {
		panic(err)
	}

	//Set two keys
	c.Do("SET", "best_car_ever", "Tesla Model S")
	c.Do("SET", "worst_car_ever", "Geo Metro")

	//Get a key
	best_car_ever, err := redis.String(c.Do("GET", "best_car_ever"))
	if err != nil {
		fmt.Println("best_car_ever not found")
	} else {
		//Print our key if it exists
		fmt.Println("best_car_ever exists: " + best_car_ever)
	}

	//Delete a key
	c.Do("DEL", "worst_car_ever")

	//Try to retrieve the key we just deleted
	worst_car_ever, err := redis.String(c.Do("GET", "worst_car_ever"))
	if err != nil {
		fmt.Println("worst_car_ever not found", err)
	} else {
		//Print our key if it exists
		fmt.Println(worst_car_ever)
	}

	defer c.Close()

}
package controllers

import (
	"github.com/stevenmahana/OrganizationServiceTemplate/src/models"
	"encoding/json"
	"net/http"
	"log"
)

type (
	MainController struct{}
)

// NewController exposes all of the controller methods
func NewController() *MainController {
	return &MainController{}
}

// main controller
func (mc MainController) Controller(msg []byte) ([]byte) {

	// incoming message
	data := models.MessagePayload{}
	err := json.Unmarshal([]byte(msg), &data)
	if err != nil {
		log.Fatal(err) // TODO: handle this error
	}

	switch data.Http_method {

	case "GET":
		resp := mc.Get(data)
		r, e := json.Marshal(resp)
		if e != nil {
			log.Println(e) // TODO: handle this error
		}
		return r

	case "POST":
		resp := mc.Create(data)
		r, e := json.Marshal(resp)
		if e != nil {
			log.Println(e)
		}
		return r

	case "PUT":
		resp := mc.Update(data)
		r, e := json.Marshal(resp)
		if e != nil {
			log.Println(e)
		}
		return r

	case "DELETE":
		resp := mc.Remove(data)
		r, e := json.Marshal(resp)
		if e != nil {
			log.Println(e)
		}
		return r

	default:
		var s []byte
		return s

	}

}


/*
	GET Methods
 */
func (mc MainController) Get(data models.MessagePayload) (models.ResponseObject) {

	// access model manager
	m := models.ModelManager()

	switch data.Method {

	case "profile":
		resp := m.GetObjectProfile(data)  // response is map
		return resp

	case "search":
		resp := m.GetObjects(data) // response is slice of maps
		return resp

	default:
		// Response Object with error
		return models.ResponseObject{
			Error: true,
			Message: "did not recognize method",
			Code: http.StatusBadRequest,
		}

	}

}


/*
	POST Methods
 */
func (mc MainController) Create(data models.MessagePayload) (models.ResponseObject) {

	// access model manager
	m := models.ModelManager()

	switch data.Method {

	case "create":
		resp := m.CreateObject(data)  // response is map
		return resp

	case "image":
		resp := m.CreateObjectFile(data) // response is slice of maps
		return resp

	default:
		// Response Object with error
		return models.ResponseObject{
			Error: true,
			Message: "did not recognize method",
			Code: http.StatusBadRequest,
		}
	}

}


/*
	PUT Methods
 */
func (mc MainController) Update(data models.MessagePayload) (models.ResponseObject) {

	// access model manager
	m := models.ModelManager()

	switch data.Method {

	case "profile":
		resp := m.UpdateObjectProfile(data)  // response is map
		return resp

	case "bulk":
		resp := m.UpdateObjects(data) // response is map
		return resp

	default:
		// Response Object with error
		return models.ResponseObject{
			Error: true,
			Message: "did not recognize method",
			Code: http.StatusBadRequest,
		}
	}
}


/*
	DELETE Methods
 */
func (mc MainController) Remove(data models.MessagePayload) (models.ResponseObject) {

	// access model manager
	m := models.ModelManager()

	switch data.Method {

	case "remove":
		resp := m.RemoveObject(data)  // response is nil
		return resp

	case "file":
		resp := m.RemoveObjectFile(data) // response is nil
		return resp

	default:
		// response object with error
		return models.ResponseObject{
			Error: true,
			Message: "did not recognize method",
			Code: http.StatusBadRequest,
		}
	}

}
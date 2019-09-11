package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Company data and attributes
type Company struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Application string `json:"application,omitempty" bson:"application,omitempty"`
}

// CreateCompanyEndpoint creates a new company
func CreateCompanyEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var company Company
	_ = json.NewDecoder(request.Body).Decode(&company)
	collection := client.Database("RESTful").Collection("companies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, company)
	json.NewEncoder(response).Encode(result)
}

// GetCompaniesEndpoint retrieves all companies in the database
func GetCompaniesEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var companies []Company
	collection := client.Database("RESTful").Collection("companies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var company Company
		cursor.Decode(&company)
		companies = append(companies, company)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	}
	json.NewEncoder(response).Encode(companies)
}

// GetCompanyEndpoint retrieves a single existing company in the database
func GetCompanyEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	name, _ := params["name"]
	var company Company
	collection := client.Database("RESTful").Collection("companies")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, Company{Name: name}).Decode(&company)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(company)
}

// DeleteCompanyEndpoint deletes a company from the database
func DeleteCompanyEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	name, _ := params["name"]
	var company Company
	collection := client.Database("RESTful").Collection("companies")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOneAndDelete(ctx, Company{Name: name}).Decode(&company)
	if err == nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(company)
}

// // UpdateCompanyEndpoint updates a company application url in the database
// func UpdateCompanyEndpoint(response http.ResponseWriter, request *http.Request) {
// 	response.Header().Set("content-type", "application/json")
// 	params := mux.Vars(request)
// 	name, _ := params["name"]
// 	var company Company
// 	_ = json.NewDecoder(request.Body).Decode(&company)
// 	collection := client.Database("RESTful").Collection("companies")
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()
// 	filter := bson.D{{"name", name}}
// 	update := bson.D{{"$set", bson.D{{"application", company.Application}}}}
// 	err := collection.FindOneAndUpdate(
// 		ctx,
// 		filter,
// 		update).Decode(&company)
// 	if err == nil {
// 		response.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(response).Encode(company)
// }

// UpdateCompanyEndpoint updates a company application url in the database
func UpdateCompanyEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	name, _ := params["name"]
	var company Company
	_ = json.NewDecoder(request.Body).Decode(&company)
	filter := bson.D{{"name", name}}
	fmt.Println(name)
	update := bson.D{{"$set", bson.D{{"application", company.Application}}}}
	collection := client.Database("RESTful").Collection("companies")
	doc := collection.FindOneAndUpdate(
		context.Background(),
		filter,
		update,
		nil)
	fmt.Println(doc)
}

func main() {
	fmt.Println("Starting the application...")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/company", CreateCompanyEndpoint).Methods("POST")
	router.HandleFunc("/companies", GetCompaniesEndpoint).Methods("GET")
	router.HandleFunc("/company/{name}", GetCompanyEndpoint).Methods("GET")
	router.HandleFunc("/delete/{name}", DeleteCompanyEndpoint).Methods("DELETE")
	router.HandleFunc("/update/{name}", UpdateCompanyEndpoint).Methods("PUT")
	http.ListenAndServe(":8080", router)
}

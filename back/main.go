package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

var collec *mongo.Collection
var info Person
var emailExists bool
var mySigningKey = []byte("key")

type Person struct { //creates a person struct
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func init() {
	loadTheEnv()
	createDb()
}

var person []Person

func loadTheEnv() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func createDb() {
	connectionString := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	colName := os.Getenv("DB_COLLECTION_NAME")
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collec = client.Database(dbName).Collection(colName)

	fmt.Println("Collection instance created!")
}

//this function prints the people after the post function executes
//the method insertOne only executes when it gets the data in this method that's why it belongs in post
func getPeople(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, person) //this prints the array of items on the page in json format
}

//this function adds the request body into the slice
func post(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	fmt.Println("BODY ", string(body))
	if err != nil {
		panic(err)
	}
	defer c.Request.Body.Close()
	//puts the body in the json array values
	json.Unmarshal(body, &info)
	//adds info to the slice
	person = append(person, Person{Name: info.Name, Email: info.Email, Password: info.Password})
	for i := range person {
		fmt.Println("Array of people ", person[i])
	}
	fmt.Println("Person's data ", info)
	fmt.Println("Slice of people ", person)
	findEmail() //this method checks if an email exists
	if emailExists {
		fmt.Println("Email Exists")
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "email exists", // cast it to string before showing
		})
	} else {
		insertOneTask(info) //inserts data to mongodb on post request not get request
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Success", // cast it to string before showing
		})
	}
}

//this function retrieves all data from the database
func retrieveAllData(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collec.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var episodes []bson.M
	if err = cursor.All(ctx, &episodes); err != nil {
		log.Fatal(err)
	}
	fmt.Println("ALL DATA ", episodes)
	c.IndentedJSON(http.StatusOK, episodes)
}

//this function checks if users posted email exists in database
func findEmail() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collec.Find(ctx, bson.M{"email": info.Email}) //finds only the data with the email listed
	if err != nil {
		log.Fatal(err)
	}
	var episodes []bson.M
	if err = cursor.All(ctx, &episodes); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(episodes); i++ {
		fmt.Println("LOOPED DATA: ", episodes[i])
	}
	if episodes == nil {
		emailExists = false
	} else {
		emailExists = true
	}
	fmt.Println("Checking if email exists ", episodes)
}

func printEmailResponse(c *gin.Context) {
	if emailExists {
		c.IndentedJSON(http.StatusOK, "Email already exists")
		fmt.Println("Email Exists = ", emailExists)
		return
	} else {
		c.IndentedJSON(http.StatusOK, "")
		fmt.Println("You can use it")
	}
	n, _ := CreateToken(5) //the token is created here
	fmt.Println("JWT token ", n)
}

//THE ERROR WITH CHECKING IF EMAIL EXISTS: is the data has to post twice before not being allowed in the database
func main() {
	// Set the router as the default one shipped with Gin
	router := gin.Default()
	router.SetTrustedProxies([]string{"192.168.1.2"})
	router.Use(cors.Default())
	router.Use(static.Serve("/", static.LocalFile("./index.html", false)))
	// Setup route group for the API
	api := router.Group("/api")
	api.POST("/postform", post)

	//the below method gets the people from the post method and prints them
	api.GET("/", getPeople)
	api.GET("/data", retrieveAllData)
	api.GET("/email", printEmailResponse)

	router.Run(":8000")
}

func CreateToken(userid uint64) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(mySigningKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

//the method below accepts a type person to be inserted in the database
func insertOneTask(task Person) { //parameters accept the struct in the models file
	//inserts the type person in the database
	insertResult, err := collec.InsertOne(context.Background(), task)
	//below checks for error
	if err != nil {
		log.Fatal(err)
	}
	//else prints this line
	fmt.Println("Inserted a Single Record ", insertResult.InsertedID)
}

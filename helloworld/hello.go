//peoplesvc, a CRUD RestAPI Demonstrating GORM persistence
// curl -d '{"FirstName":"Gordon", "LastName":"Young", "City":"Gilbert"}' -H "Content-Type: application/json" -X POST http://localhost:8080/people

package main

import (
	"google.golang.org/appengine"
	"net/http"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

var db *gorm.DB
var err error

type Person struct {
	ID        uint   `json:”id”`
	FirstName string `json:”firstname”`
	LastName  string `json:”lastname”`
	City      string `json:”city”`
}

func main() {
	db, err = gorm.Open("sqlite3", "./people.db")

	defer db.Close()
	db.AutoMigrate(&Person{})

	//setup our controller
	r := gin.Default() 
	r.GET("/people /", GetPeople)
	r.GET("/people /:id", GetPerson)
	r.POST("/people", CreatePerson)
	r.PUT("/people /:id", UpdatePerson)
	r.DELETE("/people /:id", DeletePerson)

	http.HandleFunc("/", indexHandler)
	appengine.Main()
	//r.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}


func indexHandler(w http.ResponseWriter, req *http.Request) {
    fmt.Fprint(w, "done")
}

func DeletePerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person Person
	d := db.Where("id = ?", id).Delete(&person)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
func UpdatePerson(c *gin.Context) {
	var person Person
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&person).Error;
		err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&person)
	db.Save(&person)
	c.JSON(200, person)
}
func CreatePerson(c *gin.Context) {
	var person Person
	c.BindJSON(&person)
	db.Create(&person)
	c.JSON(200, person)
}
func GetPerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person Person
	if err := db.Where("id = ?", id).First(&person).Error;
		err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, person)
	}
}
func GetPeople(c *gin.Context) {
	var people []Person
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, people)
	}
}


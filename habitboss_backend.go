package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type Habit struct {
	Id            int
	IntervalType  int
	Description   string
	LastPerformed string
}

func exampleHabit() Habit {
	return Habit{Id: 12, IntervalType: 0, Description: "Walk the dog", LastPerformed: "2014-10-10T08:49:53+00:00"}
}

// ROOT
func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ROOT")
}

// API, basic routing
func currentHabits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	habits, _ := load()
	habitJson, _ := json.Marshal(habits)
	w.Write(habitJson)
	return
}

func createHabit(w http.ResponseWriter, r *http.Request) {
	habits, _ := load()

	// TODO extract params from request params
	id := rand.Intn(100000)
	intervalType := 0
	description := "Work out"
	lastPerformed := "2014-10-10T08:49:53+00:00"

	newHabit := Habit{Id: id, IntervalType: intervalType, Description: description, LastPerformed: lastPerformed}
	habits = append(habits, newHabit)

	err := persist(habits)
	if err == nil {
		fmt.Println("Successfully added: ", newHabit)
	} else {
		fmt.Println("Error while adding habit: ", newHabit)
	}
}

func updateHabit(w http.ResponseWriter, r *http.Request) {

}

func deleteHabit(w http.ResponseWriter, r *http.Request) {

}

/* TODO Routing to each endpoint

   Read: /currentHabits
   Create: /createHabit?intervalType=0&description="Do the laundry"&lastPerformed="2014-10-10T08:49:53+00:00"
   Update: /updateHabit?id=67&description="Do the laundry"&lastPerformed="2014-10-10T08:49:53+00:00"
   Delete: /deleteHabit?id=67

*/

/*
TODO gorilla mux:
prettier routing, match and preserve parts of url, match on methods:

/api/user/23/habit GET, PUT, POST, DELETE

*/

// TODO Add actual persistence behind each endpoint

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/webconsole", webconsole)

	http.HandleFunc("/api/currentHabits", currentHabits)
	http.HandleFunc("/api/createHabit", createHabit)
	http.HandleFunc("/api/updateHabit", updateHabit)
	http.HandleFunc("/api/deleteHabit", deleteHabit)

	http.ListenAndServe(":8080", nil)
}

// WEB
func webconsole(w http.ResponseWriter, r *http.Request) {
	habits, _ := load()

	t, err := template.ParseFiles("webconsole.html")
	if err != nil {
		showErrorPage(w, err)
		return
	}
	err = t.Execute(w, habits)
	if err != nil {
		showErrorPage(w, err)
	}
}

func showErrorPage(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// http://blog.golang.org/json-and-go
// One file per habit, or just save/load all of them as a habit array?

const persistedFilename string = "habits.json"

func persist(habits []Habit) error {
	filebody, err := json.Marshal(habits)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(persistedFilename, filebody, 0600)
}

func load() ([]Habit, error) {
	filebody, err := ioutil.ReadFile(persistedFilename)
	if err != nil {
		return []Habit{}, err
	}

	filebodyBytes := []byte(filebody)
	var habit []Habit

	err = json.Unmarshal(filebodyBytes, &habit)
	if err != nil {
		return []Habit{}, err
	}

	return habit, err
}

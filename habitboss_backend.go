package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

type Habit struct {
	Id            int
	IntervalType  string
	Description   string
	LastPerformed string
}

func exampleHabit() Habit {
	return Habit{Id: 12, IntervalType: "0", Description: "Walk the dog", LastPerformed: "2014-10-10T08:49:53+00:00"}
}

// ROOT
func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ROOT")
}

// API

// Example: http://localhost:8080/api/getHabits
func getHabits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	habits, _ := load()
	habitJson, _ := json.Marshal(habits)
	w.Write(habitJson)
	return
}

// Example: http://localhost:8080/api/createHabit?intervalType=%222%22&description=%22Take%20a%20swim%22&lastPerformed=%222014-10-10T08:49:53+00:00%22
func createHabit(w http.ResponseWriter, r *http.Request) {
	intervalType := r.URL.Query().Get("intervalType")
	description := r.URL.Query().Get("description")
	lastPerformed := r.URL.Query().Get("lastPerformed")

	newHabit := Habit{Id: rand.Intn(100000000), IntervalType: intervalType, Description: description, LastPerformed: lastPerformed}

	habits, _ := load()
	habits = append(habits, newHabit)
	err := persist(habits)

	if err == nil {
		fmt.Println("Successfully added: ", newHabit)
	} else {
		fmt.Println("Error while adding habit: ", newHabit)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func updateHabit(w http.ResponseWriter, r *http.Request) {
	idToUpdate, _ := strconv.Atoi(r.URL.Query().Get("id"))
	updatedDescription := r.URL.Query().Get("description")
	updatedLastPerformed := r.URL.Query().Get("lastPerformed")

	habits, _ := load()

	// Where in the habit list is it?
	indexOfDeletion := -1
	for index, element := range habits {
		if element.Id == idToUpdate {
			indexOfDeletion = index
		}
	}

	if indexOfDeletion == -1 {
		fmt.Println("No habit with id: ", idToUpdate)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update it
	oldHabit := habits[indexOfDeletion]
	newHabit := Habit{Id: oldHabit.Id, IntervalType: oldHabit.IntervalType, Description: updatedDescription, LastPerformed: updatedLastPerformed}
	oldHabit.Description = updatedDescription
	oldHabit.LastPerformed = updatedLastPerformed
	habits[indexOfDeletion] = newHabit

	err := persist(habits)
	if err == nil {
		fmt.Println("Successfully updated habit: ", idToUpdate)
	} else {
		fmt.Println("Error while updating habit: ", idToUpdate)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Example:    http://localhost:8080/api/deleteHabit?id=98081
func deleteHabit(w http.ResponseWriter, r *http.Request) {
	idToDelete, _ := strconv.Atoi(r.URL.Query().Get("id"))

	habits, _ := load()

	// Where in the habit list is it?
	indexOfDeletion := -1
	for index, element := range habits {
		if element.Id == idToDelete {
			indexOfDeletion = index
		}
	}

	if indexOfDeletion == -1 {
		fmt.Println("No habit with id: ", idToDelete)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Remove it. Not exactly pretty, but most ideomatic I could find
	habits = append(habits[:indexOfDeletion], habits[indexOfDeletion+1:]...)

	err := persist(habits)
	if err == nil {
		fmt.Println("Successfully removed habit: ", idToDelete)
	} else {
		fmt.Println("Error while removing habit: ", idToDelete)
		w.WriteHeader(http.StatusInternalServerError)
	}
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

	http.HandleFunc("/api/getHabits", getHabits)
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

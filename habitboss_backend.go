package main

import (
	"encoding/json"
	"errors"
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

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/webconsole", webconsole)
	http.HandleFunc("/api/getHabits", getHabits)
	http.HandleFunc("/api/createHabit", createHabit)
	http.HandleFunc("/api/updateHabit", updateHabit)
	http.HandleFunc("/api/deleteHabit", deleteHabit)
	http.ListenAndServe(":8080", nil)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ROOT")
}

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

// API ENDPOINTS

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

// Example: http://localhost:8080/api/updateHabit?id=27887&description=%22Clean%20the%20bathroom%22&lastPerformed=%222014-10-10T08:49:53+00:00%22
func updateHabit(w http.ResponseWriter, r *http.Request) {
	idToUpdate, _ := strconv.Atoi(r.URL.Query().Get("id"))
	updatedDescription := r.URL.Query().Get("description")
	updatedLastPerformed := r.URL.Query().Get("lastPerformed")

	habits, _ := load()

	i, findErr := indexOf(idToUpdate, habits)
	if findErr != nil {
		fmt.Println("No habit with id: ", idToUpdate)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update it
	oldHabit := habits[i]
	newHabit := Habit{Id: oldHabit.Id, IntervalType: oldHabit.IntervalType, Description: updatedDescription, LastPerformed: updatedLastPerformed}
	oldHabit.Description = updatedDescription
	oldHabit.LastPerformed = updatedLastPerformed
	habits[i] = newHabit

	err := persist(habits)
	if err == nil {
		fmt.Println("Successfully updated habit: ", idToUpdate)
	} else {
		fmt.Println("Error while updating habit: ", idToUpdate)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Example: http://localhost:8080/api/deleteHabit?id=98081
func deleteHabit(w http.ResponseWriter, r *http.Request) {
	idToDelete, _ := strconv.Atoi(r.URL.Query().Get("id"))

	habits, _ := load()

	i, findErr := indexOf(idToDelete, habits)
	if findErr != nil {
		fmt.Println("No habit with id: ", idToDelete)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Remove it. Not exactly pretty, but most ideomatic I could find
	habits = append(habits[:i], habits[i+1:]...)

	err := persist(habits)
	if err == nil {
		fmt.Println("Successfully removed habit: ", idToDelete)
	} else {
		fmt.Println("Error while removing habit: ", idToDelete)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// PERSISTENCE

func dbGetHabits() ([]Habit, error) {
	return load()
}

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

func indexOf(habitId int, habits []Habit) (int, error) {
	indexOfElement := -1
	for index, element := range habits {
		if element.Id == habitId {
			indexOfElement = index
		}
	}
	if indexOfElement == -1 {
		return indexOfElement, errors.New("failed to find habit in slice")
	} else {
		return indexOfElement, nil
	}
}

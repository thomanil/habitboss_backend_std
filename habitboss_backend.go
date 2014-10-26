package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Habit struct {
	Id            string
	IntervalType  int
	Description   string
	LastPerformed string
}

func main() {
	habit := exampleHabit()
	habitJson, _ := asJsonString(habit)
	fmt.Printf("%#v", habitJson)
}

/* 1. TODO Routing to each endpoint
// (cut out user/id stuff if too much work)

   Read: GET /api/allHabits

   Create: PUT /api/habit/?intervalType=0&description="Do the laundry"&lastPerformed="2014-10-10T08:49:53+00:00"

   Update: POST /api/habit/67&description="Do the laundry"&lastPerformed="2014-10-10T08:49:53+00:00"

   Delete: DELETE /api/habit/67

   (Return url in text for each one)
*/

// 2. TODO CRUD persistence operations

// http://blog.golang.org/json-and-go
// One file per habit, or just save/load all of them as a habit array?

func exampleHabit() Habit {
	return Habit{Id: "12", IntervalType: 0, Description: "Walk the dog", LastPerformed: "2014-10-10T08:49:53+00:00"}
}

const persistedFilename = "habits.json"

func asJsonString(habit Habit) (string, error) {
	jsonBytes, err := json.Marshal(habit)
	if err != nil {
		return "", err
	}
	jsonString := string(jsonBytes)
	return jsonString, err
}

func saveToFile(habit *Habit) error {
	filebody, err := json.Marshal(habit)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(persistedFilename, filebody, 0600)
}

func loadFromFile() (Habit, error) {
	filebody, err := ioutil.ReadFile(persistedFilename)
	if err != nil {
		return Habit{}, err
	}

	filebodyBytes := []byte(filebody)
	var habit Habit

	err = json.Unmarshal(filebodyBytes, &habit)
	if err != nil {
		return Habit{}, err
	}

	//TODO print first, see if it is intact
	return habit, err
}

// 3. TODO web console (template + image/css assets)

// 4. TODO Write brief summary of what net/http and friends gives us

// 5. TODO Deploy to a remote server

// Rerunning this on the fly:
// go get github.com/pilu/fresh
// run fresh in this dir to start process that restarts app on each change in go filename
// (turn off flycheck if it interfers here)

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Habit struct {
	Id            string
	IntervalType  int
	Description   string
	LastPerformed string
}

// ROOT
func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ROOT: %s!", r.URL.Path[1:])
}

// API
func showHabits(w http.ResponseWriter, r *http.Request) {
	// TODO avoid escaped quotes/str literals in json output?
	habit := exampleHabit()
	habitJson, _ := asJsonString(habit)
	fmt.Fprintf(w, "%#v", habitJson)
}

/* TODO Routing to each endpoint
   Read: GET /api/allHabits
   Create: PUT /api/habit/?intervalType=0&description="Do the laundry"&lastPerformed="2014-10-10T08:49:53+00:00"
   Update: POST /api/habit/67&description="Do the laundry"&lastPerformed="2014-10-10T08:49:53+00:00"
   Delete: DELETE /api/habit/67
*/

/* TODO Actual persistence behind each endpoint */

// WEB
func webconsole(w http.ResponseWriter, r *http.Request) {
	// TODO show prettified html representation of habits
	// TODO add image/css assets (/public dir? /template dir? Route urls to those assets?)
	habit := exampleHabit()
	t, err := template.ParseFiles("webconsole.html")
	if err != nil {
		showErrorPage(w, err)
		return
	}
	err = t.Execute(w, habit)
	if err != nil {
		showErrorPage(w, err)
	}
}

func showErrorPage(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/webconsole", webconsole)
	http.HandleFunc("/api/allHabits", showHabits)
	http.ListenAndServe(":8080", nil)
}

// http://blog.golang.org/json-and-go
// One file per habit, or just save/load all of them as a habit array?

func exampleHabit() Habit {
	return Habit{Id: "12", IntervalType: 0, Description: "Walk the dog", LastPerformed: "2014-10-10T08:49:53+00:00"}
}

const persistedFilename string = "habits.json"

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

package main

import (
	"fmt"
)

func main() {
	fmt.Printf(`
Usage: slapshot <RELEASE MESSAGE>
`)
}

/* 1. TODO Routing to each endpoint
// (cut out user/id stuff if too much work)

   Read: GET /api/user/23/allHabits

   Create: PUT /api/user/23?intervalType=0&description="Do the laundry"&lastPerformed="2014-10-10T08:49:53+00:00"

   Update: POST /api/user/23/habit/67&description="Do the laundry"&lastPerformed="2014-10-10T08:49:53+00:00"

   Delete: DELETE /api/user/23/habit/67

   (Return url in text for each one)
*/

// 2. TODO Model/json export of habit

// 3. TODO CRUD persistence operations

// 4. TODO web console (template + image/css assets)

// 5. TODO Write brief summary of what net/http and friends gives us

// 6. TODO Deploy to a remote server

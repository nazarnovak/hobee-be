package main

import (
	"fmt"
	"os"

	"github.com/nazarnovak/hobee-be/config"
	"github.com/nazarnovak/hobee-be/pkg/log"
	"github.com/nazarnovak/hobee-be/pkg/socket"
	"github.com/nazarnovak/hobee-be/pkg/email"
)

/*
TODO:

BIIIIITCH I HAS A CHAT!

0) Try to align the current code with the simplest code possible to make it run. Throw all the extra shit away or on the
backburner/issues in github
1) No need for extensive routing, that's all done in the FE already, a simple 404 maybe to an unknown path and then /api stuff

1) Rewrite tests and log to make it simpler! Bonus: write tests for log, errors!
Couple of functions:
* Write an error message - it should save the stacktrace where it was added, and pass it along to the final place where we log it
* Wrap values to existing error - add key-val pairs to the error so we can debug more easily
2) Include WS, and see how it works
3) Move Unit to a shared package? Probably something like User
4) Resolve naming for a User model which will be kept in DB and the User that is present in the websocket, with the
WS connection alive and paired channel to know when the user is paired
5) Check for memory leaks, shit's on fire, yo
6) Refactor the shit out of everything, since it's ugly! Is that important tho? Better to have something working than perfect?
7) Check why my simplistic router can be worse? Any benefit with the router we're currently using?
8) Make an init in Router? And then you can group similar requests, for example: routeTable ([]route <- make it this?) append
	"/ws" http.MethodGet, "/ws" http.MethodPost, etc
9) How to persist users if lights go out, do we lose the whole pool? Is that dangerous? Should we save the ID's of people
who are searching, any benefit of that?
10) Instead of using the existing mderrors way, why not just past an error that is an herror (struct) instead? Will that
change anything?

*/
func main() {
	c, err := config.Load()
	if err != nil {
		fmt.Printf("Config init fail: %s", err.Error())
		return
	}

	if err := log.Init(c.Log.Out); err != nil {
		fmt.Printf("Log init fail: %s", err.Error())
		return
	}

	if err := email.Init(c.Email.ApiKey, c.Email.Domain); err != nil {
		fmt.Printf("Email init fail: %s", err.Error())
		return
	}

	//if err := db.Init(c.DB.Connection); err != nil {
	//	fmt.Printf("DB init fail: %s", err)
	//	return
	//}

	usersPool := make(chan [2]*socket.User)

	go socket.Matcher(usersPool)

	socket.Rooms(usersPool)

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println(("$PORT must be set"))
		return
	}

	port = ":" + port

	fmt.Println("Running on port", port)

	s := NewServer()
	s.Start(c.Secret, port)

	//defer s.Stop()
}

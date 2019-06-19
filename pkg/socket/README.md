## Structure

1) When setting Cookie we set UUID (or if existing Cookie, we reuse UUID)
2) When creating a new socket -> add it to map of users -> map[UUID]User (User has UUID, Sockets []Socket). So append a new Socket here to the UUID
3) When searching - have [2]User that will be passed to the matcher
4) Matcher creates a room with [2]User's and adds the RoomID to each User. Now map[UUID]User will be able to tell in which room the user is in, so when you open the new tab with the same cookie like in step 2 - we just add the new socket to the User
5) When sending a message - we check if the socket (User) UUID is the same, and if it is - we send it as a "own" message, otherwise - buddy message

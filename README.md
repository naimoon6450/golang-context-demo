# TO RUN DEMO

In 1 terminal run the server
```bash
go run ./server/main.go
```
In another terminal make a request from the client
```bash
go run ./client/main.go
```

# Contexts

What problem are they trying to solve?
- Preventing unnecessary work from being done via proper cancellation propogation between related services / processes
- Request scoped state (ie. req id, trace id, auth tokens, etc.) 

What does a world without context look like?
- Context, imo, isn't intuitive if you haven't work with concurrent programming langs/libs
- A thread local store would help with access to globals such as db_connection or requests object and are accessible by every thread
- No goroutine local store and that's where contexts come in, so go routines can be aware of other routines with the same root context


relevant resources:
https://stackoverflow.com/questions/35036653/why-doesnt-this-golang-code-to-select-among-multiple-time-after-channels-work
https://rotational.io/blog/contexts-in-go-microservice-chains/
https://github.com/rotationalio/ctxms
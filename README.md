### README

#### Overview
This is a simple web framework named `waku`, which provides basic routing, middleware, and request/response handling capabilities. It is designed to be lightweight and easy to use, suitable for building small to medium-sized web applications.

#### Features
1. **Routing**: Support for GET and POST requests. Routes can be defined using patterns, including dynamic parameters.
2. **Middleware**: Middleware functions can be registered to handle cross - cutting concerns such as error recovery.
3. **Request Handling**: Provide methods to handle form data, JSON data, and query parameters.
4. **Response Handling**: Support for returning responses in various formats, including plain text, JSON, and HTML.

#### Code Structure
- **`waku.go`**: Contains the main `Engine` struct, which is the core of the framework. It provides methods for adding routes (`Get`, `Post`), starting the server (`Run`), and registering middleware (`Use`).
- **`router.go`**: Defines the `Router` struct, which is responsible for handling route matching and dispatching requests to the appropriate handlers.
- **`context.go`**: Defines the `Context` struct, which encapsulates the request and response information. It provides methods for handling requests (e.g., `PostForm`, `PostJSON`, `Query`) and sending responses (e.g., `String`, `JSON`, `HTML`).
- **`trie.go`**: Implements a trie data structure for efficient route matching.
- **`recover.go`**: Provides a middleware function for error recovery, which can catch panics and return an appropriate error response.

#### Usage

##### 1. Create a new engine
```go
engine := NewEngine()
```

##### 2. Add routes
```go
// Add a GET route
engine.Get("/hello", func(c *Context) {
    c.String(http.StatusOK, "Hello, World!")
})

// Add a POST route
engine.Post("/login", func(c *Context) {
    username := c.PostForm("username")
    password := c.PostForm("password")
    // Handle login logic
    c.JSON(http.StatusOK, H{"message": "Login successful"})
})
```

##### 3. Register middleware
```go
engine.Use(Recover())
```

##### 4. Start the server
```go
err := engine.Run(":8080")
if err != nil {
    log.Fatalf("Failed to start server: %v", err)
}
```

#### API Reference

##### Engine
- `NewEngine()`: Create a new engine instance.
- `Get(routePattern string, handleFunc handleFunc)`: Add a GET route.
- `Post(routePattern string, handleFunc handleFunc)`: Add a POST route.
- `Run(addr string)`: Start the server and listen on the specified address.
- `Use(middleware ...handleFunc)`: Register middleware functions.

##### Context
- `Param(key string)`: Get the value of a dynamic parameter in the route.
- `PostForm(key string)`: Get the value of a form field in a POST request.
- `PostJSON(dest T)`: Parse JSON data from a POST request.
- `Query(key string)`: Get the value of a query parameter in a GET request.
- `Status(code int)`: Set the HTTP status code of the response.
- `SetHeader(key string, value string)`: Set a header in the response.
- `String(code int, format string, values ...interface{})`: Return a plain text response.
- `JSON(code int, obj interface{})`: Return a JSON response.
- `HTML(code int, name string, data interface{})`: Return an HTML response.
- `Fail(code int, err string)`: Return an error response in JSON format.
- `Data(code int, data []byte)`: Return a binary data response.
- `Next()`: Call the next middleware or handler function.

##### Router
- `handle(c *Context)`: Handle a request by matching the route and calling the appropriate handler.

##### Recover
- `Recover()`: Return a middleware function for error recovery.

#### Conclusion
The `waku` framework provides a simple and flexible way to build web applications. It can be easily extended to meet more complex requirements by adding custom middleware and handlers.

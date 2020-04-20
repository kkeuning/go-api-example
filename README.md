# go-api-example
Learning example demonstrating an api server and client.

## Project status: 

A basic server example is provided using net/http and gorilla/mux.

A client is provided, or the server can be tested from curl or a web browser.  If testing from the browser, do not require an api key.

To avoid the need for a database, this example uses as its data a package variable in models.Users initialized with User records and methods similar to database operations.

The net/http ServeMux does not easily support path variables, so the example initially used a query parameter to specifiy the id of a user.  With the introduction of gorilla/mux, path variables are now supported.

Testing from the browser:
```
bin/server
```
1. http://localhost:8090 returns an html Hello, World! message, demonstrating the simplest way that static content can be served.
2. http://localhost:8090/api/v1/users returns a list of users in json format.
3. http://localhost:8090/api/v1/users?id=2 returns a single user with id=2 in json format.
4. http://localhost:8090/api/v1/users/2 returns a single user with id=2 in json format using the path variable.

Client usage, demonstrating api key:
```
bin/server -apikey=secret
bin/client list-users -apikey=secret
bin/client get-user -id=2 -apikey=secret
```

Logging from the gorilla/handlers package is in use.

Simple security middleware has been added to demonstrate both http middleware and api security.  While not production ready the example does provide a starting point for later adding security based on jwt or other mechanisms.

# Go Example WebApp

A simple go web app illustrating:

* Reading TOML config files
* Connecting to mySQL DB
* Simple authentication
* Gorilla mux
* Gorilla sessions
* Storing passwords with bcrypt
* Application Context rather than global variables

## Directions

1. ```git clone git@github.com:andrewaustin/go-example-webapp.git```
2. ```cd go-example-webapp```
3. ```cp config.example.toml config.toml```
4. Modify config.toml to suit your needs.
5. ```cat schemas/*.sql | mysql -u root -p -D <database name>```
6. ```go run *.go``` 

## Things I will eventually fix:

* Lock down routes
* REST endpoints
* Tests
* Use https://github.com/juju/errgo

## Thanks
Special thanks to Matt Silverlock (https://github.com/elithrar) for writing some excellent articles on web apps in go.

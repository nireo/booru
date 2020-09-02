# booru: A Classic Imageboard

## Setup

The project is built with golang, but you can create a custom client in any language you want, since booru also works as a rest api. By default booru will serve html which leverages golang templates.

```
# Installing dependencies
go get -u github.com/jinzhu/gorm

# Creating database (postgresql)
su - postgres
createdb dbname
```

### Configuration file

Most of the fields are quite self explanatory. But the `restApi` field means that the server won't serve html templates, but instead JSON. This means you can build your own client, but use booru as the back-end.

```json
// config.json
{
    "port": "8080",
    "adminAccess": true,
    "databaseHost": "localhost",
    "databasePort": "5432",
    "databaseUser": "postgres",
    "databaseName": "dbname",
    "adminLogin": {
        "username": "admin",
        "password": "password"
    },
    "restApi": true
}
```

### Running the code

```
go run main.go
```

## Contributions

Feel free to create a pull request if you want to change anything, we can go from there.

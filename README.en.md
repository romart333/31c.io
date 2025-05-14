# Test Task:
**Enhance the mini-server using Go, MySQL, GORM, Gin, and OpenAPI (Swagger).**


**The assignment consists of two parts:**:
- _Part 1_: Homework – At home, you need to familiarize yourself with the project structure and tests.
  It is also recommended to complete a series of simple tasks using the aforementioned technologies.
  Completing these recommended tasks will help you gain a deeper understanding of the project and
  simplify the Live Coding part (Part 2). This part of the assignment is described in this document.


- _Part 2_:  In real-time, you will be required to complete a simple task.
  This part of the assignment will be given during the in-person interview.

**Note:** We have prepared a server project template based on the Gin framework using GORM.
  In this template, the basic modules, controllers, and services are implemented.
  Additionally, 3 methods in the controllers are implemented and fully covered by integration tests.
  The integration tests will help you understand the purpose of the project.


**Note:** There is no need to run the server or connect to real database
  to complete the task. Testing should be carried out by running the integration tests.


**Description of the tasks for Part 1: Homework**

_1. Understand the project structure and tests._

1.1. Ensure you are using Go version 1.23.
```bash
    $ go version
```

1.2. Install the packages
```bash
    $ go install github.com/swaggo/swag/cmd/swag@latest
    $ go install github.com/cespare/reflex@latest
```

1.3. Install the dependencies
```bash
    $ go mod vendor
``` 

1.4. Generate the OpenAPI (Swagger) documentation
```bash
    $ swag init
``` 

1.5.Install/start the MySQL server (used in tests).

Approach 1:

- Install Docker: `https://docs.docker.com/get-docker/`;
- Run the script with admin privileges (a user with elevated rights): `docker start db-test || docker run -d -p 3406:3306 --env MYSQL_ROOT_PASSWORD=root_password --name db-test --rm mysql:latest`;
- In the `.env` file, do not provide the connection details for the test database (the default values correspond to the parameters of the MySQL container).;

Approach 2:

- Install the MySQL server: `https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/`;
- Create a MySQL user with administrative privileges;
- In the `.env` file, fill in the connection details for the test database.;

1.6. Create a `.env` file.
```bash
    APP_NAME={YOUR_APP_NAME} # Optional parameter, default value is `TestApp`
    APP_HOST={YOUR_APP_HOST} # Optional parameter, default value is `undefined`
    PORT={YOUR_APP_PORT} # Optional parameter, default value is `3000`
    IS_DEBUG={IS_DEBUG} # Optional parameter, default value is `true`
    ADMIN_X_API_KEY={ADMIN_X_API_KEY} # Required parameter; any non-empty string will do 
    CRON_X_API_KEY={CRON_X_API_KEY} # Required parameter; any non-empty string will do
    REQUEST_TIMEOUT_SEC={REQUEST_TIMEOUT_SEC} # Optional parameter, default value is `20`

    # Parameters for connecting to MySQL (required for running the application, not used in tests)
    DB_HOST={DB_HOST} # Optional parameter, default value is `localhost`
    DB_PORT={DB_PORT} # Optional parameter, default value is `3306`
    DB_USERNAME={DB_USERNAME} # Optional parameter, default value is `username`
    DB_PASSWORD={DB_PASSWORD} # Optional parameter, default value is `password`
    DB_SCHEMA={DB_SCHEMA} # Optional parameter, default value is `database`

    # Parameters for connecting to MySQL (required for running tests, not used in the application)
    TEST_DB_HOST={TEST_DB_HOST} # Optional parameter, default value is `localhost`
    TEST_DB_PORT={TEST_DB_PORT} # Optional parameter, default value is `3406`
    TEST_DB_USERNAME={TEST_DB_USERNAME} # Optional parameter, default value is `root`
    TEST_DB_PASSWORD={TEST_DB_PASSWORD} # Optional parameter, default value is `root_password`
    TEST_DB_SCHEMA={TEST_DB_SCHEMA} # Optional parameter, default value is `server`
``` 

1.7. Run the tests – if everything is configured correctly, all 18 tests should pass:
```bash
    $ go test -v
```


_2. Recommended tasks to complete._

_2.1 Enhance the entity: account._
Add 3 new fields to the `account` entity:
- `name` (varchar(255)),
- `rank` (tinyint, from 0 to 100),
- `memo` (text, nullable).

_2.2. Enhance the existing GET /account method – the method for retrieving a list of accounts with pagination.

2.2.1. Add 3 new fields to the dto `AccountDto`.

2.2.2. Add an optional parameter – the field search. Use this field to search for matches in the `address`, `name`, and `memo` fields of the `account` table.

2.2.3. Add the ability to sort not only by `id` or `updated_at`, but also by the fields `address`, `name` and `rank`.

2.2.4. Update the integration tests. Ensure that all the above changes are taken into account.

_2.3. Enhance the POST /account method – the method for creating a new account in the `account` table.

2.3.1. In the method’s `body`, add 3 new fields: the field `memo` should be optional, while the fields `name` and `rank` are required. Be sure to validate the new fields.

2.3.2. When creating a new account in the `account` table, take the new parameters from the `body` into account.

2.3.3. Update the integration tests. Ensure that all the above changes are taken into account.

_2.4. Adjust the OpenAPI (Swagger) documentation to reflect all the above changes._

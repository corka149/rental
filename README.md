# Rental

> A general purpose rental administration system.

Prerequisites
-------------
- Go 1.22
- Task 3

Installation
------------

This will install all the dependencies.

```bash
task setup
```

The configuration can be done by the following environment variables:
ADMIN_USER: The username for the admin user. Default is "admin".
ADMIN_PASSWORD: The password for the admin user. Default is "password".
DB_URL: The URL for the database. Default is "postgres://myadmin:mypassword@localhost:5432/rental_db".
PORT: The port for the application. Default is "8081".
SECRET: The secret for the application. Default is "secret".
MODE: The mode for the application. Default is "DEV" which logs and returns extra details.
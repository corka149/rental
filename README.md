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

## Possible integration as iframe

Example

```html
 <div style="display: flex; justify-content: center; align-items: center;">
    <iframe src="http://localhost:8081/calendar" width="1px" height="1px" id="calendar-iframe" loading="lazy">
    </iframe>
    <script>
        document.querySelector("iframe").onload = (event) => {
            const width = "248px";
            const height = "420px";

            event.target.width = width;
            event.target.height = height;
        };
    </script>
 </div>
```

# Runway

Runway is a simple, lightweight, auto-generated backend framework written in Go.
It enables rapid development by automatically generating and Admin Console
using [Gorm](https://gorm.io/)
Designed for flexibility and ease of use, Runway allows developers to quickly
scaffold robust backends with minimal configuration, while still supporting
custom logic and extensions. Its efficient, modular architecture makes it ideal
for building modern web applications, prototypes, and internal tools with
minimal overhead.

![Gopher above runway](assets/hero.png "Happy Take Off")

<!--toc:start-->

- [Runway](#runway)
  - [Usage](#usage)
  - [Contribution](#contribution)
    - [How to Contribute](#how-to-contribute)

## Usage

Install the tool and run your new admin project:

```sh
# Create a new directory for your project
mkdir myrunwayapp
cd myrunwayapp

# Initialize the project
go mod init myrunwayapp

# Get Runway
go get github.com/panzerit/runway

cat >main.go <<EOF
package main

import (
        "log"

        "github.com/panzerit/runway"
        "gorm.io/driver/sqlite"
        "gorm.io/gorm"
)

func main() {
        db, err := gorm.Open(sqlite.Open("db.sqlite3"), &gorm.Config{})
        if err != nil {
                log.Fatalf("failed connecting to sqlite3: %v", err)
        }

        app := runway.New("Runway Tester", "1234567890123456", db)

        // this would override the default homepage
        // app.Server.GET("/", func(c echo.Context) error {
        //      return c.String(200, "Hello from Runway!")
        // })

        app.Start()
}
EOF

# Clean up
go mod tidy

# Run the app
go run main.go
```

Now just [open your browser](http://localhost:1323).

## Contribution

I welcome contributions! Whether you want to fix bugs, improve documentation,
or suggest new features, your input is valuable. If you have an idea for a
new generator template, or want to help make the admin experience even more
myrunwayapp, open an issue or submit a pull request.

I believe that great admin tools are built by the community, for the community.
Let's make something myrunwayapp together!

### How to Contribute

- **Open Issues:** If you find a bug or have a feature request, please open an
  issue with clear details.
- **Pull Requests:** Fork the repo, create a branch, and submit a pull request.
  Please describe your changes and reference any related issues.
- **Code Style:** Follow idiomatic Go conventions. Keep changes focused and
  well-documented.
- **Respectful Communication:** Be kind and constructive in discussions and code
  reviews.
- **Tests:** Add or update tests where appropriate to ensure reliability.

All skill levels are welcome—whether you're new to open source or a seasoned
maintainer. Your ideas and feedback help shape the future of this project!

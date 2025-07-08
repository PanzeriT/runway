# Runway

Runway is a simple, lightweight, auto-generated backend framework written in Go.
It enables rapid development by automatically generating and Admin Console from
your [Ent](https://entgo.io) schema definitions. Designed for flexibility and
ease of use, Runway allows developers to quickly scaffold robust backends with
minimal configuration, while still supporting custom logic and extensions. Its
efficient, modular architecture makes it ideal for building modern web
applications, prototypes, and internal tools with minimal overhead.

<!--toc:start-->

- [Usage](#usage)
- [Contribution](#contribution)
  - [How to Contribute](#how-to-contribute)
  <!--toc:end-->

## Usage

Install the tool and run your new admin project:

```sh
# Create a new directory for your project
mkdir awesome
cd awesome

# Initialize the project using runway
./runway init awesome -f

# Generate the code
./runway generate

# Tidy up Go modules
go mod tidy

# Set your JWT secret
echo "JWT_SECRET=secret" > .env

# Run the application
go run main.go

# Open the admin interface in your browser
open localhost:1323
```

## Contribution

I welcome contributions! Whether you want to fix bugs, improve documentation,
or suggest new features, your input is valuable. If you have an idea for a
new generator template, or want to help make the admin experience even more
awesome, open an issue or submit a pull request.

I believe that great admin tools are built by the community, for the community.
Let's make something awesome together!

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

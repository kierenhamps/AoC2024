# Advent of Code 2024

My first attempt at Advent of Code!

## Goals

Instead of just hacking together a quick solution, which would be messy and unreadable, I want to use good cleanliness practices.

Therefore my goals for each solution will aim to be:

- Idiomatic Go: Using standard practices around golang
- Domain-Driven Design (DDD): to break down the domain of the solution into the domain language and data ownership
- Test Driven Development (TDD): to drive the development of the code using Red-Green-Refactor

So hopefully if anyone (including me) comes across this repo in the future they have a chance of understanding what I did to code up the domain, and how it works.

## Running

Each day is broken down into its own folder that can be run standalone.

To run the various challenges, first `cd` into the required day. Then simply run `go run main.go`.

## Testing

To run tests, you can either run `go test ./...` from the root directory. Alternatively, you can `cd` into the required day and run `go test .`

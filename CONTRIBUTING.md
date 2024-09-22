# Contributing to ZOG

Contributions are welcome, and they are greatly appreciated!

However, until we reach v1.0.0, since we are still in the early stages of the project:

- They will take longer to be merged
- You will have to do more ground work of reading the code since I won't write a detailed contributing guide yet.

## About the codebase

- Run `make test` to run all tests
- The zog package holds all the library code that is exported for use
- The zhttp package holds all the code that is used to parse http requests
- The zenv package holds all the code that is used to parse environment variables
- the primitives package holds Zog functions and types that are the building blocks of the library. This is mostly internal code that shouldn't be used directly outside of zog as it may break in the future. However I want to keep it exported so if someone wants to use it they can.

## Understanding the code

- I suggest you read the README.md section on the core design decisions & Zog Parsing Execution Structure

## Code Style

- **Test Names**: Test{USER_API_BEING_TESTED}{DESCRIPTION}
  - Specially important when testing the main ZogTypes, for example: `TestStringRequired`

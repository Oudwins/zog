# Contributing to ZOG

Contributions are welcome, and they are greatly appreciated!

## Initial steps

Before you start working on a contribution, create an issue describing what you want to build. It's possible someone else is already working on something similar, or perhaps there is a reason that feature isn't implemented. The maintainers will point you in the right direction.

<!-- ## Submitting a Pull Request

- Fork the repo
- Clone your forked repository: `git clone git@github.com:{your_username}/zog.git`
- Enter the zog directory: `cd zog`
- Create a new branch off the `master` branch: `git checkout -b your-feature-name`
- Implement your contributions (see the Development section for more information)
- Push your branch to the repo: `git push origin your-feature-name`
- Go to https://github.com/colinhacks/zog/compare and select the branch you just pushed in the "compare:" dropdown
- Submit the PR. The maintainers will follow up ASAP. -->

## Prerequisites

I encourage you to at least read through the zog.dev docs. Specially the [Core Concepts](https://zog.dev/category/core-concepts) and the [Core Design Decisions](https://zog.dev/core-design-decisions) sections.

### Setting up the development environment

The project really only depends on Golang and Nodejs (for the docs). So if you have those two installed, you should be good to go.

However, if you want the full development environment with some nice make commands to lint, test, etc. you can use Nix. Here is how you can do it:

1. Install Nix -> use this https://github.com/DeterminateSystems/nix-installer
2. Install direnv -> [guide](https://direnv.net/docs/installation.html)
3. Navigate to the project and run `direnv allow`. This will load the Nix environment into the project every time you navigate to it. Alternatively you can run `nix develop` to do the same manually.

## Documentation Only Contributions

All documentation is in the [docs](./docs) folder. It is built using [Docusaurus](https://docusaurus.io/). You can run `make docs-install` to install the dependencies and `make docs-dev` to run the docs locally.

## Code Contributions

### Architecture

- Packages
  - `zog` holds primary user facing API
  - `conf` holds user facing configuration API
  - `zconst` holds constant values used by the library. This is both internal and user facing.
  - `internals` holds code that should only be used internally by the library.
  - Optional Packages (should not be imported by ANY other package)
    - `zhttp` holds all the code that is used to parse http requests.
    - `zenv` holds all the code that is used to parse environment variables.
    - `i18n` holds all the code that is used to support internationalization.

### Coding Style

Most important thing is to please try to keep the code style consistent with the rest of the codebase. Here are some notes:

- **Test Names**: `Test{USER_API_BEING_TESTED}{DESCRIPTION}`
  - Specially important when testing the main ZogTypes, for example: `TestStringRequired`
- **Test File Names**: `{file_being_tested}_test.go`

### Tests

- All features should have tests, preferably with 100% coverage.

### Merging PRs

PRs will only be merged if:

- All CI checks pass
- At least one maintainer approves the PR

## Commands

**`make docs-install`**

- installs the dependencies for the docs

**`make docs-dev`**

- Runs docs in dev mode

**`make test`**

- runs all tests and generates pretty output (requires `gotestsum`)

**`make test-watch`**

- runs all tests in watch mode (requires `gotestsum`)

**`make test-cover`**

- runs all tests and generates a coverage report. It also shows any file that is not 100% covered.

**`make lint`**

- Runs golangcilint

## License

By contributing your code to the zog GitHub repository, you agree to
license your contribution under the MIT license.

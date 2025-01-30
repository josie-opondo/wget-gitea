# wgetApp

## Table of Contentp
- [Description](#description)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Testing](#testing)
- [License](#license)
- [Authors](#authors)

## Description
`wgetApp` is a Go-based application designed to download files from the web efficiently, leveraging concurrency and rate-limiting features. It provides a background download mechanism, logging capabilities, and URL processing to manage multiple downloads simultaneously.

## Features
- Download files in the background.
- Custom rate limiting for downloads.
- Concurrent processing with controlled parallel execution.
- Logging mechanism to track download status.
- Error handling for invalid URLs and failed downloads.

## Installation

1. **Clone the repository:**
- git clone https://learn.zone01kisumu.ke/git/josopondo/wget
- cd wgetApp

2. **Install dependencies:**
   ```sh
   go mod tidy
   ```

## Usage

Run the application:
```sh
   go run main.go -url="http://example.com/file.txt" -rate="1M"
```

## Testing

To run the test suite:
```sh
   go test ./...
```

## License
This project is licensed under the MIT License.

## Authors

- joseopondo
- hannapiko
- qochieng
- masman

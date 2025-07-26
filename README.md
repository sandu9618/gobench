# GoBench

**GoBench** is a simple HTTP benchmarking CLI tool written in Go.  
It lets you send multiple concurrent HTTP requests and get min/max/avg response times — perfect for load testing APIs.

## Features

- Support for different HTTP methods (GET, POST, PUT, DELETE, etc.)
- Request body support for POST/PUT requests
- Configurable content types (JSON, form data, etc.)
- Detailed error reporting with failure reasons
- Status code breakdown for failed requests
- Configurable concurrency levels
- Progress bar with real-time feedback
- Detailed statistics including min/max/average response times
- Requests per second calculation

## Installation

```bash
go install github.com/sandu9618/gobench@latest
```

## Usage

```bash
# Basic GET request benchmark
gobench -u https://example.com -n 100 -c 5

# POST request with JSON body
gobench -u https://api.example.com/users -m POST -d '{"name":"John","email":"john@example.com"}' -H "application/json" -n 50 -c 10

# POST request with form data
gobench -u https://api.example.com/users -m POST -d "name=John&email=john@example.com" -H "application/x-www-form-urlencoded" -n 30 -c 5

# PUT request with JSON body
gobench -u https://api.example.com/users/1 -m PUT -d '{"name":"Jane","email":"jane@example.com"}' -H "application/json" -n 30 -c 3

# DELETE request (no body needed)
gobench -u https://api.example.com/users/1 -m DELETE -n 20 -c 2
```

## Flags

- `-u, --url`: URL to benchmark (required)
- `-m, --method`: HTTP method to use (default: GET)
- `-d, --data`: Request body data (for POST/PUT requests)
- `-H, --header`: Content-Type header (defaults to application/json for POST/PUT with data)
- `-n, --requests`: Total number of requests (default: 1)
- `-c, --concurrency`: Number of concurrent workers (default: 1)

## Example Output

```
Sending requests 100% [########################################] (100/100)

====== GoBench Result ======
URL : https://api.example.com/users
Method : POST
Body : {"name":"John","email":"john@example.com"}
Content-Type : application/json
Total Requests : 100
Success : 95
Failed : 5

--- Failure Details ---
Status Codes:
  400: 2 requests
  500: 3 requests
Errors:
  HTTP 400: 400 Bad Request: 2 requests
  HTTP 500: 500 Internal Server Error: 3 requests

Min Time : 45.2ms
Max Time : 1.2s
Avg Time : 234.5ms
Total Duration: 23.45s
Requests/sec : 4.26
```

### Error Types

The tool now provides detailed error information:

- **HTTP Status Errors**: Shows specific status codes (400, 404, 500, etc.) with counts
- **Network Errors**: Connection failures, DNS resolution issues, timeouts
- **Request Errors**: Malformed URLs, invalid request creation
- **Response Errors**: Server errors, client errors, and their specific messages

## Content Types

The tool supports various content types for request bodies:

- **JSON**: `-H "application/json"` (default for POST/PUT with data)
- **Form Data**: `-H "application/x-www-form-urlencoded"`
- **XML**: `-H "application/xml"`
- **Plain Text**: `-H "text/plain"`
- **Custom**: Any other content type you specify

# Chapter 2: Foundations

## 2.1 Project setup and creating a module

- In Go, module paths should match a URL you own.

- For shareable projects, the module path should equal the code's download URL, e.g., `github.com/foo/bar` if hosted there.

## 2.2 Web application basics

```go
package main

import (
    "log"
    "net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello from Snippetbox"))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)

    log.Print("starting server on :4000")

    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
```

- `home` is a standard Go func with 2 params.

- `http.ListenAndServe()` always returns non-nil errors.
- Go's servemux treats `"/"` as a catch-all.

## 2.3 Routing requests

```go
package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a specific snippet..."))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a form for creating a new snippet..."))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/{$}", home)
    mux.HandleFunc("/snippet/view", snippetView)
    mux.HandleFunc("/snippet/create", snippetCreate)

    log.Print("starting server on :4000")

    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
```

**Route Patterns in Go**

- **Exact match**: No trailing slash (`/foo`) → matches only exact paths.
- **Subtree path**: Trailing slash (`/foo/`) → matches `/foo/**` (prefix match).
- **Disable subtree**: Append `{$}` (`/foo/{$}`) → forces exact match.
- **Path cleaning**: Redirects `.`, `..`, or `//` to clean URLs (e.g., `/foo/bar/..//baz` → `/foo/baz`).
- **Trailing slash redirect**: `/foo` → `/foo/` if `/foo/` is registered.
- **Host-specific routes**: Checked first (e.g., `foo.example.org/` before `/`).

**Default ServeMux**

```go
func main() {
    http.HandleFunc("/", home)
    http.HandleFunc("/snippet/view", snippetView)
    http.HandleFunc("/snippet/create", snippetCreate)

    log.Print("starting server on :4000")

    err := http.ListenAndServe(":4000", nil)
    log.Fatal(err)
}
```

- `http.Handle()` / `http.HandleFunc()` register routes in `http.DefaultServeMux`.
- Using `nil` in `ListenAndServe()` defaults to `DefaultServeMux`.
- **Best practice**: Avoid `DefaultServeMux` for clarity and security.

## 2.4 Wildcard route patterns

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)
    w.Write([]byte(msg))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a form for creating a new snippet..."))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/{$}", home)
    mux.HandleFunc("/snippet/view/{id}", snippetView)
    mux.HandleFunc("/snippet/create", snippetCreate)

    log.Print("starting server on :4000")

    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
```

**Wildcard Routing in Go**

- **Syntax**: `/{name}` (no patterns like `/c_{category}`, `/{y}-{m}-{d}`, or `/{slug}.html`).
- **Access values**: Use `r.PathValue("name")` (returns string; validate before use).

**Precedence Rules**

- Most specific path wins (`/post/edit` > `/post/{id}`).
- Conflicts (e.g., `/post/new/{id}` vs. `/post/{author}/latest`) cause runtime panic → avoid overlaps.

**Subtree Wildcards**

- `"/user/{id}/"` matches `/user/1/`, `/user/2/a/b`, etc.

**Remainder Wildcards**

- `/{path...}` captures all remaining segments (e.g., `/post/a/b/c` → `r.PathValue("path")` = `"a/b/c"`).

## 2.5 Method-based routing

**Method-Specific Routes**

- Prefix pattern with uppercase method (`GET /path`):
  ```go
  mux.HandleFunc("GET /{$}", home)
  mux.HandleFunc("POST /snippet/create", snippetCreatePost)
  ```
- `GET` matches `HEAD`; other methods require exact match.
- Multiple methods per path allowed:
  ```go
  mux.HandleFunc("GET /create", snippetCreate)
  mux.HandleFunc("POST /create", snippetCreatePost)
  ```

**Automatic Responses**

- Unsupported methods → `405 Method Not Allowed` (with `Allow` header).

**Precedence Rules**

- Method-specific routes (`POST /path`) override method-agnostic (`/path`).

**Limitations**

- **Not supported**:
  - Regex/advanced wildcards
  - Custom 404/405 pages
  - Multi-method routes
  - `OPTIONS` auto-handling
  - Header-based routing
- **Alternatives**: `chi`, `httprouter`, `gorilla/mux` ([comparison](https://www.alexedwards.net/blog/which-go-router-should-i-use)).

---

**Example Code**

**File: `main.go`**

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)
    w.Write([]byte(msg))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a form for creating a new snippet..."))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Save a new snippet..."))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /{$}", home)
    mux.HandleFunc("GET /snippet/view/{id}", snippetView)
    mux.HandleFunc("GET /snippet/create", snippetCreate)
    mux.HandleFunc("POST /snippet/create", snippetCreatePost)

    log.Print("starting server on :4000")

    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
```

## 2.6 Customizing responses

**Default Response Behavior**

- By default, responses have:

  - Status code: `200 OK`
  - Headers: `Date`, `Content-Length`, and `Content-Type` (sniffed).

- Example:

  ```bash
  $ curl -i localhost:4000/
  HTTP/1.1 200 OK
  Date: Wed, 18 Mar 2024 11:29:23 GMT
  Content-Length: 21
  Content-Type: text/plain; charset=utf-8

  Hello from Snippetbox
  ```

**Custom Status Codes**

```go
w.WriteHeader(http.StatusCreated)  // 201
// OR
w.WriteHeader(404)  // Must be called before Write()
```

- **Note**: `w.WriteHeader` can only be called once per response.

**Headers**

- **Set before writing**:

  ```go
  w.Header().Set("Content-Type", "application/json")
  w.Header().Add("Cache-Control", "public")
  ```

- **Methods**:

  - `Set()`: Overwrite
  - `Add()`: Append
  - `Del()`: Remove
  - `Get()`/`Values()`: Read

**Writing Response Bodies**

- For any functions that satisfy the `io.Writer` interface, you can pass in your `http.ResponseWriter` value.

- WHY?

  - `http.ResponseWriter` implements the `io.Writer` interface.
  - This means you can pass `http.ResponseWriter` to any function that expects an `io.Writer`.

- Example:

  ```go
  w.Write([]byte("Hello world"))

  io.WriteString(w, "Hello world")
  fmt.Fprint(w, "Hello world")
  ```

**JSON Responses**

- Go automatically sets the `Content-Type` header except for JSON.

- For JSON responses, manually set the `Content-Type` header:
  ```go
  w.Header().Set("Content-Type", "application/json")
  w.Write([]byte(`{"status":"ok"}`))
  ```

**Header Canonicalization**

- Names auto-capitalized (`cache-control` → `Cache-Control`)

- **Bypass**: Direct map access:

  ```go
  w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
  ```

- HTTP/2: Headers converted to lowercase

---

**Example Code**

**File: `main.go`**

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")

    w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a form for creating a new snippet..."))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusCreated)

    w.Write([]byte("Save a new snippet..."))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /{$}", home)
    mux.HandleFunc("GET /snippet/view/{id}", snippetView)
    mux.HandleFunc("GET /snippet/create", snippetCreate)
    mux.HandleFunc("POST /snippet/create", snippetCreatePost)

    log.Print("starting server on :4000")

    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
```

## 2.7 Project structure and organization

**Standard Layout**

```
project/
├── cmd/          # Entry points (e.g., `cmd/web/main.go`)
│   └── web/
│       ├── main.go
│       └── handlers.go
├── internal/     # Private reusable code (enforced import protection)
└── ui/           # UI assets
    ├── html/     # Templates
    └── static/   # CSS/JS
```

**Key Benefits**

- **Clean separation**: Go code (`cmd/`, `internal/`) vs. assets (`ui/`)

- **Scalable**: Add new executables (e.g., `cmd/cli`)
- **Safety**: `internal/` prevents external imports

**Usage**

```bash
go run ./cmd/web  # Run from project root
```

------

**Refactored Code**

**File: `main.go`**

```go
package main

import (
	"log"
	"net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /{$}", home)
    mux.HandleFunc("GET /snippet/view/{id}", snippetView)
    mux.HandleFunc("GET /snippet/create", snippetCreate)
    mux.HandleFunc("POST /snippet/create", snippetCreatePost)

    log.Print("starting server on :4000")

    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
```

**File: `handlers.go`**

```go
package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")
    w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a form for creating a new snippet..."))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Save a new snippet..."))
}
```

## 2.9 Serving static files

**Basic File Server Setup**

```go
fileServer := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
```

- `http.StripPrefix` removes `/static` from URL paths before file lookup
- Example: `/static/css/main.css` → serves `./ui/static/css/main.css`

**Key Features**

- **Security**

  - Automatic path sanitization prevents directory traversal attacks
  - For custom file serving, always sanitize paths:

  ```go
  http.ServeFile(w, r, filepath.Clean("./ui/static/"+userProvidedPath))
  ```

- **Performance Optimizations**

  - Supports HTTP Range requests (206 Partial Content)
  - Automatic `304 Not Modified` responses using `Last-Modified` headers
  - Content-Type detection from file extensions

- **Directory Listings**

  - Disable by adding empty `index.html` files
  - Advanced method: [Custom FileSystem implementation](https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings)

**Single File Serving**

```go
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./ui/static/report.pdf")
}
```

⚠️ **Warning**: Unlike `FileServer`, `ServeFile` doesn't automatically sanitize paths

**Example cURL Test**

```bash
# Test partial content support
curl -H "Range: bytes=0-100" http://localhost:4000/static/large-file.zip
```

---

**Example Code**

**File: `main.go`**

```go
package main

import (
	"log"
	"net/http"
)

func main() {
    mux := http.NewServeMux()

    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

    mux.HandleFunc("GET /{$}", home)
    mux.HandleFunc("GET /snippet/view/{id}", snippetView)
    mux.HandleFunc("GET /snippet/create", snippetCreate)
    mux.HandleFunc("POST /snippet/create", snippetCreatePost)

    log.Print("starting server on :4000")

    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
```

**File: `ui/html/base.tmpl`**

```html
{{define "base"}}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>{{template "title" .}} - Snippetbox</title>
    <link rel="stylesheet" href="/static/css/main.css" />
    <link
      rel="shortcut icon"
      href="/static/img/favicon.ico"
      type="image/x-icon"
    />
    <link
      rel="stylesheet"
      href="https://fonts.googleapis.com/css?family=Ubuntu+Mono:400,700"
    />
  </head>
  <body>
    <header>
      <h1><a href="/">Snippetbox</a></h1>
    </header>
    {{template "nav" .}}
    <main>{{template "main" .}}</main>
    <footer>Powered by <a href="https://golang.org/">Go</a></footer>
    <script src="/static/js/main.js" type="text/javascript"></script>
  </body>
</html>
{{end}}
```

## 2.10 The http.Handler interface

**Handlers**

- A handler satisfies `http.Handler`:
  ```go
  type Handler interface {
      ServeHTTP(ResponseWriter, *Request)
  }
  ```

- Struct example:
  ```go
  type home struct{}
  func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
      w.Write([]byte("This is my home page"))
  }
  mux.Handle("/", &home{})
  ```

**Handler Functions**

- Function example:
  ```go
  func home(w http.ResponseWriter, r *http.Request) {
      w.Write([]byte("This is my home page"))
  }
  ```

- Conversion methods:
  ```go
  mux.Handle("/", http.HandlerFunc(home))  // Adapter
  mux.HandleFunc("/", home)               // Shortcut
  ```

**Handler Chain**

- Servemux is a handler (implements `ServeHTTP`)
- Flow:
  1. Server calls servemux's `ServeHTTP()`
  2. Servemux routes to matching handler
  3. Handler's `ServeHTTP()` generates response

**Concurrency**

- Each request handled in separate goroutine
- Enables high efficiency
- Requires careful shared resource management

# Chapter 3: Configuration and error handling

## 3.1 Managing configuration settings

**Configuration Issues**

- Hard-coded settings in `main.go`:
  - Server address `":4000"`
  - Static files path `"./ui/static"`
- Problems:
  - No config/code separation
  - Can't change at runtime

**Command-Line Flags**

```go
addr := flag.String("addr", ":4000", "HTTP network address")
```
- Usage:
  ```bash
  go run ./cmd/web -addr=":80"
  ```
- Notes:
  - Ports 0-1023 need root
  - `-help` shows all flags

**Flag Types**

- Available functions:
  ```go
  flag.Int()
  flag.Bool()
  flag.Float64()
  flag.Duration()
  ```
- Failed conversions exit app

**Environment Variables**

```go
addr := os.Getenv("SNIPPETBOX_ADDR")
```
- Limitations:
  - No defaults
  - No type conversion
  - No help text
- Combined approach:
  ```bash
  export SNIPPETBOX_ADDR=":9999"
  go run ./cmd/web -addr=$SNIPPETBOX_ADDR
  ```

**Boolean Flags**

```bash
go run ./example -flag    # = true
go run ./example -flag=false
```

**Struct Configuration**

```go
type config struct {
    addr      string
    staticDir string
}

var cfg config
flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP address")
flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Static files path")
flag.Parse()
```

---

**Example Code**

**File: `main.go`**

```go
package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
    flag.Parse()

    mux := http.NewServeMux()

    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

    mux.HandleFunc("GET /{$}", home)
    mux.HandleFunc("GET /snippet/view/{id}", snippetView)
    mux.HandleFunc("GET /snippet/create", snippetCreate)
    mux.HandleFunc("POST /snippet/create", snippetCreatePost)

    log.Printf("starting server on %s", *addr)

    err := http.ListenAndServe(*addr, mux)
    log.Fatal(err)
}
```

## 3.2 Structured logging

**Structured Logging Overview**

- Replace `log` package with `slog`
- Features:
  - Timestamp (ms precision)
  - Severity levels (`Debug`, `Info`, `Warn`, `Error`)
  - Key-value attributes

**Logger Creation**

```go
// Text format
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

// JSON format
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
```
- Output examples:
  ```plaintext
  time=2024-03-18T11:29:23.000+00:00 level=INFO msg="starting server" addr=:4000
  ```
  ```json
  {"time":"...","level":"ERROR","msg":"address already in use"}
  ```

**Handler Options**

```go
&slog.HandlerOptions{
    Level: slog.LevelDebug,  // Log level threshold
    AddSource: true,         // Include file/line info
}
```
- Log levels: `Debug` < `Info` < `Warn` < `Error`

**Logging Methods**

```go
logger.Info("message", "key", value)          // Basic
logger.Error("message", slog.String("k", v))  // Type-safe
```
- Supported attribute helpers:
  `slog.String()`, `slog.Int()`, `slog.Bool()`, `slog.Time()`

**Implementation Example**

```go
func main() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    logger.Info("starting server", "addr", *addr)

    err := http.ListenAndServe(*addr, mux)
    logger.Error(err.Error())
    os.Exit(1)  // Manual exit replacement for log.Fatal()
}
```

**Log Management**

- Write to stdout for flexibility
- Redirect to file:
  ```bash
  go run ./cmd/web >> /tmp/web.log  # Append
  go run ./cmd/web > /tmp/web.log   # Overwrite
  ```

**Concurrency**

- `slog` loggers are thread-safe
- Shared destinations must have concurrency-safe `Write()`

---

**Example Code**

**File: `main.go`**

```go
package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
    flag.Parse()

    // Use the slog.New() function to initialize a new structured logger, which
    // writes to the standard out stream and uses the default settings.
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    mux := http.NewServeMux()

    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

    mux.HandleFunc("GET /{$}", home)
    mux.HandleFunc("GET /snippet/view/{id}", snippetView)
    mux.HandleFunc("GET /snippet/create", snippetCreate)
    mux.HandleFunc("POST /snippet/create", snippetCreatePost)

    // Use the Info() method to log the starting server message at Info severity
    // (along with the listen address as an attribute).
    logger.Info("starting server", "addr", *addr)

    err := http.ListenAndServe(*addr, mux)
    // And we also use the Error() method to log any error message returned by
    // http.ListenAndServe() at Error severity (with no additional attributes),
    // and then call os.Exit(1) to terminate the application with exit code 1.
    logger.Error(err.Error())
    os.Exit(1)
}
```

## 3.3 Dependency injection

**Problem & Solution**

- **Problem**: Handlers use `log.Print` instead of structured logger
- **Solution**: Dependency injection via application struct rather than global variables

**Implementation**

```go
// cmd/web/main.go
type application struct {
    logger *slog.Logger
}

func main() {
    app := &application{
        logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
    }
    mux.HandleFunc("GET /{$}", app.home) // Updated handler registration
}
```

```go
// cmd/web/handlers.go
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    // ... rest of handler
}

// All other handlers updated similarly:
// func (app *application) snippetView(...)
// func (app *application) snippetCreate(...)
```

**Alternative: Closure Approach**

```go
// package config
type Application struct {
    Logger *slog.Logger
}

// package foo
func ExampleHandler(app *config.Application) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err != nil {
            app.Logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
            // ...
        }
    }
}

// package main
mux.Handle("/", foo.ExampleHandler(app))
```

- [This](https://gist.github.com/alexedwards/5cd712192b4831058b21) is a more concrete example.

---

**Example Code**

**File: `cmd/web/main.go`**

```go
package main

import (
    "flag"
    "log/slog"
    "net/http"
    "os"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses.
type application struct {
    logger *slog.Logger
}

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    // Initialize a new instance of our application struct, containing the
    // dependencies (for now, just the structured logger).
    app := &application{
        logger: logger,
    }

    mux := http.NewServeMux()

    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

    // Swap the route declarations to use the application struct's methods as the
    // handler functions.
    mux.HandleFunc("GET /{$}", app.home)
    mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
    mux.HandleFunc("GET /snippet/create", app.snippetCreate)
    mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)


    logger.Info("starting server", "addr", *addr)

    err := http.ListenAndServe(*addr, mux)
    logger.Error(err.Error())
    os.Exit(1)
}
```

**File: `cmd/web/handlers.go`**

```go
package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strconv"
)

// Change the signature of the home handler so it is defined as a method against
// *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")

    files := []string{
        "./ui/html/base.tmpl",
        "./ui/html/partials/nav.tmpl",
        "./ui/html/pages/home.tmpl",
    }

    ts, err := template.ParseFiles(files...)
    if err != nil {
        // Because the home handler is now a method against the application
        // struct it can access its fields, including the structured logger. We'll
        // use this to create a log entry at Error level containing the error
        // message, also including the request method and URI as attributes to
        // assist with debugging.
        app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = ts.ExecuteTemplate(w, "base", nil)
    if err != nil {
        // And we also need to update the code here to use the structured logger
        // too.
        app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

// Change the signature of the snippetView handler so it is defined as a method
// against *application.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Change the signature of the snippetCreate handler so it is defined as a method
// against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a form for creating a new snippet..."))
}

// Change the signature of the snippetCreatePost handler so it is defined as a method
// against *application.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Save a new snippet..."))
}
```

## 3.4 Centralized error handling

**Key Notes:**

- `serverError` helper:
  - Logs error with method, URI and stack trace
  - Returns 500 response
- `clientError` helper:
  - Returns specific status code with standard text
- `http.StatusText()` provides standard status descriptions
- `debug.Stack()` captures stack trace for debugging

---

**Example Code**

**File: `cmd/web/helpers.go`**

```go
package main

import (
    "net/http"
)

// The serverError helper writes a log entry at Error level (including the request
// method and URI as attributes), then sends a generic 500 Internal Server Error
// response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
    var (
        method = r.Method
        uri    = r.URL.RequestURI()
        trace  = string(debug.Stack())
    )

    app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}
```

**File: `cmd/web/handlers.go`**

```go
package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")

    files := []string{
        "./ui/html/base.tmpl",
        "./ui/html/partials/nav.tmpl",
        "./ui/html/pages/home.tmpl",
    }

    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, r, err) // Use the serverError() helper.
        return
    }

    err = ts.ExecuteTemplate(w, "base", nil)
    if err != nil {
        app.serverError(w, r, err) // Use the serverError() helper.
    }
}

...
```

## 3.5 Isolating the application routes

**Improved Code Organization**

- Make `main()` only focused on:

  ```go
  // 1. Config parsing
  addr := flag.String("addr", ":4000", "HTTP network address")
  flag.Parse()

  // 2. Dependency setup
  app := &application{
      logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
  }

  // 3. Server execution
  err := http.ListenAndServe(*addr, app.routes())
  ```

**Key Notes:**

- `ListenAndServe` should accept a `*ServeMux` containing all routes

---

**Example Code**

**File: `cmd/web/main.go`**

```go
package main

...

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    app := &application{
        logger: logger,
    }

    logger.Info("starting server", "addr", *addr)

    // Call the new app.routes() method to get the servemux containing our routes,
    // and pass that to http.ListenAndServe().
    err := http.ListenAndServe(*addr, app.routes())
    logger.Error(err.Error())
    os.Exit(1)
}
```

**File: `cmd/web/routes.go`**

```go
package main

import "net/http"

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() *http.ServeMux {
    mux := http.NewServeMux()

    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

    mux.HandleFunc("GET /{$}", app.home)
    mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
    mux.HandleFunc("GET /snippet/create", app.snippetCreate)
    mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

    return mux
}
```

# Chapter 4: Database-driven responses

## 4.3 Modules and reproducible builds

**Go Modules Overview**

- Uses `go.mod` and `go.sum` for dependency management
- Ensures reproducible builds with exact versions

**The go.mod File**

```plaintext
module snippetbox.libra.dev
go 1.24.1

require (
  github.com/jackc/pgx/v5 v5.7.3 // indirect
  golang.org/x/crypto v0.31.0 // indirect
)
```
- Specifies module path, Go version, and dependencies
- `// indirect` marks transitive dependencies
- Projects can use different versions without conflicts

**The go.sum File**

```plaintext
github.com/jackc/pgx/v5 v5.7.3 h1:...
github.com/jackc/pgx/v5 v5.7.3/go.mod h1:...
```
- Contains cryptographic checksums
- Verifies downloaded package integrity

**Dependency Commands**

```bash
go mod verify    # Check package integrity
go mod download  # Download all dependencies
go mod tidy      # Clean unused dependencies
```

**Updating Dependencies**

```bash
go get -u github.com/foo/bar  # Upgrade to latest
go get foo/bar@v1.2.3         # Specific version
go get foo/bar@none           # Remove dependency
```

## 4.4 Creating a database connection pool

**Using `sql.Open()`**

```go
db, err := sql.Open("pgx", "postgres://user:pass@host:port/db?param1=value1&param2=value2")
if err != nil {
    ...
}
```

- **Parameters**:
  - Driver name (e.g., `"pgx"`)
  - Data Source Name (DSN) string
    - MySQL tip: Add `parseTime=true` for `time.Time` conversion

- **Key behavior**:
  - Creates a connection pool, not single connection
  - Thread-safe for concurrent use
  - Designed for long-lived usage (initialize once in `main()`)

**Implementation Example**

```go
import (
    "database/sql"
	...

    _ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
    dsn := flag.String("dsn", "postgres://...", "PostgreSQL DSN")
    db, err := openDB(*dsn)
    defer db.Close()
    ...
}

func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        db.Close()
        return nil, err
    }
    return db, nil
}
```

**Key Implementation Notes**:
- Use blank import for driver: `_ "github.com/jackc/pgx/v5/stdlib"`
- `Ping()` verifies working connection
- `defer db.Close()` for cleanup (good practice)

**Common Issues**:
- Missing dependencies may require:
  ```bash
  go mod tidy
  # or
  go get github.com/jackc/pgx/v5/pgxpool@v5.7.3
  ```

**Testing**:
```bash
go run ./cmd/web
# Should show server start message
```

---

**Example Code**

**File: `cmd/web/main.go`**

```go
package main

import (
    "database/sql" // New import
    "flag"
    "log/slog"
    "net/http"
    "os"

    _ "github.com/jackc/pgx/v5/stdlib" // New import
)

...

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
    // Define a new command-line flag for the postgresql DSN string.
		dsn := flag.String("dsn", "postgresql://postgres:cZYDwJ4zQK13ERGT@db.wjiqrmuhaespqbmerqhi.supabase.co:5432/postgres", "PostgreSQL data source name")
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    // To keep the main() function tidy I've put the code for creating a connection
    // pool into the separate openDB() function below. We pass openDB() the DSN
    // from the command-line flag.
    db, err := openDB(*dsn)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    // We also defer a call to db.Close(), so that the connection pool is closed
    // before the main() function exits.
    defer db.Close()

    app := &application{
        logger: logger,
    }

    logger.Info("starting server", "addr", *addr)

    // Because the err variable is now already declared in the code above, we need
    // to use the assignment operator = here, instead of the := 'declare and assign'
    // operator.
    err = http.ListenAndServe(*addr, app.routes())
    logger.Error(err.Error())
    os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, err
    }

    return db, nil
```

## 4.5 Designing a database model

**Database Model Structure**

```go
// internal/models/snippets.go
type Snippet struct {
    ID      int
    Title   string
    Content string
    Created time.Time
    Expires time.Time
}

type SnippetModel struct {
    DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {}
func (m *SnippetModel) Get(id int) (Snippet, error) {}
func (m *SnippetModel) Latest() ([]Snippet, error) {}
```

**Integration in Application**

```go
// cmd/web/main.go
type application struct {
    logger   *slog.Logger
    snippets *models.SnippetModel
}

func main() {
    db := openDB(*dsn)
    app := &application{
        logger: logger,
        snippets: &models.SnippetModel{DB: db},
    }
}
```

**Key Benefits**
- Clear separation between database and HTTP logic
- Self-contained database operations in `SnippetModel`
- Easy mocking for testing
- Runtime database configuration via DSN flag

**Implementation Notes**

- Model methods match database operations (CRUD)
- Application struct holds model instance
- Connection pool passed to model during initialization

---

**Example Code**

**File: `internal/models/snippets.go`**

```go
File: internal/models/snippets.go
package models

import (
    "database/sql"
    "time"
)

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our database snippets
// table.
type Snippet struct {
    ID      int
    Title   string
    Content string
    Created time.Time
    Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
    DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
    return 0, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
    return Snippet{}, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
    return nil, nil
}
```

**File: `cmd/web/main.go`**

```go
package main

import (
    "database/sql"
    "flag"
    "log/slog"
    "net/http"
    "os"


    _ "github.com/go-sql-driver/mysql"

    // Import the models package that we just created. You need to prefix this with
    // whatever module path you set up so that the import statement looks like this:
    // "{your-module-path}/internal/models".
    "snippetbox.libra.dev/internal/models"
)

// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
type application struct {
    logger   *slog.Logger
    snippets *models.SnippetModel
}

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
		dsn := flag.String("dsn", "postgresql://postgres:cZYDwJ4zQK13ERGT@db.wjiqrmuhaespqbmerqhi.supabase.co:5432/postgres", "PostgreSQL data source name")
		flag.Parse()

		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

		db, err := openDB(*dsn)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
		defer db.Close()

    // Initialize a models.SnippetModel instance containing the connection pool
    // and add it to the application dependencies.
    app := &application{
        logger:   logger,
        snippets: &models.SnippetModel{DB: db},
    }

    logger.Info("starting server", "addr", *addr)

    err = http.ListenAndServe(*addr, app.routes())
    logger.Error(err.Error())
    os.Exit(1)
}

...
```

## 4.6 Executing SQL statements

**Query Methods**

- `DB.Query()`: For `SELECT` returning multiple rows
- `DB.QueryRow()`: For `SELECT` returning single row
- `DB.Exec()`: For non-row-returning statements (`INSERT/UPDATE/DELETE`)

**Using DB.Exec()**

```go
result, err := db.Exec(query, args...)
// Or ignore result:
_, err := db.Exec(query, args...)
```
- Returns `sql.Result` with:
  - `LastInsertId()` (not supported in PostgreSQL)
  - `RowsAffected()`

**MySQL Implementation**
```go
// models/snippets.go
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
    stmt := `INSERT INTO snippets (title, content, created, expires)
             VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

    result, err := m.DB.Exec(stmt, title, content, expires)
    id, err := result.LastInsertId()
    return int(id), err
}

// handlers.go
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    id, err := app.snippets.Insert("Title", "Content", 7)
    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
```

**PostgreSQL Differences**
```go
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
    stmt := `INSERT INTO snippets (title, content, created, expires)
             VALUES($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 day' * $3)
             RETURNING id`

    var id int
    err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
    return id, err
}
```
- Uses `$1,$2,$3` placeholders
- `RETURNING` clause + `QueryRow().Scan()` to get ID
- Different timestamp handling

**Key Notes**
- Always use placeholder parameters (`?` or `$N`) to prevent SQL injection
- MySQL uses `?` placeholders, PostgreSQL uses `$1,$2,...`
- PostgreSQL doesn't support `LastInsertId()` - use `RETURNING` instead
- Test with curl:
  ```bash
  curl -iL -d "" http://localhost:4000/snippet/create
  ```

---

**Example Code**

**File: `internal/models/snippets.go`**

```go
package models

...

type SnippetModel struct {
    DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
    // Write the SQL statement we want to execute. I've split it over two lines
    // for readability (which is why it's surrounded with backquotes instead
    // of normal double quotes).
    stmt := `INSERT INTO snippets (title, content, created, expires)
        VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
    `

    // Use the Exec() method on the embedded connection pool to execute the
    // statement. The first parameter is the SQL statement, followed by the
    // values for the placeholder parameters: title, content and expiry in
    // that order. This method returns a sql.Result type, which contains some
    // basic information about what happened when the statement was executed.
    result, err := m.DB.Exec(stmt, title, content, expires)
    if err != nil {
        return 0, err
    }

    // Use the LastInsertId() method on the result to get the ID of our
    // newly inserted record in the snippets table.
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    // The ID returned has the type int64, so we convert it to an int type
    // before returning.
    return int(id), nil
}

...
```

**File: `cmd/web/handlers.go`**

```go
package main

...

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    // Create some variables holding dummy data. We'll remove these later on
    // during the build.
    title := "O snail"
    content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
    expires := 7

    // Pass the data to the SnippetModel.Insert() method, receiving the
    // ID of the new record back.
    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Redirect the user to the relevant page for the snippet.
    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
```

## 4.7 Single-record SQL queries

**Querying Single Rows**
```go
// models/snippets.go
func (m *SnippetModel) Get(id int) (Snippet, error) {
    stmt := `SELECT id, title, content, created, expires
             FROM snippets
             WHERE expires > UTC_TIMESTAMP() AND id = ?`

    var s Snippet
    err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

    if errors.Is(err, sql.ErrNoRows) {
        return Snippet{}, ErrNoRecord
    }
    return s, err
}

// models/errors.go
var ErrNoRecord = errors.New("models: no matching record found")

// handlers.go
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.PathValue("id"))
    snippet, err := app.snippets.Get(id)

    if errors.Is(err, models.ErrNoRecord) {
        http.NotFound(w, r)
        return
    }
    fmt.Fprintf(w, "%+v", snippet)
}
```

**Key Points**
- Use `QueryRow()` for single-row queries
- `Scan()` copies column values to struct fields
- Handle `sql.ErrNoRows` with custom error
- MySQL requires `parseTime=true` for time conversions

**PostgreSQL Differences**

```go
stmt := `SELECT ... WHERE expires > CURRENT_TIMESTAMP AND id = $1`
```

**Error Handling**

- Prefer `errors.Is()` over `==` for error checking because it handles wrapped errors
- Encapsulate the model and avoid exposing datastore-specific errors (e.g., `sql.ErrNoRows`) to handlers

---

**Example Code**

**File: `internal/models/snippets.go`**

```go
package models

import (
    "database/sql"
    "errors" // New import
    "time"
)

...

func (m *SnippetModel) Get(id int) (Snippet, error) {
    // Write the SQL statement we want to execute. Again, I've split it over two
    // lines for readability.
    stmt := `SELECT id, title, content, created, expires FROM snippets
        WHERE expires > UTC_TIMESTAMP() AND id = ?
    `

    // Use the QueryRow() method on the connection pool to execute our
    // SQL statement, passing in the untrusted id variable as the value for the
    // placeholder parameter. This returns a pointer to a sql.Row object which
    // holds the result from the database.
    row := m.DB.QueryRow(stmt, id)

    // Initialize a new zeroed Snippet struct.
    var s Snippet

    // Use row.Scan() to copy the values from each field in sql.Row to the
    // corresponding field in the Snippet struct. Notice that the arguments
    // to row.Scan are *pointers* to the place you want to copy the data into,
    // and the number of arguments must be exactly the same as the number of
    // columns returned by your statement.
    err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
    if err != nil {
        // If the query returns no rows, then row.Scan() will return a
        // sql.ErrNoRows error. We use the errors.Is() function check for that
        // error specifically, and return our own ErrNoRecord error
        // instead (we'll create this in a moment).
        if errors.Is(err, sql.ErrNoRows) {
            return Snippet{}, ErrNoRecord
        } else {
            return Snippet{}, err
        }
    }

    // If everything went OK, then return the filled Snippet struct.
    return s, nil
}

...
```

**File: `internal/models/errors.go`**

```go
package models

import (
    "errors"
)

var ErrNoRecord = errors.New("models: no matching record found")
```

**File: `cmd/web/handlers.go`**

```go
package main

import (
    "errors" // New import
    "fmt"
    "html/template"
    "net/http"
    "strconv"

    "snippetbox.libra.dev/internal/models" // New import
)

...

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    // Use the SnippetModel's Get() method to retrieve the data for a
    // specific record based on its ID. If no matching record is found,
    // return a 404 Not Found response.
    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            http.NotFound(w, r)
        } else {
            app.serverError(w, r, err)
        }
        return
    }

    // Write the snippet data as a plain-text HTTP response body.
    fmt.Fprintf(w, "%+v", snippet)
}

...
```

## 4.8 Multiple-record SQL queries

**Querying Multiple Rows**

```go
// models/snippets.go
func (m *SnippetModel) Latest() ([]Snippet, error) {
    stmt := `SELECT id, title, content, created, expires
             FROM snippets
             WHERE expires > UTC_TIMESTAMP()
             ORDER BY id DESC LIMIT 10`

    rows, err := m.DB.Query(stmt)
    if err != nil {
        return nil, err
    }
    defer rows.Close() // Critical for connection management

    var snippets []Snippet
    for rows.Next() {
        var s Snippet
        if err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
            return nil, err
        }
        snippets = append(snippets, s)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }
    return snippets, nil
}

// handlers.go
func (app *application) home(w http.ResponseWriter, r *http.Request) {
    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    for _, s := range snippets {
        fmt.Fprintf(w, "%+v\n", s)
    }
}
```

**Key Points**

1. Always `defer rows.Close()` immediately after checking query error
2. Process rows with `rows.Next()`/`rows.Scan()` loop
3. Check for iteration errors with `rows.Err()`

**PostgreSQL Differences**

```go
stmt := `SELECT ... WHERE expires > CURRENT_TIMESTAMP ...`
```

---

**Example Code**

**File: `internal/models/snippets.go`**

```go
package models

...

func (m *SnippetModel) Latest() ([]Snippet, error) {
    // Write the SQL statement we want to execute.
    stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

    // Use the Query() method on the connection pool to execute our
    // SQL statement. This returns a sql.Rows resultset containing the result of
    // our query.
    rows, err := m.DB.Query(stmt)
    if err != nil {
        return nil, err
    }

    // We defer rows.Close() to ensure the sql.Rows resultset is
    // always properly closed before the Latest() method returns. This defer
    // statement should come *after* you check for an error from the Query()
    // method. Otherwise, if Query() returns an error, you'll get a panic
    // trying to close a nil resultset.
    defer rows.Close()

    // Initialize an empty slice to hold the Snippet structs.
    var snippets []Snippet

    // Use rows.Next to iterate through the rows in the resultset. This
    // prepares the first (and then each subsequent) row to be acted on by the
    // rows.Scan() method. If iteration over all the rows completes then the
    // resultset automatically closes itself and frees-up the underlying
    // database connection.
    for rows.Next() {
        // Create a new zeroed Snippet struct.
        var s Snippet
        // Use rows.Scan() to copy the values from each field in the row to the
        // new Snippet object that we created. Again, the arguments to row.Scan()
        // must be pointers to the place you want to copy the data into, and the
        // number of arguments must be exactly the same as the number of
        // columns returned by your statement.
        err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
        if err != nil {
            return nil, err
        }
        // Append it to the slice of snippets.
        snippets = append(snippets, s)
    }

    // When the rows.Next() loop has finished we call rows.Err() to retrieve any
    // error that was encountered during the iteration. It's important to
    // call this - don't assume that a successful iteration was completed
    // over the whole resultset.
    if err = rows.Err(); err != nil {
        return nil, err
    }

    // If everything went OK then return the Snippets slice.
    return snippets, nil
}
```

**File: `cmd/web/handlers.go`**

```go
package main

import (
    "errors"
    "fmt"
    // "html/template"
    "net/http"
    "strconv"

    "snippetbox.libra.dev/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")

    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    for _, snippet := range snippets {
        fmt.Fprintf(w, "%+v\n", snippet)
    }

    // files := []string{
    //     "./ui/html/base.tmpl",
    //     "./ui/html/partials/nav.tmpl",
    //     "./ui/html/pages/home.tmpl",
    // }

    // ts, err := template.ParseFiles(files...)
    // if err != nil {
    //     app.serverError(w, r, err)
    //     return
    // }

    // err = ts.ExecuteTemplate(w, "base", nil)
    // if err != nil {
    //     app.serverError(w, r, err)
    // }
}

...
```

## 4.9 Transactions and other details

**`database/sql` Package**

- Works with any SQL database (MySQL, PostgreSQL, SQLite)
- Each driver may have specific behaviors - check documentation

**Handling NULL Values**

- Go does not handle `NULL` values in database records well by default

```go
// Solution 1: sql.Null*
type Book struct {
    Isbn  string
    Title  sql.NullString  // For nullable string fields
    Author sql.NullString
    Price  sql.NullFloat64 // For nullable numeric fields
}

...

for _, bk := range bks {
  var price string
  // Check validity before use:
  if bk.Price.Valid {
    price = fmt.Sprintf("£%.2f", bk.Price.Float64)
  } else {
    price = "PRICE NOT SET"
  }
  fmt.Printf("%s, %s, %s, %s\n", bk.Isbn, bk.Title.String, bk.Author.String, price)
}

// Should output:

// 978-1503261969, Emma, Jayne Austen, £9.44
// 978-1514274873, Journal of a Soldier, , £5.49
// 978-1503379640, The Prince, Niccolò Machiavelli, PRICE NOT SET
```

```go
// Solution 2: database constraints
CREATE TABLE books (
    isbn    char(14) NOT NULL,
    title   varchar(255),
    author  varchar(255),
    price   decimal(5,2)
);

INSERT INTO books (isbn, title, author, price) VALUES
    ('978-1503261969', 'Emma', 'Jayne Austen', 9.44),
    ('978-1514274873', 'Journal of a Soldier', NULL, 5.49),
    ('978-1503379640', 'The Prince', 'Niccolò Machiavelli', NULL);
```

Full example [here](https://gist.github.com/alexedwards/dc3145c8e2e6d2fd6cd9).

**Transactions**

- Use transactions for atomic operations

```go
func (m *Model) TransferFunds() error {
    tx, err := m.DB.Begin()
    defer tx.Rollback() // Safety net

    _, err = tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = 1")
    _, err = tx.Exec("UPDATE accounts SET balance = balance + 100 WHERE id = 2")

    return tx.Commit() // Or error will trigger Rollback
}
```

Key points:

- Always call Rollback or Commit
- Use same tx for all operations in transaction

**Prepared Statements**

- Prepared statements improve performance for repeated queries

```go
type UserModel struct {
    DB        *sql.DB
    InsertStmt *sql.Stmt
}

func NewUserModel(db *sql.DB) (*UserModel, error) {
    stmt, err := db.Prepare("INSERT INTO users(...) VALUES(...)")
    return &UserModel{DB: db, InsertStmt: stmt}, err
}

func (m *UserModel) Insert(user User) error {
    _, err := m.InsertStmt.Exec(user.Name, user.Email)
    return err
}
```

Key points:

- Prepare during app startup
- Close statements on shutdown
- Only use for frequently repeated queries (Under heavy load, statements may be re-prepared on different connections, reducing performance gains)

---

**Example Code**

```go
type ExampleModel struct {
    DB *sql.DB
}

func (m *ExampleModel) ExampleTransaction() error {
    // Calling the Begin() method on the connection pool creates a new sql.Tx
    // object, which represents the in-progress database transaction.
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }

    // Defer a call to tx.Rollback() to ensure it is always called before the
    // function returns. If the transaction succeeds it will be already be
    // committed by the time tx.Rollback() is called, making tx.Rollback() a
    // no-op. Otherwise, in the event of an error, tx.Rollback() will rollback
    // the changes before the function returns.
    defer tx.Rollback()

    // Call Exec() on the transaction, passing in your statement and any
    // parameters. It's important to notice that tx.Exec() is called on the
    // transaction object just created, NOT the connection pool. Although we're
    // using tx.Exec() here you can also use tx.Query() and tx.QueryRow() in
    // exactly the same way.
    _, err = tx.Exec("INSERT INTO ...")
    if err != nil {
        return err
    }

    // Carry out another transaction in exactly the same way.
    _, err = tx.Exec("UPDATE ...")
    if err != nil {
        return err
    }

    // If there are no errors, the statements in the transaction can be committed
    // to the database with the tx.Commit() method.
    err = tx.Commit()
    return err
}
```

```go
// We need somewhere to store the prepared statement for the lifetime of our
// web application. A neat way is to embed it in the model alongside the
// connection pool.
type ExampleModel struct {
    DB         *sql.DB
    InsertStmt *sql.Stmt
}

// Create a constructor for the model, in which we set up the prepared
// statement.
func NewExampleModel(db *sql.DB) (*ExampleModel, error) {
    // Use the Prepare method to create a new prepared statement for the
    // current connection pool. This returns a sql.Stmt object which represents
    // the prepared statement.
    insertStmt, err := db.Prepare("INSERT INTO ...")
    if err != nil {
        return nil, err
    }

    // Store it in our ExampleModel struct, alongside the connection pool.
    return &ExampleModel{DB: db, InsertStmt: insertStmt}, nil
}

// Any methods implemented against the ExampleModel struct will have access to
// the prepared statement.
func (m *ExampleModel) Insert(args...) error {
    // We then need to call Exec directly against the prepared statement, rather
    // than against the connection pool. Prepared statements also support the
    // Query and QueryRow methods.
    _, err := m.InsertStmt.Exec(args...)

    return err
}

// In the web application's main function we will need to initialize a new
// ExampleModel struct using the constructor function.
func main() {
    db, err := sql.Open(...)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
    defer db.Close()

    // Use the constructor function to create a new ExampleModel struct.
    exampleModel, err := NewExampleModel(db)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    // Defer a call to Close() on the prepared statement to ensure that it is
    // properly closed before our main function terminates.
    defer exampleModel.InsertStmt.Close()
}
```

# Chapter 6: Middleware

## 6.1 How middleware works

**Middleware Basics**

- Sits between server and handlers in request chain
- Must call `next.ServeHTTP()` to continue chain

**Implementation Patterns**

```go
// Closure style
func myMiddleware(next http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        // Logic
        next.ServeHTTP(w, r)
    }
    return http.HandlerFunc(fn)
}

// Anonymous function style
func myMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Logic
        next.ServeHTTP(w, r)
    })
}
```

- **Placement Impact**
  - *Before servemux*: Affects all requests (e.g., logging)

  ```
  middleware -> servemux -> handler
  ```

- *After servemux*: Route-specific (e.g., auth)

  ```
  servemux -> middleware -> handler
  ```

**Key Points**

- Position determines scope (global vs route-level)
- Common uses: logging, auth, error handling, headers
- Keep middleware lightweight for performance

##  6.2 Setting common headers

**Implementation**

```go
// cmd/web/middleware.go
func commonHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Server", "Go")
        ...
        next.ServeHTTP(w, r)
    })
}

// cmd/web/routes.go
func (app *application) routes() http.Handler {
    mux := http.NewServeMux()
    mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static/"))))
    mux.HandleFunc("GET /{$}", app.home)
    mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
    mux.HandleFunc("GET /snippet/create", app.snippetCreate)
    mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)
    return commonHeaders(mux) // Wrap servemux with middleware
}
```

**Key Points**

- `routes()` should return `http.Handler` instead of `*http.ServeMux`

**Middleware Flow**

```plaintext
commonHeaders → servemux → handler → servemux → commonHeaders
```

**Execution Order**

```go
func myMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Downstream logic (before next handler)
        next.ServeHTTP(w, r)
        // Upstream logic (after handler completes)
    })
}
```

**Early Termination**

- Can stop chain execution and control will flow back upstream (e.g., for auth failures)

```go
func myMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // If the user isn't authorized, send a 403 Forbidden status and
        // return to stop executing the chain.
        if !isAuthorized(r) {
            w.WriteHeader(http.StatusForbidden)
            return
        }

        // Otherwise, call the next handler in the chain.
        next.ServeHTTP(w, r)
    })
}
```

---

**[Example Code](https://github.com/Libra2021/snippet-box/commit/b8755411ea573e1811f38f8ae37644e7718f126a)**

## 6.3 Request logging

**Implementation**

```go
// middleware.go
func (app *application) logRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var (
            ip     = r.RemoteAddr
            proto  = r.Proto
            method = r.Method
            uri    = r.URL.RequestURI()
        )

        app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

        next.ServeHTTP(w, r)
    })
}

// routes.go
...

return app.logRequest(commonHeaders(mux))

...
```

**Key Points**

- Implementing the middleware as a method against `application` to access the handler dependencies

**Middleware Flow**

```plaintext
logRequest ↔ commonHeaders ↔ servemux ↔ application handler
```

---

[Example Code](https://github.com/Libra2021/snippet-box/commit/4d9101b657ce12fbf05add468f34a09b573f1186)

## 6.4 Panic recovery

**Panic Recovery Middleware**

- Prevents server crashes from handler panics
- Returns 500 response instead of empty reply

**Implementation**

```go
// middleware.go
func (app *application) recoverPanic(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                w.Header().Set("Connection", "close")  // Force connection close
                app.serverError(w, r, fmt.Errorf("%s", err))  // Return 500
            }
        }()
        next.ServeHTTP(w, r)
    })
}

// routes.go
func (app *application) routes() http.Handler {
    mux := http.NewServeMux()
    // ... route setup
    return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
```

**Key Points**

- Use a deferred function

- `recover()` catches panics in current goroutine only
- `Connection: close` ensures clean termination
- Middleware order matters (place `recoverPanic` first)

**Background Goroutines**

- Our middleware will only recover panics that happen in the same goroutine that executed the `recoverPanic()` middleware

```go
func (app *application) myHandler(w http.ResponseWriter, r *http.Request) {
    ...

    // Spin up a new goroutine to do some background processing.
    go func() {
        defer func() {
            if err := recover(); err != nil {
                app.logger.Error(fmt.Sprint(err))  // Must handle separately
            }
        }()

        doSomeBackgroundProcessing()
    }()

    w.Write([]byte("OK"))
}
```

**Response Example**
```bash
HTTP/1.1 500 Internal Server Error
Connection: close
...
Content-Length: 22
Internal Server Error
```

**Note**

- Panics in handlers terminate the goroutine of current request but the main goroutine still work

---

**[Example Code](https://github.com/Libra2021/snippet-box/commit/77214415f634bd228987eeba49fbfcf7c69996ae)**

## 6.5 Composable middleware chains

**`justinas/alice` Package**

- Helps us manage our middleware/handler chains

```go
// Without alice
Middleware1(Middleware2(Middleware3(myHandler)))

// with alice
alice.New(Middleware1, Middleware2, Middleware3).Then(myHandler)
```
- Can be used to create middleware chains that can be assigned to variables, appended to, and reused

```go
myChain := alice.New(myMiddlewareOne, myMiddlewareTwo)
myOtherChain := myChain.Append(myMiddleware3)
return myOtherChain.Then(myHandler)
```

---

**[Example Code](https://github.com/Libra2021/snippet-box/commit/aa3b0b09bd936b5b2ce330a45f460f6d87b3bd48)**

# Chapter 7: Processing forms

## 7.2 Parsing form data

**Form Handling Essentials**

```go
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    title := r.PostForm.Get("title")       // String value
    content := r.PostForm.Get("content")   // String value
    expires, err := strconv.Atoi(r.PostForm.Get("expires"))  // Convert to int
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
```

**Key Points**

- **`r.ParseForm()`**
  - Parses POST/PUT/PATCH bodies into `r.PostForm`
  - Idempotent (safe to call multiple times)
  - Returns an error if the body is empty or too large
  - 10MB limit (except multipart forms `enctype="multipart/form-data"`)

- **`r.PostForm.Get()`**

  - Only returns the first value as *string*
  - Requires manual type conversion (e.g., `strconv.Atoi()`)

  - The type is `url.Values` which has the underlying type `map[string][]string`
  - **`r.PostFormValue()`** is a shorthand but it silently ignores parse errors

**Advanced Handling**

- **Multiple Values**:

  ```go
  for i, item := range r.PostForm["items"] {  // Checkboxes, multi-select
      fmt.Fprintf(w, "%d: Item %s\n", i, item)
  }
  ```

- **Size Limits**:

  ```go
  r.Body = http.MaxBytesReader(w, r.Body, 4096)  // Set 4KB limit
  err := r.ParseForm()
  ```

- **Query Strings**:

  ```go
  search := r.URL.Query().Get("q")  // GET /search?q=foo
  ```

**Data Sources**

| Method                | Source                   | Best For             |
| :-------------------- | :----------------------- | :------------------- |
| `r.PostForm.Get()`    | POST body only           | Form submissions     |
| `r.URL.Query().Get()` | URL query string         | Search/filter params |
| `r.Form.Get()`        | POST body + query string | Mixed data (avoid)   |

**Best Practices**

- Always call `r.ParseForm()` first
- Use explicit sources (`PostForm` or `Query`)
- Avoid `PostFormValue()` (hides parse errors)
- Handle type conversions explicitly
- Set size limits for security

---

[**Example Code**](https://github.com/Libra2021/snippet-box/commit/01001026914f409d7441a517134a75156706dd12)

## 7.3 Validating form data

**Form Validating Essentials**

```go
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    title := r.PostForm.Get("title")
    content := r.PostForm.Get("content")

    expires, err := strconv.Atoi(r.PostForm.Get("expires"))
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    fieldErrors := make(map[string]string)

    if strings.TrimSpace(title) == "" {
        fieldErrors["title"] = "This field cannot be blank"
    } else if utf8.RuneCountInString(title) > 100 {
        fieldErrors["title"] = "This field cannot be more than 100 characters long"
    }

    if strings.TrimSpace(content) == "" {
        fieldErrors["content"] = "This field cannot be blank"
    }

    if expires != 1 && expires != 7 && expires != 365 {
        fieldErrors["expires"] = "This field must equal 1, 7 or 365"
    }

    if len(fieldErrors) > 0 {
        fmt.Fprint(w, fieldErrors)
        return
    }

    id, err := app.snippets.Insert(title, content, expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
```

**Key Points**

- **Validation**:
  - Use `strings.TrimSpace()` to check for blank values
  - Use `utf8.RuneCountInString()` to check for long values
- **Error Handling**:
  - Use a map to store validation errors
  - Return the errors to the user

**`utf8.RuneCountInString()` vs `len()`**

- `len(str)` → Byte count (O(1))
- `utf8.RuneCountInString(str)` → Unicode character count (O(n))

```go
utf8.RuneCountInString("Hello, 世界") // 8
len("Hello, 世界") // 13

utf8.RuneCountInString("Zoë") // 3
len("Zoë") // 4
```

**Additional Resources**

- [Validation Snippets for Go](https://www.alexedwards.net/blog/validation-snippets-for-go)

---

[**Example Code**](https://github.com/Libra2021/snippet-box/commit/b24d0b44118297855f35166b88aa42e91c036a16)

## 7.4 Displaying errors and repopulating fields

**Template Data Structure**

```go
// templates.go
// Add a Form field with the type "any".
type templateData struct {
    CurrentYear int
    Snippet     models.Snippet
    Snippets    []models.Snippet
    Form        any
}

// handlers.go
type snippetCreateForm struct {
    Title       string
    Content     string
    Expires     int
    FieldErrors map[string]string
}
```

**Form Handling Flow**

```go
// Display form (GET)
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)
    data.Form = snippetCreateForm{Expires: 365}  // Set defaults
    app.render(w, r, http.StatusOK, "create.tmpl", data)
}

// Process form (POST)
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    // Parse and validate
    form := snippetCreateForm{
        Title:   r.PostForm.Get("title"),
        Content: r.PostForm.Get("content"),
        Expires: expires,
        FieldErrors: map[string]string{},
    }

    // Validation
    if strings.TrimSpace(form.Title) == "" {
        form.FieldErrors["title"] = "Cannot be blank"
    } else if utf8.RuneCountInString(form.Title) > 100 {
        form.FieldErrors["title"] = "Max 100 characters"
    }
    // ... other validations

    if len(form.FieldErrors) > 0 {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
        return
    }

    // Process valid form
    id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
    // ... handle success/error
}
```

**Template Implementation**

```html
<form action='/snippet/create' method='POST'>
    <!-- Title Field -->
    <div>
        <label>Title:</label>
        {{with .Form.FieldErrors.title}}<label class='error'>{{.}}</label>{{end}}
        <input type='text' name='title' value='{{.Form.Title}}'>
    </div>

    <!-- Content Field -->
    <div>
        <label>Content:</label>
        {{with .Form.FieldErrors.content}}<label class='error'>{{.}}</label>{{end}}
        <textarea name='content'>{{.Form.Content}}</textarea>
    </div>

    <!-- Expires Field -->
    <div>
        <label>Delete in:</label>
        {{with .Form.FieldErrors.expires}}<label class='error'>{{.}}</label>{{end}}
        <input type='radio' name='expires' value='365' {{if (eq .Form.Expires 365)}}checked{{end}}> One Year
        <input type='radio' name='expires' value='7' {{if (eq .Form.Expires 7)}}checked{{end}}> One Week
        <input type='radio' name='expires' value='1' {{if (eq .Form.Expires 1)}}checked{{end}}> One Day
    </div>
</form>
```

**Key Points**

- **Form Struct Pattern**

  - Combines form data + validation errors
  - Exported fields for template access
    - `{{.Form.Title}}`
  - Map key names don’t have to be capitalized
    - `{{.Form.FieldErrors.title}}`

- **Initial State**

  - Initialize `snippetCreateForm` fields in `snippetCreate` in case `Form` is `nil`
  - Set default values here

  ```go
  data.Form = snippetCreateForm{
    Expires: 365
  }
  ```

- **Validation Flow**

  - Parse → Validate → Render or Process
  - 422 status for validation errors
  - Automatic form repopulation

- **Template Techniques**

  - `with` for conditional error display
  - `value`/`checked` attributes for repopulation
  - `eq` for comparing values

---

**[Example Code]**(https://github.com/Libra2021/snippet-box/commit/02843fe5cba059de2656b724b53df2f3eccbb049)

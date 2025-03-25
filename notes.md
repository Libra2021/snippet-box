# Chapter 2: Foundations

## 2.1 Project setup and creating a module

- In the Go community, a common convention is to base your module paths on a URL that you own.

- If you’re creating a project which can be downloaded and used by other people and programs, then it’s good practice for your module path to equal the location that the code can be downloaded from.

  For instance, if your package is hosted at `https://github.com/foo/bar` then the module path for the project should be `github.com/foo/bar`.

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

- The `home` handler function is just a regular Go function with two parameters.

- Any error returned by `http.ListenAndServe()` is always non-nil.

- Go’s servemux treats the route pattern `"/"` like a catch-all.

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

- When a pattern doesn’t have a trailing slash, it will only be matched when the request URL path exactly matches the pattern in full.

- When a route pattern ends with a trailing slash — like `"/"` or `"/static/"` — it is known as a subtree path pattern. Subtree path patterns are matched whenever the start of a request URL path matches the subtree path. You can think of subtree paths as acting a bit like they have a wildcard at the end, like `"/**"` or `"/static/**"`.

- Append the special character sequence `{$}` to the end of the pattern — like `"/{$}"` or `"/static/{$}"` to prevent subtree path patterns.

- Request URL paths are automatically sanitized. If the request path contains any `.` or `..` elements or repeated slashes, the user will automatically be redirected to an equivalent clean URL. For example, if a user makes a request to `/foo/bar/..//baz` they will automatically be sent a `301 Permanent Redirect` to `/foo/baz` instead.

- If a subtree path has been registered and a request is received for that subtree path without a trailing slash, then the user will automatically be sent a `301 Permanent Redirect` to the subtree path with the slash added. For example, if you have registered the subtree path /foo/, then any request to `/foo` will be redirected to `/foo/`.

- It’s possible to include host names in your route patterns. This can be useful when you want to redirect all HTTP requests to a canonical URL, or if your application is acting as the back end for multiple sites or services. For example:

  ```go
  mux := http.NewServeMux()
  mux.HandleFunc("foo.example.org/", fooHandler)
  mux.HandleFunc("bar.example.org/", barHandler)
  mux.HandleFunc("/baz", bazHandler)
  ```

  When it comes to pattern matching, any host-specific patterns will be checked first and if there is a match the request will be dispatched to the corresponding handler. Only when there isn’t a host-specific match found will the non-host specific patterns also be checked.

- `http.Handle()` and `http.HandleFunc()` allow you to register routes without explicitly declaring a servemux, like this:

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

  Behind the scenes, these functions register their routes with something called the _default servemux_. This is just a regular servemux like we’ve already been using, but which is initialized automatically by Go and stored in the `http.DefaultServeMux` global variable.

  If you pass `nil` as the second argument to `http.ListenAndServe()`, the server will use `http.DefaultServeMux` for routing.

  But for the sake of clarity, maintainability and security, it’s generally a good idea to avoid `http.DefaultServeMux` and the corresponding helper functions.

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

- Wildcard segments in a route pattern are denoted by an wildcard identifier inside `{}` brackets. Like this:

  ```go
  mux.HandleFunc("/products/{category}/item/{itemID}", exampleHandler)
  ```

  Patterns like `"/products/c_{category}"`, `/date/{y}-{m}-{d}` or `/{slug}.html` are not valid.

- Inside your handler, you can retrieve the corresponding value for a wildcard segment using its identifier and the `r.PathValue()` method. For example:

  ```go
  func exampleHandler(w http.ResponseWriter, r *http.Request) {
      category := r.PathValue("category")
      itemID := r.PathValue("itemID")

      ...
  }
  ```

  The `r.PathValue()` method always returns a string value and you should validate or sanity check the value before doing anything important with it.

**Precedence and conflicts**

- The most specific route pattern wins. For example: `/post/edit` will match `"/post/edit"` instead of `"/post/{id}"`.

- `/post/new/latest` won't match `"/post/new/{id}"` and `"/post/{author}/latest"` because Go’s servemux considers the patterns to conflict, and will panic at runtime when initializing the routes.
- keep overlaps to a minimum or avoid them completely.

**Subtree path patterns with wildcards**

- The routing rules we described in the previous chapter still apply. For example: `"/user/{id}/"` will match `/user/1/`, `/user/2/a` and `/user/2/a/b/c`.

**Remainder wildcards**

If a route pattern ends with a wildcard, and this final wildcard identifier ends in `...`, then the wildcard will match any and all remaining segments of a request path. For example, `"/post/{path...}"` will match `/post/a`, `/post/a/b`, `/post/a/b/c` — very much like a subtree path pattern. But you can access the entire wildcard part via `r.PathValue()`. In this example, `r.PathValue("path")` would return `"a/b/c"`.

## 2.5 Method-based routing

1. **Restricting Routes to Specific HTTP Methods**

- To restrict a route to a specific HTTP method, prefix the route pattern with the HTTP method.

- Example:

  ```go
  mux.HandleFunc("GET /{$}", home)
  mux.HandleFunc("GET /snippet/view/{id}", snippetView)
  mux.HandleFunc("GET /snippet/create", snippetCreate)
  ```

- **Note**: HTTP methods in route patterns are **case-sensitive** and must be **uppercase**, followed by at least one whitespace character.

- A route registered with the `GET` method matches both `GET` and `HEAD` requests.
- Other methods (`POST`, `PUT`, `DELETE`) require an exact match.

- It's acceptable to declare multiple routes with the same pattern but different HTTP methods.
  ```go
  mux.HandleFunc("GET /snippet/create", snippetCreate)
  mux.HandleFunc("POST /snippet/create", snippetCreatePost)
  ```

2. **Handling Unsupported Methods**

- If a request uses an unsupported method, Go's `ServeMux` automatically sends a `405 Method Not Allowed` response which includes an `Allow` header listing supported methods.

3. **Testing with `curl`**

- **GET Request**:

  ```bash
  curl -i localhost:4000/snippet/create
  ```

- **HEAD Request**:

  ```bash
  curl --head localhost:4000/snippet/create
  ```

- **POST Request**:

  ```bash
  curl -i -d "" localhost:4000/snippet/create
  ```

- **DELETE Request**:
  ```bash
  curl -i -X DELETE localhost:4000/snippet/create
  ```

4. **Method Precedence**

- The most specific pattern wins.
- A route pattern without a method (e.g., `"/article/{id}"`) matches requests with any method.
- A route pattern with a method (e.g., `"POST /article/{id}"`) takes precedence over a pattern without a method.

5. **Handler Naming Conventions**

- No strict rules for naming handlers in Go.
- Common conventions:
  - Postfix `Post` for handlers dealing with `POST` requests (e.g., `snippetCreatePost`).
  - Prefix `get` or `post` (e.g., `getSnippetCreate`, `postSnippetCreate`).
  - Use descriptive names (e.g., `newSnippetForm`, `createSnippet`).

6. **Limitations of Standard Library Routing**

- **Not Supported**:

  - Custom `404 Not Found` and `405 Method Not Allowed` responses.
  - Regular expressions in route patterns or wildcards.
  - Matching multiple HTTP methods in a single route declaration.
  - Automatic support for `OPTIONS` requests.
  - Routing based on HTTP request headers.

- **Recommended Third-Party Routers**:
  - `httprouter`
  - `chi`
  - `flow`
  - `gorilla/mux`
    Comparison and guidance in this [blog](https://www.alexedwards.net/blog/which-go-router-should-i-use).

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

1. **Default Response Behavior**

- By default, responses have:

  - Status code: `200 OK`
  - Headers: `Date`, `Content-Length`, and `Content-Type` (automatically generated).

- Example:

  ```bash
  $ curl -i localhost:4000/
  HTTP/1.1 200 OK
  Date: Wed, 18 Mar 2024 11:29:23 GMT
  Content-Length: 21
  Content-Type: text/plain; charset=utf-8

  Hello from Snippetbox
  ```

2. **Customizing HTTP Status Codes**

- Use `w.WriteHeader(statusCode)` to set a custom status code.

- Example:

  ```go
  func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
      w.WriteHeader(http.StatusCreated) // w.WriteHeader(201)
      w.Write([]byte("Save a new snippet..."))
  }
  ```

- Use constants from the `net/http` package for clarity and to avoid typos.

- **Important**:
  - `w.WriteHeader()` can only be called **once** per response.
  - If not called explicitly, the first `w.Write()` will send a `200 OK` status.

3. **Customizing Response Headers**

- Use `w.Header().Add()` to add custom headers.

- Example:

  ```go
  func home(w http.ResponseWriter, r *http.Request) {
      w.Header().Add("Server", "Go")
      w.Write([]byte("Hello from Snippetbox"))
  }
  ```

- **Important**:
  - Headers must be set **before** calling `w.WriteHeader()` or `w.Write()`. Changes to headers after calling `w.WriteHeader()` or `w.Write()` are ignored.

4. **Writing Response Bodies**

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

5. **Content Sniffing**

- Go automatically sets the `Content-Type` header using `http.DetectContentType()` based on the response body except for JSON.

- For JSON responses, manually set the `Content-Type` header:
  ```go
  w.Header().Set("Content-Type", "application/json")
  w.Write([]byte(`{"name":"Alex"}`))
  ```

6. **Manipulating the Header Map**

- Methods for manipulating headers:

  - `Set()`: Overwrites an existing header.
  - `Add()`: Appends a new header (can be called multiple times).
  - `Del()`: Deletes a header.
  - `Get()`: Retrieves the first value of a header.
  - `Values()`: Retrieves all values of a header.

- Example:

  ```go
  w.Header().Set("Cache-Control", "public, max-age=31536000")

  w.Header().Add("Cache-Control", "public")
  w.Header().Add("Cache-Control", "max-age=31536000")

  w.Header().Del("Cache-Control")

  w.Header().Get("Cache-Control")

  w.Header().Values("Cache-Control")
  ```

7. **Header Canonicalization**

- Header names are canonicalized (e.g., `cache-control` → `Cache-Control`) using `textproto.CanonicalMIMEHeaderKey()` which means that when calling these methods (`Set()`, `Add()`, `Del()`, `Get()` and `Values()`) the header name is case-insensitive.

- To avoid canonicalization, modify the header map directly:

  ```go
  w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
  ```

- If a HTTP/2 connection is being used, Go will always automatically convert the header names and values to lowercase for you when writing the response.

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

1. **Project Structure Overview**

- There’s no single "right" way to structure Go projects, but a common approach is:

  - `cmd/`: Contains application-specific code (e.g., `cmd/web` for the web application).
  - `internal/`: Contains reusable, non-application-specific code (e.g., validation helpers, database models).
  - `ui/`: Contains user-interface assets (e.g., HTML templates in `ui/html`, static files in `ui/static`).

- Example:
  ```
  snippetbox/
  ├── cmd/
  │   └── web/
  │       ├── main.go
  │       └── handlers.go
  ├── internal/
  ├── ui/
  │   ├── html/
  │   └── static/
  ```

2. **Benefits of This Structure**

- **Separation of Concerns**:

  - Go code lives under `cmd` and `internal`.
  - Non-Go assets (e.g., HTML, CSS) live under `ui`.

- **Scalability**:
  - Easy to add additional executables (e.g., a CLI under `cmd/cli`).
  - Reusable code in `internal` can be shared across executables.

3. **Running the Application**

- Use `go run` to start the application:
  ```bash
  cd $HOME/code/snippetbox
  go run ./cmd/web
  ```

4. **The `internal` Directory**

- Packages under `internal` can only be imported by code inside the parent of the `internal` directory.
- Prevents external codebases from importing and relying on internal packages.

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

## 2.8 HTML templating and inheritance

1. **Template File Creation**

**File: `ui/html/pages/home.tmpl`**

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>Home - Snippetbox</title>
  </head>
  <body>
    <header>
      <h1><a href="/">Snippetbox</a></h1>
    </header>
    <main>
      <h2>Latest Snippets</h2>
      <p>There's nothing to see here yet!</p>
    </main>
    <footer>Powered by <a href="https://golang.org/">Go</a></footer>
  </body>
</html>
```

2. **Rendering the Template in Go**

- Use the `html/template` package to parse and render the template.

**File: `cmd/web/handlers.go`**

```go
package main

import (
    "html/template"
    "log"
    "net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")

    ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
    if err != nil {
        log.Print(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = ts.Execute(w, nil)
    if err != nil {
        log.Print(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
```

- **Note**:
  - `template.ParseFiles()` function must either be relative to your current working directory, or an absolute path.
  - The last parameter to `Execute()` represents any dynamic data that we want to pass in.

3. **Template Composition**

**File: `ui/html/base.tmpl`**

```html
{{define "base"}}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>{{template "title" .}} - Snippetbox</title>
  </head>
  <body>
    <header>
      <h1><a href="/">Snippetbox</a></h1>
    </header>
    {{template "nav" .}}
    <main>{{template "main" .}}</main>
    <footer>Powered by <a href="https://golang.org/">Go</a></footer>
  </body>
</html>
{{end}}
```

**File: `ui/html/pages/home.tmpl`**

```html
{{define "title"}}Home{{end}} {{define "main"}}
<h2>Latest Snippets</h2>
<p>There's nothing to see here yet!</p>
{{end}}
```

**File: `ui/html/partials/nav.tmpl`**

```html
{{define "nav"}}
<nav>
  <a href="/">Home</a>
</nav>
{{end}}
```

**File: `cmd/web/handlers.go`**

```go
package main

...

func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")

    files := []string{
        "./ui/html/base.tmpl",
        "./ui/html/partials/nav.tmpl",
        "./ui/html/pages/home.tmpl",
    }

    ts, err := template.ParseFiles(files...)
    if err != nil {
        log.Print(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = ts.ExecuteTemplate(w, "base", nil)
    if err != nil {
        log.Print(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

...
```

- **Note**:
  - The dot at the end of the {{template "title" .}} action represents any dynamic data that you want to pass to the invoked template.
  - The `{{block}}` action allows you to specify default content if a template doesn’t exist. For example:
    ```html
    {{define "base"}}
    <h1>An example template</h1>
    {{block "sidebar" .}}
    <p>My default sidebar content</p>
    {{end}} {{end}}
    ```
    You can also leave the default content empty which means the invoked template acts like it’s ‘optional’.

## 2.9 Serving static files

```bash
cd $HOME/code/snippetbox
curl https://www.alexedwards.net/static/sb-v2.tar.gz | tar -xvz -C ./ui/static/
```

1. **The http.Fileserver handler**

- Go's `net/http` package provides `http.FileServer` to serve files over HTTP from a specific directory.

  ```go
  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
  ```

- `http.StripPrefix` removes `/static` from the URL path before passing it to the file server.
  - `/static/css/style.css` → `css/style.css` → `./ui/static/css/style.css`

2. **Additional Features of `http.FileServer`**

- **Path Sanitization**:
  Automatically cleans paths using `path.Clean()` to prevent directory traversal attacks.

- **Range Requests**:
  Supports partial content requests for large files (e.g., resumable downloads).

  ```bash
  $ curl -i -H "Range: bytes=100-199" --output - http://localhost:4000/static/img/logo.png
  HTTP/1.1 206 Partial Content
  Accept-Ranges: bytes
  Content-Length: 100
  Content-Range: bytes 100-199/1075
  Content-Type: image/png
  Last-Modified: Wed, 18 Mar 2024 11:29:23 GMT
  Date: Wed, 18 Mar 2024 11:29:23 GMT
  [binary data]
  ```

- **Caching**:
  Uses `Last-Modified` and `If-Modified-Since` headers to send `304 Not Modified` responses for unchanged files.

- **Content-Type Detection**:
  Automatically sets `Content-Type` based on file extensions using `mime.TypeByExtension()`.

3. **Serving Single Files**

- Use `http.ServeFile()` to serve individual files:

  ```go
  func downloadHandler(w http.ResponseWriter, r *http.Request) {
      http.ServeFile(w, r, "./ui/static/file.zip")
  }
  ```

- **Warning**: `http.ServeFile()` does not automatically sanitize the file path. Always sanitize file paths with `filepath.Clean()` when constructing paths from user input to prevent directory traversal attacks.

4. **Disabling Directory Listings**

- Simple Method: Add an empty `index.html` file to directories you want to hide.
- Advanced Method: Create a custom `http.FileSystem` implementation that returns `os.ErrNotExist` for directories. Refer to this [blog](https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings) for details.

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

1. **Overview**

A handler is an object that satisfies the `http.Handler` interface:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Here’s an example of a handler using a struct:

```go
type home struct {}

func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is my home page"))
}
```

This handler can be registered with a servemux using the `Handle` method:

```go
mux := http.NewServeMux()
mux.Handle("/", &home{})
```

2. **Handler Functions**

Creating a struct just to implement `ServeHTTP()` is often unnecessary. Instead, it's more common to write handlers as normal functions:

```go
func home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is my home page"))
}
```

However, this function isn't a handler by itself. To convert it into a handler, use the http.`HandlerFunc()` adapter:

```go
mux := http.NewServeMux()
mux.Handle("/", http.HandlerFunc(home))
```

The `HandleFunc()` method simplifies the process by transforming a function into a handler and registering it in one step:

```go
mux := http.NewServeMux()
mux.HandleFunc("/", home)
```

3. **Chaining Handlers**

The `http.ListenAndServe()` function takes a `http.Handler` as its second parameter:

```go
func ListenAndServe(addr string, handler Handler) error
```

A servemux can be passed to `ListenAndServe` because it also implements the `ServeHTTP()` method, satisfying the `http.Handler` interface. This allows the servemux to act as a special kind of handler that routes requests to other handlers.

- How it works:
  - The server receives an HTTP request.
  - It calls the servemux's `ServeHTTP()` method.
  - The servemux looks up the appropriate handler based on the request method and URL path.
  - The handler's `ServeHTTP()` method is called to generate the response.

In essence, a Go web application is a chain of ServeHTTP() methods being called sequentially.

4. **Requests are handled concurrently**

All incoming HTTP requests are served in their own goroutine. This makes Go servers highly efficient but also requires careful handling of shared resources to avoid race conditions.

### Key Takeaways

- A handler is any object that implements the `http.Handler` interface.
- Handlers can be created using structs or normal functions.
- The `http.HandlerFunc()` adapter converts functions into handlers.
- Servemuxes are special handlers that route requests to other handlers.
- Requests are handled concurrently, so be cautious with shared resources.

# Chapter 3: Configuration and error handling

## 3.1 Managing configuration settings

1. **Overview**

Our web application’s `main.go` file currently contains hard-coded configuration settings:

- The network address for the server to listen on (`":4000"`)
- The file path for the static files directory (`"./ui/static"`)

Hard-coding these settings isn’t ideal because:

- There’s no separation between configuration and code.
- We can’t change settings at runtime (important for different environments like development, testing, and production).

2. **Command-Line Flags**

A common and idiomatic way to manage configuration settings in Go is to use **command-line flags**. For example:

```bash
go run ./cmd/web -addr=":80"
```

To accept and parse a command-line flag, use the `flag.String()` function:

```go
addr := flag.String("addr", ":4000", "HTTP network address")
```

- **Name**: `addr`
- **Default value**: `":4000"`
- **Help text**: Explains what the flag controls.

**Note**:

- Ports 0-1023 are restricted and require root privileges. Attempting to use them may result in a `bind: permission denied` error.
- Command-line flags are optional. If no `-addr` flag is provided, the server falls back to the default value.
- Use the `-help` flag to list all available command-line flags and their help text:
  ```bash
  $ go run ./cmd/web -help
  Usage of /tmp/go-build3672328037/b001/exe/web:
    -addr string
          HTTP network address (default ":4000")
  ```

3. **Type Conversions**

- If the conversion fails, the application exits with an error.

- Go also provides other flag functions for different types:
  - `flag.Int()`
  - `flag.Bool()`
  - `flag.Float64()`
  - `flag.Duration()`

4. **Environment Variables**

You can store your configuration settings in environment variables and access them by using the `os.Getenv()`:

```go
addr := os.Getenv("SNIPPETBOX_ADDR")
```

But this has some drawbacks:

- No default values (`os.Getenv()` returns an empty string if the variable doesn’t exist).
- No automatic type conversion.
- No `-help` functionality.

However, you can combine both approaches by passing environment variables as command-line flags:

```bash
$ export SNIPPETBOX_ADDR=":9999"
$ go run ./cmd/web -addr=$SNIPPETBOX_ADDR
2024/03/18 11:29:23 starting server on :9999
```

5. **Boolean Flags**

For flags defined with `flag.Bool()`, omitting a value is equivalent to setting it to `true`:

```bash
go run ./example -flag=true
go run ./example -flag

go run ./example -flag=false
```

6. **Pre-Existing Variables**

You can parse command-line flag values into pre-existing variables using functions like `flag.StringVar()`, `flag.IntVar()`, and `flag.BoolVar()`. This is useful for storing all configuration settings in a single struct:

```go
type config struct {
    addr      string
    staticDir string
}

...

var cfg config

flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

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

### Key Takeaways

- Use **command-line flags** for runtime configuration.
- Leverage **default values** for convenience during development.
- Use **type-specific flag functions** (`flag.Int()`, `flag.Bool()`, etc.) for automatic type conversion.
- Combine **environment variables** with command-line flags for flexibility.
- Store configuration settings in a **struct** for better organization.

## 3.2 Structured logging

1. **Overview**

- Replace `log.Printf()` and `log.Fatal()` with structured logging using `log/slog`.
- Structured logs include:
  - Timestamp (millisecond precision)
  - Severity level (`Debug`, `Info`, `Warn`, `Error`)
  - Log message
  - Optional key-value attributes

2. **Creating a Structured Logger**

- Use `slog.NewTextHandler()` to create a handler for plaintext logs.
- Use `slog.NewJSONHandler()` for JSON-formatted logs.

- Example:

  ```go
  logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{...}))
  ```

  ```go
  logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
  ```

  Output:

  ```plaintext
  {"time":"2024-03-18T11:29:23.00000000+00:00","level":"INFO","msg":"starting server","addr":":4000"}
  {"time":"2024-03-18T11:29:23.00000000+00:00","level":"ERROR","msg":"listen tcp :4000: bind: address already in use"}
  ```

- **Handler Options**:

  - `os.Stdout`: Writes logs to standard output.
  - `slog.HandlerOptions`: Customizes handler behavior (e.g., log level, source location). Pass `nil` for the default.

    - **log level**:

      - `slog.LevelDebug`: Logs all messages.
      - `slog.LevelInfo`: Logs info, warn, and error messages.
      - `slog.LevelWarn`: Logs warn and error messages.
      - `slog.LevelError`: Logs only error messages.
      - ```go
        logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelDebug,
        }))
        ```

    - **Caller location**

      - Include filename and line number in logs:

        ```go
        logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            AddSource: true,
        }))
        ```

        ```plaintext
        time=2024-03-18T11:29:23.000+00:00 level=INFO source=/home/alex/code/snippetbox/cmd/web/main.go:32 msg="starting server" addr=:4000
        ```

3. **Using a Structured Logger**

- Log entries are created using `Debug()`, `Info()`, `Warn()`, or `Error()` methods which can accept an arbitrary number of additional attributes (key-value pairs).

- Example:

  ```go
  logger.Info("request received", "method", "GET", "path", "/")
  ```

  Output:

  ```plaintext
  time=2024-03-18T11:29:23.000+00:00 level=INFO msg="request received" method=GET path=/
  ```

- **Note**: If your attribute keys, values, or log message contain `"` or `=` characters or any whitespace, they will be wrapped in double quotes in the log output.

- **Safer Attributes**:

  - Bad Example:

    ```go
    logger.Info("starting server", "addr") // Oops, the value for "addr" is missing
    ```

    ```plaintext
    time=2024-03-18T11:29:23.000+00:00 level=INFO msg="starting server" !BADKEY=addr
    ```

  - Use `slog.Any()`, `slog.String()`, `slog.Int()`, `slog.Bool()`, `slog.Time()` and `slog.Duration()` for type-safe attributes.
  - Example:
    ```go
    logger.Info("starting server", slog.String("addr", ":4000"))
    ```

4. **Adding structured logging to our application**

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

- There is no structured logging equivalent to the `log.Fatal()`. Instead, `logger.Error()` and manually calling `os.Exit(1)` to terminate the application.

5. **Decoupled Logging**:

- Write logs to `os.Stdout` for flexibility.
- Redirect logs to a file in production:
  ```bash
  go run ./cmd/web >> /tmp/web.log
  ```
- **Note**:
  - The double arrow `>>` will append to an existing file.
  - The single arrow `>` will overwrite an existing file.

6. **Concurrent Logging**:

- `slog.New()` loggers are concurrency-safe. You can share a single logger across multiple goroutines.
- Ensure the destination's `Write()` method is also concurrency-safe if multiple loggers share the same destination.

### Key Takeaways

- **Structured logging** improves log readability, filtering, and parsing compared to standard logging.
- Use the `log/slog` package to create structured loggers with **severity levels** (`Debug`, `Info`, `Warn`, `Error`).
- Customize log output with **handler options**:
  - Set a **minimum log level** (e.g., `Debug`, `Info`).
  - Include **caller location** (filename and line number) for debugging.
  - Output logs in **JSON format** for machine readability.
- **Decouple logging** from application logic by writing logs to `os.Stdout` or redirect logs to files for persistent storage in production.
- Use **type-safe attributes** (e.g., `slog.String()`, `slog.Int()`) to avoid errors and improve code reliability.
- Structured loggers are **concurrency-safe**, making them suitable for use across multiple goroutines.

## 3.3 Dependency injection

1. **Overview**

- **Problem**:

  - The `home` handler in `handlers.go` still uses Go's standard logger (`log.Print`) instead of the structured logger.
  - Example:
    ```go
    func home(w http.ResponseWriter, r *http.Request) {
        ts, err := template.ParseFiles(files...)
        if err != nil {
            log.Print(err.Error()) // Still using the standard logger.
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
    }
    ```

- **Question**: How can we make the structured logger (and other dependencies) available to handlers?

- **Solution**:
  - Avoid global variables.
  - _Inject dependencies_ into handlers for better explicitness, testability, and reduced errors.
    - Use a **custom application struct** to hold dependencies and define handlers as methods on this struct.

2. **Implementation**

**Define the Application Struct**

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
    ...
}
```

**Update Handlers to Use the Struct**

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

**Wire Everything in `main.go`**

**File: `cmd/web/main.go`**

```go
package main

import (
    "flag"
    "log/slog"
    "net/http"
    "os"
)

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

3. **Closures for Dependency Injection**

- Use closures if handlers are spread across multiple packages.
- Example:

  ```go
  // package config

  type Application struct {
      Logger *slog.Logger
  }
  ```

  ```go
  // package foo

  func ExampleHandler(app *config.Application) http.HandlerFunc {
      return func(w http.ResponseWriter, r *http.Request) {
          ...
          ts, err := template.ParseFiles(files...)
          if err != nil {
              app.Logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
              http.Error(w, "Internal Server Error", http.StatusInternalServerError)
              return
          }
          ...
      }
  }
  ```

  ```go
  // package main

  func main() {
      app := &config.Application{
          Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
      }
      ...
      mux.Handle("/", foo.ExampleHandler(app))
      ...
  }
  ```

  - [This](https://gist.github.com/alexedwards/5cd712192b4831058b21) is a more concrete example.

### Key Takeaways

- Avoid global variables for dependencies; inject them instead.
- **Dependency injection** makes code more explicit, testable, and maintainable.
- Use a **custom application struct** to hold dependencies (e.g., logger, database connection).
- Define handlers as methods on the struct to access dependencies.
- For multi-package applications, use **closures** to inject dependencies into handlers.

## 3.4 Centralized error handling

1. **Overview**

- To improve code organization and reduce repetition, centralize error handling logic into helper methods.

2. **Implementation**

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

- **Note**:
  - `http.StatusText()` provides human-readable descriptions for HTTP status codes (e.g., `400` → `"Bad Request"`).
  - Use `debug.Stack()` to include a stack trace in log entries for better debugging. This returns a byte slice, which we need to convert to a string so that it's readable in the log entry.
    ```go
    var trace  = string(debug.Stack())
    ```

### Key Takeaways

- **Centralized error handling** reduces code duplication and improves maintainability.
- Use helper methods like `serverError` and `clientError` to handle errors consistently.
- Log errors with context (e.g., HTTP method, URI) for easier debugging.
- Include **stack traces** with `debug.Stack()` in logs for detailed debugging information.
- Use `http.StatusText()` to provide human-readable HTTP status descriptions.

## 3.5 Isolating the application routes

1. **Overview**

- To improve code organization, make `main()` only focused on:

  - Parsing runtime configuration.

    ```go
    addr := flag.String("addr", ":4000", "HTTP network address")
    flag.Parse()
    ```

  - Establishing dependencies.

    ```go
    app := &application{
        logger: logger,
    }
    ```

  - Running the HTTP server.
    ```go
    err := http.ListenAndServe(*addr, app.routes())
    ```

2. **Implementation**

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

### Key Takeaways

- Isolate route declarations in a separate `routes.go` file for better organization.
  - Use a `routes()` method on the `application` struct to return a `*ServeMux` containing all routes.
- Keep the `main()` function focused on configuration, dependency setup, and server startup.

# Chapter 4: Database-driven responses

## 4.3 Modules and reproducible builds

1. **Overview**

- Go modules ensure reproducible builds by specifying exact versions of dependencies.
- The `go.mod` and `go.sum` files work together to manage dependencies and verify their integrity.

2. **The `go.sum` File**

- The `go.mod` file lists the module path, Go version, and dependencies with their exact versions.

- Example:

  ```plaintext
  module snippetbox.libra.dev

  go 1.24.1

  require (
    github.com/jackc/pgpassfile v1.0.0 // indirect
    github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
    github.com/jackc/pgx/v5 v5.7.3 // indirect
    golang.org/x/crypto v0.31.0 // indirect
    golang.org/x/text v0.21.0 // indirect
  )
  ```

- **Key Points**:
  - The `require` block specifies the exact versions of dependencies.
  - Commands like `go run`, `go test`, and `go build` use the versions listed in `go.mod`.
  - The `// indirect` annotation indicates that the package is not directly imported in the codebase.
  - Different projects can use different versions of the same package without conflicts.

3. **The `go.sum` File**

- The `go.sum` file contains cryptographic checksums for the required packages.

- Example:

  ```plaintext
  github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
  github.com/jackc/pgpassfile v1.0.0 h1:/6Hmqy13Ss2zCq62VdNG8tM1wchn8zjSGOBJ6icpsIM=
  github.com/jackc/pgpassfile v1.0.0/go.mod h1:CEx0iS5ambNFdcRtxPj5JhEz+xB6uRky5eyVu/W2HEg=
  github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 h1:iCEnooe7UlwOQYpKFhBabPMi4aNAfoODPEFNiAnClxo=
  ...
  ```

- **Purpose**:
  - Ensures the integrity of downloaded packages.
  - Verifies that the content of dependencies matches the checksums in `go.sum`.

4. **Verifying and Downloading Dependencies**

- Use `go mod verify` to ensure the downloaded packages match the checksums in `go.sum`.

  ```bash
  $ go mod verify
  all modules verified
  ```

- Use `go mod download` to download all dependencies listed in `go.mod`.

5. **Managing Dependencies**

- Use `go get -u` to upgrade to the latest minor or patch release.

  ```bash
  $ go get -u github.com/foo/bar
  ```

- Use `go get` with the `@version` suffix to upgrade to a specific version.

  ```bash
  $ go get github.com/foo/bar@v1.2.3
  ```

- Use `go get` with the `@none` suffix to remove a specific package.

  ```bash
  $ go get github.com/foo/bar@none
  ```

- Use `go mod tidy` to automatically remove any unused packages from `go.mod` and `go.sum` files.
  ```bash
  $ go mod tidy
  ```

### Key Takeaways

- **`go.mod`** specifies the module path, Go version, and dependencies.
- **`go.sum`** ensures the integrity of dependencies using cryptographic checksums.
- Use **`go mod verify`** to verify the integrity of downloaded packages.
- Use **`go mod download`** to download exact versions of dependencies.
- Commands like **`go run`**, **`go test`**, and **`go build`** use the versions listed in `go.mod`.
- Upgrade dependencies with **`go get -u`** or **`go get @version`**.
- Remove unused dependencies with **`go get @none`** or **`go mod tidy`**.
- Go modules enable **reproducible builds** by ensuring consistent dependency versions across environments.

## 4.4 Creating a database connection pool

1. **Using `sql.Open()`**

- **Syntax**:

  ```go
  db, err := sql.Open("pgx", "postgres://username:password@host:port/dbname?param1=value1&param2=value2")
  if err != nil {
      ...
  }
  ```

- **Parameters**:

  - `"pgx"`: The driver name.
  - `"postgres://username:password@host:port/dbname?param1=value1&param2=value2"`: The **Data Source Name (DSN)**.
    - If you are using MYSQL, `parseTime=true`: Converts SQL `TIME` and `DATE` fields to Go `time.Time`.

- **Key Points**:
  - `sql.DB` is a **connection pool**, not a single connection.
    - Go manages connections automatically (opens/closes as needed).
  - The pool is **safe for concurrent access**.
  - The connection pool is intended to be long-lived.
    - Initialize the pool in `main()` and pass it to handlers (do not open/close in handlers).

2. **Implementation**

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

**Key Points**

- **Blank Identifier (`_`) for Driver Import**:

  - The `github.com/jackc/pgx/v5/stdlib` package is imported with a blank identifier (`_`) because the driver's `init()` function must run to register itself with `database/sql`.
  - The Go compiler raises an error if the package is imported but not used directly.

- **`db.Ping()`**:

  - Verifies the connection by creating a connection to the database.
  - Ensures the DSN is correct and the database is accessible.

- **`defer db.Close()`**:
  - Ensures the connection pool is closed when the application exits.
  - While not strictly necessary in this example (since `os.Exit()` or `Ctrl+C` skips deferred functions), it's a good habit for future extensibility (e.g., graceful shutdown).

3. **Testing the Connection**

```bash
$ go run ./cmd/web
time=2025-03-23T17:32:08.513+08:00 level=INFO msg="starting server" addr=:4000
```

4. **Tidying `go.mod`**

- Run `go mod tidy` to tidy your `go.mod` file and remove unnecessary `// indirect` annotations.

### Gotch in `pgx`

- Download the latest version of `pgx` with:

  ```bash
  go get github.com/jackc/pgx/v5
  ```

- `go.mod` file after `go get`:

  ```plaintext
  module snippetbox.libra.dev

  go 1.24.1

  require (
    github.com/jackc/pgpassfile v1.0.0 // indirect
    github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
    github.com/jackc/pgx/v5 v5.7.3 // indirect
    golang.org/x/crypto v0.31.0 // indirect
    golang.org/x/text v0.21.0 // indirect
  )
  ```

- But `go run` outputs the following error:

  ```bash
  $ go run ./cmd/web
  missing go.sum entry for module providing package github.com/jackc/puddle/v2 (imported by github.com/jackc/pgx/v5/pgxpool); to add:
          go get github.com/jackc/pgx/v5/pgxpool@v5.7.3
  ```

- **Solutions**:
  ```bash
  go mod tidy
  ```
  or
  ```bash
  go get github.com/jackc/pgx/v5/pgxpool@v5.7.3
  ```

### Key Takeaways

- Use `sql.Open()` to initialize a **connection pool** (`sql.DB`).
- The connection pool is **long-lived** and **safe for concurrent access**.
- Initialize the pool in `main()` and pass it to handlers.
- Import the driver with a **blank identifier** (`_`) to ensure its `init()` function runs.
- Use `db.Ping()` to verify the connection.
- Always close the pool with `defer db.Close()`.
- Use `go mod tidy` to clean up `go.mod` and remove unnecessary `// indirect` annotations.

## 4.5 Designing a database model

1. **Creating the Database Model**

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

2. **Using the `SnippetModel` in the Application**

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

3. **Benefits of This Structure**

- **Separation of Concerns**:

  - Database logic is decoupled from HTTP handlers.
  - Handlers focus on HTTP tasks (e.g., request validation, response writing).

- **Encapsulation**:

  - The `SnippetModel` encapsulates database operations in a single, reusable object.

- **Testability**:

  - The model can be easily mocked for unit testing.

- **Runtime Flexibility**:

  - We have total control over which database is used at runtime, just by using the -dsn command-line flag.

### Key Takeaways

- Use a **database model** to encapsulate database logic and separate concerns.
- Define a `Snippet` struct to represent individual snippets.
- Create a `SnippetModel` struct with methods for database operations (`Insert`, `Get`, `Latest`).
- Inject the `SnippetModel` into the `application` struct for use in handlers.
- Benefits include **separation of concerns**, **encapsulation**, **testability**, and **runtime flexibility**.

## 4.6 Executing SQL statements

**NOTE: This chapter uses MYSQL as the database, and the identical postgresql example is in the end.**

1. **Executing SQL Queries in Go**

- Methods for Executing Queries

  - **`DB.Query()`**: Used for `SELECT` queries that return multiple rows.
  - **`DB.QueryRow()`**: Used for `SELECT` queries that return a single row.
  - **`DB.Exec()`**: Used for statements that don’t return rows (e.g., `INSERT`, `DELETE`).

- Using `DB.Exec()` to insert a new record:

  - **Syntax**:

    ```go
    result, err := db.Exec(query, args...)
    ```

  - **Returns**:

    - A `sql.Result` object containing:
      - `LastInsertId()`: The ID (an `int64`) of the newly inserted row.
      - `RowsAffected()`: The number of rows (an `int64`) affected by the query.

  - **Note**:
    - Not all drivers and databases support the `LastInsertId()` and `RowsAffected()` methods. For example, `LastInsertId()` is not supported by PostgreSQL.
    - it is common to ignore the `sql.Result` return value if you don’t need it.
      ```go
      _, err := db.Exec(query, args...)
      ```

2. **Implementation**

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

- Testing: Use `curl` to make a `POST /snippet/create` request:

  ```bash
  $ curl -iL -d "" http://localhost:4000/snippet/create
  HTTP/1.1 303 See Other
  Location: /snippet/view/4
  Date: Wed, 18 Mar 2024 11:29:23 GMT
  Content-Length: 0

  HTTP/1.1 200 OK
  Date: Wed, 18 Mar 2024 11:29:23 GMT
  Content-Length: 39
  Content-Type: text/plain; charset=utf-8

  Display a specific snippet with ID 4...
  ```

3. **Placeholder Parameters**

- **Purpose**: Prevent SQL injection by separating SQL code from data.
- **Syntax**:
  - MySQL, SQL Server, SQLite: `?`
  - PostgreSQL: `$1`, `$2`, etc.

4. **How `DB.Exec()` Works**

- Creates a prepared statement on the database.
- Passes parameter values to the database for execution.
- Closes (or deallocates) the prepared statement.

### Key Takeaways

**(MYSQL as the database)**

- Use **`DB.QueryRow()`** for single-row `SELECT` queries.
- Use **`DB.Query()`** for multi-row `SELECT` queries.
- Use **`DB.Exec()`** for `INSERT`, `UPDATE`, and `DELETE` queries.
- Use **placeholder parameters** (`?`) to safely insert untrusted data.
- **`sql.Result`** provides methods like `LastInsertId()` and `RowsAffected()`.
- Retrieve the ID of a newly inserted record with **`LastInsertId()`**.
- Redirect the user to the new snippet's page after insertion.

---

### **POSTGRESQL IMPLEMENTATION**

**File: `internal/models/snippets.go`**

```go
package models

...

type SnippetModel struct {
    DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `
        INSERT INTO snippets (title, content, created, expires)
        VALUES($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1 day' * $3)
        RETURNING id
    `

	var id int
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

...
```

- `CURRENT_TIMESTAMP` is the function in PostgreSQL to get the current timestamp.
- `CURRENT_TIMESTAMP + INTERVAL '1 day' * $3` is used to calculate the expiration time.

- PostgreSQL uses `$1, $2, $3, ...` instead of `?` as placeholders.

- PostgreSQL supports the `RETURNING` clause to return field values (e.g., `id`) of the inserted row.
- `QueryRow` and `Scan` are used to fetch the returned `id`.
- To store the returned `id`, pass a pointer.

## 4.7 Single-record SQL queries

**NOTE: This chapter uses MYSQL as the database, and the identical postgresql example is in the end.**

1. **Implementation**

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

- Use `DB.QueryRow()` to execute the query and retrieve a single row.
- Use `row.Scan()` to copy the row data into a `Snippet` struct.
- Handle the `sql.ErrNoRows` error by returning a custom `ErrNoRecord` error.
- Shorthand single-record queries:

  ```go
  ...

  var s Snippet
  err := m.DB.QueryRow("SELECT ...", id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

  ...
  ```

2. `row.Scan()` under the hood

- Behind the scenes of `rows.Scan()` your driver will automatically convert the raw output from the SQL database to the required native Go types.

- MYSQL (`go-sql-driver`)

  - In MYSQL, we need to use the `parseTime=true` parameter to force it to convert `TIME` and `DATE` fields to `time.Time`. Otherwise it returns these as `[]byte` objects.

  - `CHAR`, `VARCHAR`, `TEXT` → `string`.
  - `BOOLEAN` → `bool`.
  - `INT` → `int`; `BIGINT` → `int64`.
  - `DECIMAL`, `NUMERIC` → `float`.
  - `TIME`, `DATE`, `TIMESTAMP` → `time.Time`.

- PostgreSQL (`pgx`)
  - `CHAR`, `VARCHAR`, `TEXT` → `string`.
  - `BOOLEAN` → `bool`.
  - `INT` → `int`; `SMALLINT` → `int16`;`BIGINT` → `int64`.
  - `REAL` → `float32`; `DOUBLE PRECISION` → `float64`; `DECIMAL`, `NUMERIC` → `float64`
  - `TIME`, `DATE`, `TIMESTAMP`, `TIMESTAMPTZ` → `time.Time`.
  - `BYTEA` → `[]byte`.
  - `JSON`, `JSONB` → `[]byte`.
  - `INT[]`, `TEXT[]`, etc. → `[]int`, `[]string`, etc.

3. **Why Use `ErrNoRecord`?**

- Encapsulate the model and avoid exposing datastore-specific errors (e.g., `sql.ErrNoRows`) to handlers.

4. Error Handling with `errors.Is()`

- Use `errors.Is()` to check for specific errors (e.g., `sql.ErrNoRows`).
- This is safer than using `==` because it handles wrapped errors introduced in Go 1.13.

### Key Takeaways

- Use **`DB.QueryRow()`** for single-row `SELECT` queries.
- Use **`row.Scan()`** to copy row data into a struct.
- Handle **`sql.ErrNoRows`** by returning a custom `ErrNoRecord` error.
- Use **`errors.Is()`** to check for specific errors (e.g., `ErrNoRecord`).
- Encapsulate the model to avoid exposing datastore-specific errors to handlers.

---

### **POSTGRESQL IMPLEMENTATION**

**File: `internal/models/snippets.go`**

```go
...

stmt := `SELECT id, title, content, created, expires FROM snippets
    	WHERE expires > CURRENT_TIMESTAMP AND id = $1
	`

...
```

## 4.8 Multiple-record SQL queries

**NOTE: This chapter uses MYSQL as the database, and the identical postgresql example is in the end.**

**Implementation**

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

**Key Points**

- Use `DB.Query()` to execute the query and retrieve a result set.

- **Closing the Result Set**

  - As long as a resultset is open it will keep the underlying database connection open.
  - Always close the result set with `defer rows.Close()` to free up the database connection.
  - Failure to close the result set can exhaust the connection pool.

- **Iterating Over the Result Set**

  - Use `rows.Next()` to iterate over the rows in the result set.
  - Use `rows.Scan()` to copy row data into a struct.

- **Error Handling**
  - Check for errors after iteration using `rows.Err()`.

### Key Takeaways

- Use **`DB.Query()`** for `SELECT` queries that return multiple rows.
- Iterate over the result set with **`rows.Next()`** and **`rows.Scan()`**.
- Always close the result set with **`defer rows.Close()`**.
- Always call **`rows.Err()`** to retrieve any error that was encountered during the iteration after the `rows.Next()` loop.

---

### **POSTGRESQL IMPLEMENTATION**

**File: `internal/models/snippets.go`**

```go
...

stmt := `SELECT id, title, content, created, expires FROM snippets
    	WHERE expires > CURRENT_TIMESTAMP ORDER BY id DESC LIMIT 10
	`

...
```

## 4.9 Transactions and other details

### The `database/sql` Package

- **Portability**: The `database/sql` package allows Go code to work with any SQL database (e.g., MySQL, PostgreSQL, SQLite).

- **Driver-Specific Quirks**: While the package standardizes interactions, drivers and databases may have unique behaviors. Always review the driver documentation.

### Verbosity in Go SQL Code

- Go's SQL code can feel verbose compared to ORMs in languages like Ruby, Python, or PHP.

- **Advantages**:

  - **Transparency**: Code is explicit and non-magical.
  - **Control**: Developers have full control over SQL queries and behavior.

- **Tools**:
  - Use libraries like [`jmoiron/sqlx`](https://github.com/jmoiron/sqlx) or [`blockloop/scan`](https://github.com/blockloop/scan) to reduce verbosity.

### Managing `NULL` Values

- Go does not handle `NULL` values in database records well by default.

- **Problem**:

  - Scanning a `NULL` value into a `string` results in an error:
    ```bash
    sql: Scan error on column index 1: unsupported Scan, storing driver.Value type <nil> into type *string
    ```

- **Solution**:

  - Use `sql.NullString` (or similar types like `sql.NullInt64`, `sql.NullBool`, etc.) for nullable fields.
  - Alternatively, avoid `NULL` values by setting `NOT NULL` constraints and providing default values on all database columns.

- **Example**:

  ```sql
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

  ```go
  type Book struct {
      Isbn  string
      Title  sql.NullString
      Author sql.NullString
      Price  sql.NullFloat64
  }

  // 1. Connect to the database
  // 2. Query rows
  // 3. Scan rows into Book structs
  ...

  for _, bk := range bks {
    var price string
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

  Full example [here](https://gist.github.com/alexedwards/dc3145c8e2e6d2fd6cd9).

### Working with Transactions

- **Purpose**: Ensure multiple SQL statements use the same database connection and execute atomically.

- **Pattern**:

  - Use `DB.Begin()` to start a transaction.
  - Use `tx.Exec()`, `tx.Query()`, or `tx.QueryRow()` to execute statements within the transaction.
  - Commit the transaction with `tx.Commit()` or roll back with `tx.Rollback()`.

- **Example**:

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

- **Key Points**:
  - Always call `Rollback()` or `Commit()` to release the connection.
  - Use `defer tx.Rollback()` to ensure cleanup in case of errors.

### Prepared Statements

- **Purpose**: Improve performance by reusing complex or repeated SQL statements.

- **How It Works**:

  - Use `DB.Prepare()` to create a prepared statement.
  - Execute the statement with `stmt.Exec()`, `stmt.Query()`, or `stmt.QueryRow()`.

- **Example**:

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

- **Considerations**:
  - Prepared statements are tied to a specific database connection.
    - The `sql.Stmt` object remembers which connection in the pool was used.
  - Under heavy load, statements may be re-prepared on different connections, reducing performance gains.
  - Complexity increases, so use prepared statements only when performance benefits are significant.

# Chapter 5: Dynamic HTML templates

## 5.1 Displaying dynamic data

1. **Updating snippetView handler**

**File: `cmd/web/templates.go`**

```go
package main

import "snippetbox.libra.dev/internal/models"

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
    Snippet models.Snippet
}
```

**File: `cmd/web/handlers.go`**

```go
package main

...

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            http.NotFound(w, r)
        } else {
            app.serverError(w, r, err)
        }
        return
    }

    // Initialize a slice containing the paths to the view.tmpl file,
    // plus the base layout and navigation partial that we made earlier.
    files := []string{
        "./ui/html/base.tmpl",
        "./ui/html/partials/nav.tmpl",
        "./ui/html/pages/view.tmpl",
    }

    // Parse the template files...
    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Create an instance of a templateData struct holding the snippet data.
    data := templateData{
        Snippet: snippet,
    }

    // Pass in the templateData struct when executing the template.
    err = ts.ExecuteTemplate(w, "base", data)
    if err != nil {
        app.serverError(w, r, err)
    }
}

...
```

2. **Creating view template**

**File: `ui/html/pages/view.tmpl`**

```html
{{define "title"}}Snippet #{{.Snippet.ID}}{{end}} {{define "main"}}
<div class="snippet">
  <div class="metadata">
    <strong>{{.Snippet.Title}}</strong>
    <span>#{{.Snippet.ID}}</span>
  </div>
  <pre><code>{{.Snippet.Content}}</code></pre>
  <div class="metadata">
    <time>Created: {{.Snippet.Created}}</time>
    <time>Expires: {{.Snippet.Expires}}</time>
  </div>
</div>
{{end}}
```

**Key Points**:

- Any data that you pass as the final parameter to `ts.ExecuteTemplate()` is represented within your HTML templates by the `.` character (referred to as dot).

- `html/template` package allows you to pass in only one item of dynamic data when rendering a template.
  - **Solution**: Create `templateData` struct to hold multiple dynamic data items.

3. **Key Features of `html/template`**

- **Dynamic content escaping**:

  - Automatically escapes any data that is yielded between `{{ }}` tags.
  - Prevents XSS attacks by escaping dynamic content.
  - Context-aware escaping for HTML, CSS, JS, and URIs.
  - Example:
    ```html
    <span>{{"<script>alert('xss attack')</script>"}}</span>
    ```
    It would be rendered harmlessly as:
    ```html
    <span>&lt;script&gt;alert(&#39;xss attack&#39;)&lt;/script&gt;</span>
    ```

- **Nested templates**:

  - Must explicitly pass dot when invoking one template from another template:
    ```html
    {{template "main" .}}
    {{block "sidebar" .}}{{end}}
    ```

- **Method calling**:

  - Can call methods on exported types:
    `{{.Snippet.Created.Weekday}}`
  - Pass parameters with spaces (no parentheses):
    `{{.Snippet.Created.AddDate 0 6 0}}`

- **HTML comments**:
  - All HTML comments are stripped for security, including conditional comments.

**NOTE**:

PostgreSQL vs. MySQL: `\n` Behavior:

| Behavior                 | PostgreSQL                        | MySQL                                             |
|--------------------------|-----------------------------------|---------------------------------------------------|
| **Default Parsing**      | **Treats `\n` as literal text**.  | **Parses `\n` as a line break**.                  |
| **Store Actual Newline** | Use `E''` prefix (e.g., `E'\n'`). | Default behavior.                                 |
| **Store Literal `\n`**   | Default behavior.                 | Escape as `\\n` or enable `NO_BACKSLASH_ESCAPES`. |
| **Example**              | `'Line1\nLine2'` → Stores `\n`.   | `'Line1\nLine2'` → Stores newline.                |

### Key Takeaways

- Use `html/template` for automatic XSS protection.
- Can only pass single data item to templates (use wrapper structs for multiple items).
- Access struct fields with dot notation (`{{.Field}}`).
- Always pipeline dot (`{{template "name" .}}`) when invoking nested templates.
- `html/template` automatically escapes any data that is yielded between `{{ }}` tags.
- HTML comments are removed from templates.
- Can call methods on template variables with proper syntax.

## 5.2 Template actions and functions

1. **Actions**

| Action | Description |
|--------|-------------|
| `{{if .Foo}} C1 {{else}} C2 {{end}}`    | Renders C1 if `.Foo` not empty, else C2 |
| `{{with .Foo}} C1 {{else}} C2 {{end}}`  | Sets dot to `.Foo` and renders C1 if not empty, else C2 |
| `{{range .Foo}} C1 {{else}} C2 {{end}}` | Loops over `.Foo` (array/slice/map/channel), renders C1 for each element, else C2 |

**Note**
- `{{else}}` clause is optional.
- Empty values: `false`, `0`, `nil`, zero-length collections.
- `with` and `range` change the value of dot.

2. **Functions**

| Function | Description |
|----------|-------------|
| `{{eq .Foo .Bar}}`             | True if `.Foo` equals `.Bar` |
| `{{ne .Foo .Bar}}`             | True if `.Foo` not equal `.Bar` |
| `{{not .Foo}}`                 | Boolean negation of `.Foo` |
| `{{or .Foo .Bar}}`             | Yields `.Foo` if not empty, else `.Bar` |
| `{{index .Foo i}}`             | Value of `.Foo` at index `i` (map/slice/array) |
| `{{printf "%s-%s" .Foo .Bar}}` | Formatted string (like `fmt.Sprintf`) |
| `{{len .Foo}}`                 | Length of `.Foo` as integer |
| `{{$bar := len .Foo}}`         | Assign length to _template variable_ `$bar` |

**Note**

- **Combining functions**:
  `{{if (gt (len .Foo) 99)}} C1 {{end}}`
  `{{if (and (eq .Foo 1) (le .Bar 20))}} C1 {{end}}`

- **Loop control**:

  Control loops with `break` and `continue`.

  ```html
  {{range .Foo}}
    // Skip this iteration if the .ID value equals 99.
    {{if eq .ID 99}}
        {{continue}}
    {{end}}
    // ...
  {{end}}
  ```
  ```html
  {{range .Foo}}
    // End the loop if the .ID value equals 99.
    {{if eq .ID 99}}
        {{break}}
    {{end}}
    // ...
  {{end}}
  ```

3. **Implementation**

**File: `ui/html/pages/view.tmpl`**

```html
{{define "title"}}Snippet #{{.Snippet.ID}}{{end}}

{{define "main"}}
    {{with .Snippet}}
    <div class='snippet'>
        <div class='metadata'>
            <strong>{{.Title}}</strong>
            <span>#{{.ID}}</span>
        </div>
        <pre><code>{{.Content}}</code></pre>
        <div class='metadata'>
            <time>Created: {{.Created}}</time>
            <time>Expires: {{.Expires}}</time>
        </div>
    </div>
    {{end}}
{{end}}
```

**File: `cmd/web/templates.go`**

```go
package main

import "snippetbox.libra.dev/internal/models"

// Include a Snippets field in the templateData struct.
type templateData struct {
    Snippet  models.Snippet
    Snippets []models.Snippet
}
```

**File: `cmd/web/handlers.go`**

```go
package main

...

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")
    
    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    files := []string{
        "./ui/html/base.tmpl",
        "./ui/html/partials/nav.tmpl",
        "./ui/html/pages/home.tmpl",
    }

    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Create an instance of a templateData struct holding the slice of
    // snippets.
    data := templateData{
        Snippets: snippets,
    }

    // Pass in the templateData struct when executing the template.
    err = ts.ExecuteTemplate(w, "base", data)
    if err != nil {
        app.serverError(w, r, err)
    }
}

...
```

**File: `ui/html/pages/home.tmpl`**

```html
{{define "title"}}Home{{end}}

{{define "main"}}
    <h2>Latest Snippets</h2>
    {{if .Snippets}}
     <table>
        <tr>
            <th>Title</th>
            <th>Created</th>
            <th>ID</th>
        </tr>
        {{range .Snippets}}
        <tr>
            <td><a href='/snippet/view/{{.ID}}'>{{.Title}}</a></td>
            <td>{{.Created}}</td>
            <td>#{{.ID}}</td>
        </tr>
        {{end}}
    </table>
    {{else}}
        <p>There's nothing to see here... yet!</p>
    {{end}}
{{end}}
```

### Key Takeaways

- Use `if` for conditional rendering.
- `with` changes dot context for its block.
- `range` iterates over slices/arrays/maps.
- Template functions enable complex logic.
- Combine functions with parentheses.
- Control loops with `break` and `continue`.
- Variables can be declared with `$var := value`.

## 5.3 Caching templates

1. **Optimization Goals**
- Avoid repeatedly parsing template files by implementing an in-memory cache.
- Reduce code duplication in handlers with a helper function.

2. **Template Cache Implementation**

**File: `cmd/web/templates.go`**

```go
package main

import (
    "html/template" // New import
    "path/filepath" // New import

    "snippetbox.libra.dev/internal/models"
)

...

func newTemplateCache() (map[string]*template.Template, error) {
    // Initialize a new map to act as the cache.
    cache := map[string]*template.Template{}

    // Use the filepath.Glob() function to get a slice of all filepaths that
    // match the pattern "./ui/html/pages/*.tmpl". This will essentially gives
    // us a slice of all the filepaths for our application 'page' templates
    // like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
    pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
    if err != nil {
        return nil, err
    }

    // Loop through the page filepaths one-by-one.
    for _, page := range pages {
        // Extract the file name (like 'home.tmpl') from the full filepath
        // and assign it to the name variable.
        name := filepath.Base(page)

        // Create a slice containing the filepaths for our base template, any
        // partials and the page.
        files := []string{
            "./ui/html/base.tmpl",
            "./ui/html/partials/nav.tmpl",
            page,
        }

        // Parse the files into a template set.
        ts, err := template.ParseFiles(files...)
        if err != nil {
            return nil, err
        }

        // Add the template set to the map, using the name of the page
        // (like 'home.tmpl') as the key.
        cache[name] = ts
    }

    // Return the map.
    return cache, nil
}
```

**Note**:
- `filepath.Glob()` returns a slice of strings representing all file paths that match the pattern.
- `filepath.Base()` returns the last element of the path.

**File: `cmd/web/main.go`**

```go
package main

import (
    "database/sql"
    "flag"
    "html/template" // New import
    "log/slog"
    "net/http"
    "os"

    "snippetbox.alexedwards.net/internal/models"

    _ "github.com/jackc/pgx/v5/stdlib"
)

// Add a templateCache field to the application struct.
type application struct {
    logger        *slog.Logger
    snippets      *models.SnippetModel
    templateCache map[string]*template.Template
}

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
    dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    db, err := openDB(*dsn)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
    defer db.Close()

    // Initialize a new template cache...
    templateCache, err := newTemplateCache()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    // And add it to the application dependencies.
    app := &application{
        logger:        logger,
        snippets:      &models.SnippetModel{DB: db},
        templateCache: templateCache,
    }

    logger.Info("starting server", "addr", *addr)

    err = http.ListenAndServe(*addr, app.routes())
    logger.Error(err.Error())
    os.Exit(1)
}

...
```

3. **Render Helper Function Implementation**

**File: `cmd/web/helpers.go`**

```go
package main

import (
    "fmt" // New import
    "net/http"
)

...

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
    // Retrieve the appropriate template set from the cache based on the page
    // name (like 'home.tmpl'). If no entry exists in the cache with the
    // provided name, then create a new error and call the serverError() helper
    // method that we made earlier and return.
    ts, ok := app.templateCache[page]
    if !ok {
        err := fmt.Errorf("the template %s does not exist", page)
        app.serverError(w, r, err)
        return
    }

    // Write out the provided HTTP status code ('200 OK', '400 Bad Request' etc).
    w.WriteHeader(status)

    // Execute the template set and write the response body. Again, if there
    // is any error we call the serverError() helper.
    err := ts.ExecuteTemplate(w, "base", data)
    if err != nil {
        app.serverError(w, r, err)
    }
}
```

**File: `cmd/web/handlers.go`**

```go
package main

import (
    "errors"
    "fmt"
    "net/http"
    "strconv"

    "snippetbox.alexedwards.net/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")
    
    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Use the new render helper.
    app.render(w, r, http.StatusOK, "home.tmpl", templateData{
        Snippets: snippets,
    })
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            http.NotFound(w, r)
        } else {
            app.serverError(w, r, err)
        }
        return
    }

    // Use the new render helper.
    app.render(w, r, http.StatusOK, "view.tmpl", templateData{
        Snippet: snippet,
    })
}

...
```

4. **Improve Partial Parsing**

**File: `cmd/web/templates.go`**

```go
package main

...

func newTemplateCache() (map[string]*template.Template, error) {
    cache := map[string]*template.Template{}

    pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
    if err != nil {
        return nil, err
    }

    for _, page := range pages {
        name := filepath.Base(page)

        // Parse the base template file into a template set.
        ts, err := template.ParseFiles("./ui/html/base.tmpl")
        if err != nil {
            return nil, err
        }

        // Call ParseGlob() *on this template set* to add any partials.
        ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
        if err != nil {
            return nil, err
        }

        // Call ParseFiles() *on this template set* to add the  page template.
        ts, err = ts.ParseFiles(page)
        if err != nil {
            return nil, err
        }

        // Add the template set to the map as normal...
        cache[name] = ts
    }

    return cache, nil
}
```

**Note**:
- `template.ParseGlob()` parses all files in a directory that match the pattern.
- `ParseGlob()` and `ParseFiles()` can be called multiple times on the same template set.

### Key Takeaways

- Cache templates at startup to avoid repeated parsing.
  - Store cache in app struct for easy handler access.
- Centralize rendering with a helper method to reduce duplication.
- Automatically include partials using `ParseGlob()`.

## 5.4 Catching runtime errors

1. **Problem**

Let’s add a deliberate error to our template.

**File: `ui/html/pages/view.tmpl`**

```html
{{define "title"}}Snippet #{{.Snippet.ID}}{{end}}

{{define "main"}}
    {{with .Snippet}}
    <div class='snippet'>
        <div class='metadata'>
            <strong>{{.Title}}</strong>
            <span>#{{.ID}}</span>
        </div>
        {{len nil}} <!-- Deliberate error -->
        <pre><code>{{.Content}}</code></pre>
        <div class='metadata'>
            <time>Created: {{.Created}}</time>
            <time>Expires: {{.Expires}}</time>
        </div>
    </div>
    {{end}}
{{end}}
```

Result:

```bash
$ curl -i http://localhost:4000/snippet/view/1
HTTP/1.1 200 OK
Date: Wed, 18 Mar 2024 11:29:23 GMT
Content-Length: 734
Content-Type: text/html; charset=utf-8


<!doctype html>
<html lang='en'>
    <head>
        <meta charset='utf-8'>
        <title>Snippet #1 - Snippetbox</title>
        <link rel='stylesheet' href='/static/css/main.css'>
        <link rel='shortcut icon' href='/static/img/favicon.ico' type='image/x-icon'>
        <link rel='stylesheet' href='https://fonts.googleapis.com/css?family=Ubuntu+Mono:400,700'>
    </head>
    <body>
        <header>
            <h1><a href='/'>Snippetbox</a></h1>
        </header>
        
 <nav>
    <a href='/'>Home</a>
</nav>

        <main>
            
    
    <div class='snippet'>
        <div class='metadata'>
            <strong>An old silent pond</strong>
            <span>#1</span>
        </div>
        Internal Server Error
```

- Our application has thrown an error, but the user has wrongly been sent a `200 OK` response. And even worse, they’ve received a half-complete HTML page.

2. **Solution**

We make a ‘trial’ render by writing the template into a buffer.
- If this fails, we can respond to the user with an error message.
- If it works, we can then write the contents of the buffer to our `http.ResponseWriter`.

**File: `cmd/web/helpers.go`**

```go
package main

import (
    "bytes" // New import
    "fmt"
    "net/http"
)

...

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
    ts, ok := app.templateCache[page]
    if !ok {
        err := fmt.Errorf("the template %s does not exist", page)
        app.serverError(w, r, err)
        return
    }

    // Initialize a new buffer.
    buf := new(bytes.Buffer)

    // Write the template to the buffer, instead of straight to the
    // http.ResponseWriter. If there's an error, call our serverError() helper
    // and then return.
    err := ts.ExecuteTemplate(buf, "base", data)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // If the template is written to the buffer without any errors, we are safe
    // to go ahead and write the HTTP status code to http.ResponseWriter.
    w.WriteHeader(status)

    // Write the contents of the buffer to the http.ResponseWriter. Note: this
    // is another time where we pass our http.ResponseWriter to a function that
    // takes an io.Writer.
    buf.WriteTo(w)
}
```

Result:

```bash
$ curl -i http://localhost:4000/snippet/view/1
HTTP/1.1 500 Internal Server Error
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Wed, 18 Mar 2024 11:29:23 GMT
Content-Length: 22

Internal Server Error
```

### Key Takeaways

- Use `bytes.Buffer` to capture template output and handle errors.


## 5.5 Common dynamic data

There may be common dynamic data that you want to include on more than one webpage. For example:
- the name and profile picture of the current user
- a CSRF token in all pages with forms

Say that we want to include the current year in the footer on every page.

**Implementation**

**File: `cmd/web/templates.go`**

```go
package main

...

// Add a CurrentYear field to the templateData struct.
type templateData struct {
    CurrentYear int
    Snippet     models.Snippet
    Snippets    []models.Snippet
}

...
```

**File: `cmd/web/helpers.go`**

```go
package main

import (
    "bytes"
    "fmt"
    "net/http"
    "time" // New import
)

...

// Create an newTemplateData() helper, which returns a templateData struct 
// initialized with the current year. Note that we're not using the *http.Request 
// parameter here at the moment, but we will do later in the book.
func (app *application) newTemplateData(r *http.Request) templateData {
    return templateData{
        CurrentYear: time.Now().Year(),
    }
}

...
```

**File: `cmd/web/handlers.go`**

```go
package main

...

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")
    
    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Call the newTemplateData() helper to get a templateData struct containing
    // the 'default' data (which for now is just the current year), and add the
    // snippets slice to it.
    data := app.newTemplateData(r)
    data.Snippets = snippets

    // Pass the data to the render() helper as normal.
    app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            http.NotFound(w, r)
        } else {
            app.serverError(w, r, err)
        }
        return
    }

    // And do the same thing again here...
    data := app.newTemplateData(r)
    data.Snippet = snippet

    app.render(w, r, http.StatusOK, "view.tmpl", data)
}

...
```

**File: `ui/html/base.tmpl`**

```html
...

<footer>
    <!-- Update the footer to include the current year -->
    Powered by <a href='https://golang.org/'>Go</a> in {{.CurrentYear}}
</footer>

...
```

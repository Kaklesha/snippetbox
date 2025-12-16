# **CMD/WEB** section info 
## Description
The CMD directory will contain the application-specific code for the executable applications in the project. For now our project will have just one executable application - the web application - which will live under the cmd/web directory.
## NOTES:

### Serve a single
 Sometimes you might want to 
 file from within a handler. For this there’s the
 http.ServeFile() function, which you can use like so:
```
func downloadHandler(w http.ResponseWriter, r *http.Request) {
http.ServeFile(w, r, "./ui/static/file.zip")
}
```
http.ServerFile() does not automatically sanitize the file path.
//you must sanitize the input with "filepath.Clean()"before using it.

### Disabling directory listing
If you want to disable directory listings there are a few different approaches you can take.
The simplest way? Add a blank index.html file to the specific directory that you want to
disable listings for. This will then be served instead of the directory listing, and the user will
get a 200 OK response with no body. If you want to do this for all directories under
./ui/static you can use the command:
```$ find./ui/static -type d -exec touch {}/index.html \;```
A more complicated (but arguably better) solution is to create a custom implementation of
http.FileSystem, and have it return an os.ErrNotExist error for any directories. A full
explanation and sample code can be found in this blog post.

### Environment variables
If you’ve built and deployed web applications before, then you’re probably thinkingwhat
about environment variables? Surely it’s good-practice to store configuration settingsthere?
If you want, you can store your configuration settings in environment variables and access
them directly from your application by using the ```os.Getenv()``` function like so:
``` addr := os.Getenv("SNIPPETBOX_ADDR") ```
But this has some drawbacks compared to using command-line flags. You can’t specify a
default setting (the return value from ```os.Getenv()``` is the empty string if the environment
variable doesn’t exist), you don’t get the ```-help``` functionality that you do with commandline flags, and the return value from ```os.Getenv()``` is always a string — you don’t get
automatic type conversions like you do with ```flag.Int()```, ```flag.Bool()``` and the other
command line flag functions.
Instead, you can get the best of both worlds by passing the environment variable as a
command-line flag when starting the application. For example:
``` $ export SNIPPETBOX_ADDR=":9999"
$ go run ./cmd/web -addr=$SNIPPETBOX_ADDR
2024/03/18 11:29:23 starting server on :9999 
```


### Pre-existing variables
It’s possible to parse command-line flag values into the memory addresses of pre-existing
variables, using ```flag.StringVar()```, ```flag.IntVar()```, ```flag.BoolVar()```, and similar functions
for other types.
These functions are particularly useful if you want to store all your configuration settings in
a single struct. As a rough example:

``` type config struct {
addr string
staticDir string
}
...
var cfg config
flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
flag.Parse()
```
[![progress-banner](https://backend.codecrafters.io/progress/http-server/a1be822e-d652-4b7a-b8ec-39dfed2daa28)](https://app.codecrafters.io/users/armanchhetri?r=2qF)



This is a solution to 
["Build Your Own HTTP server" Challenge](https://app.codecrafters.io/courses/http-server/overview) on codecrafters.io.

They have got fun challenges to practice programming and software engineering.

As a solution to the HTTP server challenge, I tried to make a small version of net/http package of go standard library.

Along the way I learned the following topics:
- HTTP protocol in depth
- Go language features
   - Interfaces 
   - goroutines
   - synchronization
   - File handling
- Web technology
   - The request-response flow
   - HTTP routes and route Multiplexer


It has 4 routes by default.

<span style="background: blue">GET</span> &nbsp;   <span style="background: rgb(208 216 228);color: black"> /  </span>

```Hello There!```


<span style="background: blue">GET</span> &nbsp;   <span style="background: rgb(208 216 228);color: black"> /user-agent </span>

`echoes back whatever is in user-agent headers`

<span style="background: blue">GET</span> &nbsp;   <span style="background: rgb(208 216 228);color: black"> /echo/message </span>

`echoes back the message part`

<span style="background: blue">GET</span> &nbsp;   <span style="background: rgb(208 216 228);color: black"> /files/filename </span>

`Gets the content of the file`

<span style="background: blue">POST</span> &nbsp;   <span style="background: rgb(208 216 228);color: black"> /files/filename </span>

`Writes the post data to the filename`

USAGE
```sh
./your_server.sh --directory <path/to/a/directory>
```
`--directory: Directory path for files/<filename> route` 


Output
```sh
INFO[0000] Registering route /                          
INFO[0000] Registering route /user-agent                
INFO[0000] Registering route /echo/<msg>                
INFO[0000] Registering route /files/<filename>          
INFO[0000] Serving at: 0.0.0.0:4221    
```

Server registers the routes and listens at `0.0.0.0:4221`

One of the most challenging part for me was to continue reading on the port with a fixed sized buffer when further data is expected. It is specially relevant for POST methods that may include a body of size `Content-Size` as specified in a the header. A fixed buffer may not be enough to accumulate the whole data. I used go routine to keep reading the data from the socket and appended to the buffer until there is no data left or the connection is closed.

Another was a custom multiplexer for routes. Currently the app accepts non-consecutive path parameters. For example `/books/<id>`, `/countries/<id>/nepal/<state>` are accepted but not `/vehicles/<car>/<brand>`. I have used a word-based **Trie** for multiplexing alogrithm(character-based could have been more efficient🤔)






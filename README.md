# krugc-api
The Krug compiler server.

## how it works
The compilers stages are cut into routes of a HTTP server. The server
runs locally on port 8080. The 'driver' for krug will send requests to
the server: lex this file, parse these tokens, etc.

There is latency involved here due to having to send packets back and
forth, however the point here is not speed in compilation.

Hopefully this means that tooling is quite easy to achieve. For example
a text-editor could send over source files to lexically analyze via.
the exposed API.

## overview

The Krug programming language:

* compiles to C;
* is garbage collected;
* has type inference;
* has no generics;

## license

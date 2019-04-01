# krug-serv
This repository is the compiler for Krug. Most of the work is done here, the
frontends job is simply to talk to this server. 

You will need both the krug [frontend](//github.com/hugobrains/krug), and
this server - running locally or on the cloud somewhere - for you to 
compile krug programs.

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
The language itself is not so much the focus of this project. You can see
some code examples in the [tests](//github.com/hugobrains/krug/tree/master/tests) directory.

For a brief overview, Krug is:

* compiled (to C language - eventually LLVM);
* statically typed;
* optional? garbage collection;
* has no generics;
* includes simple type-inference;

It's based loosely on Rust, Go, and C.

## license
MIT, see the [LICENSE](/LICENSE) for more information.
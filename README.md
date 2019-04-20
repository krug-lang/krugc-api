# caasper
This repository is the compiler for Krug: Caasper. Most of the actual compilation 
work is done here, the frontends job is simply to communicate with Caasper.

You will need both the krug [frontend](//github.com/hugobrains/krug), and
an instance of Caasper running locally - or on the cloud somewhere - for you to 
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

## notes

- compiles to c (c99)
- uses stdbool, stdio, stdint and stdlib for now.
- compiled into one big c file, compiling into different
  c files is planned for some kind of conditional compilation
- garbage collection is currently unimplemented
- no generics are planned as this is out of the scope for now
- virtual machine backend is a possibility
- perhaps some rust like ownership memory model will be looked into

#### roadmap

- api route for name mangling
- api route for stripping comments out of the source files
- minification needs to be specified in the c code gen route.
- compression on the generated c code (Gzip), this will be done
  when there are test files that are big enough to measure performance.
- syntax for destructuring structures `let { a, b, c } = some_struct;`
- module system

## license
MIT, see the [LICENSE](/LICENSE) for more information.
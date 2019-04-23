# tutorial
This is a simple walk through of the Krug programming language. Given that the language
is in its early stages of development, this document is subject to change and at any point
may be inaccurate.

> Note: anything marked with a question mark (?) is subject 
to change or perhaps still under consideration!

## what is Krug?
Krug is a general purpose programming language, it compiles into C code so it's quite
capable of lower level programming, though it contains a variety of high level constructs
like garbage-collection(?), ownership semantics, type inference, etc.

## key concepts
> This section presumes that you have a solid grasp of programming concepts. This document
is not for new programmers.

There are a few key concepts that differ from languages you may be familiar with like 
Java or C++:

### mutability
Immutability is preferred by default, so something like this:

    // example 1.
    fn decrease(number int) int {
        number -= 1;
        return number;
    }

Is not allowed, and nor is this:

    // example 2.
    let x = 3;
    x = 4;

Parameters and variables are immutable by default. You must specify
the `mut` keyword as follows:

    // example 1.
    fn decrease(mut number int) int {
        number -= 1;
        return number;
    }

    // example 2.
    mut x = 3;
    x = 4;

Note that in `example 2.`, the `let` keyword was  _replaced_ with `mut`. Mutable
variables are technically their own construct.

#### structure fields
Structure fields, however, are mutable. This is intentional as having
mutable fields is messy and doesn't work too well.

    struct Felix {
        age int,
    };

    let f Felix;
    f.age = 3; // OK!

### ownership(?)
Ownership is 'enabled' by default in Krug. To illustrate with an example:

    fn main() int {
        let a = 3;
        let b = a;
        printf("%d\n", a); // ERROR! use of moved value 'a'
    }

There are a few ways to get around this error

#### borrowing
You can borrow (todo mutable borrows) values, like so:

    let a = 3;
    let b = ref!a;      // note: ref!a
                        //       ^^^^
    printf("%d\n", a);  // OK!

#### 'disable' ownership

> todo, better phrashing for this heading

You can specify that a variable does _not_ own it's value using the tilde `~`:

    // file does not own it's value
    let ~file *FILE = fopen("foo.txt", "r");

    // same_file does not own file.
    let same_file = file;

    // same_same_file does not own same_file
    let same_same_file = same_file;
    
    // this is ok!
    fclose(file);

    // ... though obviously still error prone
    fclose(same_same_file);

Note that that `same_file` or `same_same_file` don't own their
values either. The compiler knows this and so there is no ownership
to be moved/transferred.

#### refactor the code!
In the original example from above:

    fn main() int {
        let a = 3;
        let b = a;
        printf("%d\n", a); // ERROR! use of moved value 'a'
    }

This code can be refactored. You could say the program is flawed
as `b` technically does not need to exist:

    fn main() int {
        let a = 3;
        printf("%d\n", a); // OK!
    }

Though when the complexity of a program increases, and this simple
example can expand accross multiple levels of indirection or
relationships with other code, it isn't always clear.

Another way this code can be refactored is by simply using the new
owner of the value:

    fn main() int {
        let a = 3;
        let b = a;
        printf("%d\n", b); // OK!
    }

#### copying/cloning
This isn't implemented yet!

## cheat sheet/reference
Now that the key concepts are covered, the general syntax is described in this
section of the document.

## functions
Functions are denoted with the `fn` keyword. They must specify a return type.

    fn returns_nothing() void {}
    
    fn add(a int, b int) int {
        return a + b;
    }

## variables

    let age = 20;

    mut mutable_var int = 123;
    mutable_var = 1234;

## looping constructs

### `loop`

    loop {
        // forever!
    }

### `while` loop

    mut i = 0;
    while i < 10; i += 1 {
        
    }

### `next`/`continue`/`break`

    loop {
        break;
    }

## `defer`

    let file = fopen("file.txt", "r");
    defer fclose(file);

## control flow

### if/else if/else

    if 1 == 2 {
        // doesn't execute!
    } else {
        // executes.
    }

    if 1 == 2 {

    } else if 1 == 3 {

    } else {

    }

### switch

    let age = 20;
    switch age {
        // TODO!
    }

### labels & `jump`

    $foo;
        jump foo; // TODO jump $foo; ?

## types

### primitive types

## structures

    struct Person {
        name string,
        age int,
    }

### `impl`

    struct Dog {
        name string,
    }
    
    impl Dog {
        // TODO!
    }

## enums

    TODO

## traits

    TODO

## tuples

## function types

## arrays
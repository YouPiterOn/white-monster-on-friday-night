# white monster on friday night

a custom programming language interpreter written in go, featuring a complete toolchain from lexing to bytecode execution

## overview

white monster on friday night (`.wmofn`) is a virtual machine interpreter that compiles source code to bytecode and executes it. the project implements a full compiler pipeline: lexical analysis, parsing, compilation, and virtual machine execution.

## architecture

the project is organized into several key components:

- **lexer** (`internal/lexer/`) - tokenizes source code into a stream of tokens
- **ast parser** (`internal/ast/`) - builds an abstract syntax tree from tokens
- **compiler** (`internal/compiler/`) - generates bytecode instructions from the ast
- **virtual machine** (`internal/vm/`) - executes bytecode instructions

## current capabilities

### language features

- **variables and constants**
  - `var` declarations for mutable variables
  - `const` declarations for immutable constants
  - variable assignment
  - type annotations (`int`, `bool`, `null`)

- **scoping**
  - block scopes with `{ }`
  - local and upvalue (closure) variable access

- **functions**
  - function declarations with parameters and return types
  - function calls with arguments
  - closures with upvalue capture
  - return statements
  - native functions (e.g., `println`)

- **control flow**
  - `if/else` statements with conditional expressions

- **expressions**
  - integer literals
  - boolean literals (`true`, `false`)
  - null literals
  - binary operators: `+`, `-`, `*`, `/`, `==`, `!=`, `>`, `>=`, `<`, `<=`, `&&`, `||`
  - identifier references
  - function call expressions
  - statement expression optimization (pure expressions as statements are optimized away)

- **types**
  - `int` - integer values
  - `bool` - boolean values
  - `null` - null value
  - function types (closures and native functions)

### example

```javascript
var a = 20;
a = 60;

{
  const b = 123;
}

const b: int = 123;

function addToA(other: int): int {
  return (a + 0) + other;
}

function lessThenA(other: int): bool {
  return other < a;
}

1 + 1

println(1)

1 + println(1)

if (lessThenA(69)) {
  return a;
} else {
  return addToA(42);
}
```

## usage

run a `.wmofn` file using the `run` command:

```bash
go run cmd/run/main.go example/helloWorld.wmofn
```

## planned features

- **loops** - implement `for` and `while` loop constructs
- **unary operators** - support unary operators (e.g., `-`, `!`, `++`, `--`)
- **ternary operators** - add conditional expressions (`condition ? true : false`)
- **vm improvements & async** - upgrade the virtual machine with async/await support for concurrent execution
- **embedded interpreter** - compile interpreter to extern-c dll to make it embedable into other projects

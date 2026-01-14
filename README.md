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
  - type annotations (`int`, `float`, `string`, `bool`, `null`, `any`, `[]type` for arrays)

- **scoping**
  - block scopes with `{ }`
  - local and upvalue (closure) variable access

- **functions**
  - function declarations with parameters and return types
  - function calls with arguments
  - vararg support using rest operator (`...`) for variable-length arguments
  - closures with upvalue capture
  - return statements
  - native functions (`println`, `append`)

- **control flow**
  - `if/else` statements with conditional expressions

- **expressions**
  - integer literals
  - float literals
  - string literals
  - boolean literals (`true`, `false`)
  - null literals
  - array literals with `[ ]` syntax
  - array indexing with `[index]` syntax
  - binary operators: `+`, `-`, `*`, `/`, `==`, `!=`, `>`, `>=`, `<`, `<=`, `&&`, `||`
    - arithmetic operators work with `int` and `float` types
    - string concatenation with `+` operator
    - comparison operators work with `int`, `float`, `string`, and `bool` types
  - identifier references
  - function call expressions
  - statement expression optimization (pure expressions as statements are optimized away)

- **types**
  - `int` - integer values
  - `float` - floating-point values
  - `string` - string values
  - `bool` - boolean values
  - `null` - null value
  - `any` - accepts any type (for flexible typing)
  - `[]type` - array types (e.g., `[]int`, `[]string`)
  - function types (closures and native functions)

### example

```javascript
println("Hello, World!");

var a = 20;
a = 60;

{
  const b = 123;
}

const b: int = 123;
const pi: float = 3.14159;
const message: string = "hello";

function addToA(other: int): int {
  return (a + 0) + other;
}

function lessThenA(other: int): bool {
  return other < a;
}

function greet(name: string): string {
  return "hello, " + name;
}

var arr: []int = [1, 2, 3];
var first = arr[0];
arr = append(arr, 4);

1 + 1
3.14 + 2.86
"hello" + " " + "world"

function factorial(n: int): int {
  if(n == 1) {
    return n;
  }
  return n * factorial(n - 1);
}

return factorial(5);
```

## usage

run a `.wmofn` file using the `run` command:

```bash
go run cmd/cli/main.go run example/helloWorld.wmofn
```

or start an interactive repl:

```bash
go run cmd/cli/main.go repl
```

## planned features

- **loops** - implement `for` and `while` loop constructs
- **unary operators** - support unary operators (e.g., `-`, `!`, `++`, `--`)
- **ternary operators** - add conditional expressions (`condition ? true : false`)
- **vm improvements & async** - upgrade the virtual machine with async/await support for concurrent execution
- **embedded interpreter** - compile interpreter to extern-c dll to make it embedable into other projects

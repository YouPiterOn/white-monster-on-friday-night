# white monster on friday night

A custom programming language interpreter written in Go, featuring a complete toolchain from lexing to bytecode execution.

## Overview

white monster on friday night (`.wmofn`) is a virtual machine interpreter that compiles source code to bytecode and executes it. The project implements a full compiler pipeline: lexical analysis, parsing, compilation, and virtual machine execution.

## Architecture

The project is organized into several key components:

- **Lexer** (`internal/lexer/`) - Tokenizes source code into a stream of tokens
- **AST Parser** (`internal/ast/`) - Builds an Abstract Syntax Tree from tokens
- **Compiler** (`internal/compiler/`) - Generates bytecode instructions from the AST
- **Virtual Machine** (`internal/vm/`) - Executes bytecode instructions

## Current Capabilities

### Language Features

- **Variables and Constants**
  - `var` declarations for mutable variables
  - `const` declarations for immutable constants
  - Variable assignment

- **Scoping**
  - Block scopes with `{ }`
  - Local and upvalue (closure) variable access

- **Functions**
  - Function declarations with parameters
  - Function calls with arguments
  - Closures with upvalue capture
  - Return statements

- **Expressions**
  - Numeric literals (integers)
  - Binary operators: `+`, `-`, `*`, `/`
  - Identifier references
  - Function call expressions

### Example

```javascript
var a = 20;
a = 60;

{
  const b = 123;
}

const b = 321;

function addToA(other) {
  return a + other;
}

return addToA(9);
```

## Usage

Run a `.wmofn` file using the `run` command:

```bash
go run cmd/run/main.go example/helloWorld.wmofn
```

## Planned Features

- **Conditional Statements** - add `if/else` statements for control flow
- **Loops** - implement `for` and `while` loop constructs
- **Data Types** - add some data types except int and function
- **Unary Operators** - support unary operators (e.g., `-`, `!`, `++`, `--`)
- **Ternary Operators** - add conditional expressions (`condition ? true : false`)
- **VM Improvements & Async** - upgrade the virtual machine with async/await support for concurrent execution

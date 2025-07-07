# A Look Under the Hood

This section is for anyone who wants to understand how JAWT works internally. I've tried to document the code in a way that explains not just *what* it does, but *why* it does it that way. These are my notes and thoughts on the design and implementation of the toolchain.

## Internal Packages

*   [**Abstract Syntax Tree (AST)**](./ast.md): The data structure that represents your code.
*   [**Build System**](./build.md): The orchestrator for the entire build process.
*   [**Checker**](./checker.md): The semantic analyzer that makes sure your code makes sense.
*   [**Compiler**](./compiler.md): The parser that turns your JML code into an AST.
*   [**Core**](./core.md): The central nervous system of JAWT, holding the context and configuration.
*   [**Diagnostic Reporting**](./diagnostic.md): The system for reporting errors and warnings.
*   [**Emitter**](./emitter.md): The code generator that turns the AST into runnable web code.
*   [**Process Management**](./process.md): The manager for external tools like the TypeScript and Tailwind CSS compilers.

# toy vm

Pure Go implementation of a small virtual machine described in [Gynvael Coldwind's book][book].

## Usage

Build executable:

```console
$ go build
```

The above generates the `toyvm` executable.

There are a few sample programs in [examples](./examples).
To run them, you have to compile them first with Netwide Assembler (nasm).
To do it, install nasm and then run `make`.

Once you have the binaries, you can run a sample program:

```console
$ ./toyvm examples/hello.bin
Hello, World!
```

There are quite a few tests written.
If not for them, I'd have lost my sanity long time ago.
To run the tests:

```console
$ go test ./vm
```

## Instruction set

The full ISA reference lives in [`ISA.md`](./ISA.md).
In that file you'll find documentation for the instruction groups, operand conventions, jump semantics, and machine-code examples in one place.

[book]: https://github.com/gynvael/zrozumiec-programowanie

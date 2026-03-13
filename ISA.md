# ISA

`toyvm` uses a variable-length, little-endian instruction set.
The ISA is based on the VM described in [Gynvael Coldwind's book][book].

## Conventions

| Notation | Meaning                 |
| :------- | :---------------------- |
| `rdst`   | Destination register.   |
| `rsrc`   | Source register.        |
| `imm8`   | 8-bit immediate value.  |
| `imm16`  | 16-bit immediate value. |
| `imm32`  | 32-bit immediate value. |

## Machine model

- The machine has 16 general-purpose 32-bit registers: `R0` through `R15`.
- `PC` is `R14`.
- `SP` is `R15`.
- `FR` is the flag register.
- `ZF` is the zero flag.
- `CF` is the carry/borrow flag used by comparisons.
- Memory is 64 KiB, addressed from `0x0000` to `0xFFFF`.
- Multi-byte immediates are encoded in little-endian order.
- Unless noted otherwise, arithmetic and bitwise operations work on 32-bit values.
- Relative jumps use 16-bit offsets and wrap modulo `2^16`.

## 1. Data copying instructions

| Opcode | Mnemonic | Operands        | Summary                                                                         |
| :----- | :------- | :-------------- | :------------------------------------------------------------------------------ |
| `00`   | `VMOV`   | `rdst`, `rsrc`  | Copy the value of `rsrc` into `rdst`.                                           |
| `01`   | `VSET`   | `rdst`, `imm32` | Set `rdst` to a 32-bit immediate constant.                                      |
| `02`   | `VLD`    | `rdst`, `rsrc`  | Load a 32-bit value from memory at the address stored in `rsrc` into `rdst`.    |
| `03`   | `VST`    | `rdst`, `rsrc`  | Store the 32-bit value from `rsrc` into memory at the address stored in `rdst`. |
| `04`   | `VLDB`   | `rdst`, `rsrc`  | Load one byte from memory at the address stored in `rsrc` into `rdst`.          |
| `05`   | `VSTB`   | `rdst`, `rsrc`  | Store the low 8 bits of `rsrc` into memory at the address stored in `rdst`.     |

### Examples

`VMOV`

```asm
VMOV R2, R5

; machine code:
; 00 02 05
```

`VSET`

```asm
VSET R4, 0x1234

; machine code:
; 01 04 34 12 00 00
```

`VLD`

```asm
VSET R2, 0x1234
VLD  R1, R2

; machine code:
; 01 02 34 12 00 00
; 02 01 02
```

Reading a 32-bit value from address `0x1234` means reading four consecutive bytes:
`0x1234`, `0x1235`, `0x1236`, and `0x1237`.

`VST`

```asm
VSET R9, 0x1234
VSET R5, 0x12345678
VST  R9, R5

; machine code:
; 01 09 34 12 00 00
; 01 05 78 56 34 12
; 03 09 05
```

`VLDB`

```asm
VSET R3, 0x1234
VLDB R1, R3

; machine code:
; 01 03 34 12 00 00
; 04 01 03
```

`VSTB`

```asm
VSET R1, 0x41
VSET R2, 0x1234
VSTB R2, R1

; machine code:
; 01 01 41 00 00 00
; 01 02 34 12 00 00
; 05 02 01
```

## 2. Arithmetic and logic instructions

| Opcode | Mnemonic | Operands       | Summary                                                                                          |
| :----- | :------- | :------------- | :----------------------------------------------------------------------------------------------- | -------- |
| `10`   | `VADD`   | `rdst`, `rsrc` | Add `rsrc` to `rdst`.                                                                            |
| `11`   | `VSUB`   | `rdst`, `rsrc` | Subtract `rsrc` from `rdst`.                                                                     |
| `12`   | `VMUL`   | `rdst`, `rsrc` | Multiply `rdst` by `rsrc`.                                                                       |
| `13`   | `VDIV`   | `rdst`, `rsrc` | Divide `rdst` by `rsrc`. Division by zero raises interrupt `1` (`INT_DIVISION_ERROR`).           |
| `14`   | `VMOD`   | `rdst`, `rsrc` | Replace `rdst` with `rdst % rsrc`. Division by zero raises interrupt `1` (`INT_DIVISION_ERROR`). |
| `15`   | `VOR`    | `rdst`, `rsrc` | Bitwise OR: `rdst                                                                                | = rsrc`. |
| `16`   | `VAND`   | `rdst`, `rsrc` | Bitwise AND: `rdst &= rsrc`.                                                                     |
| `17`   | `VXOR`   | `rdst`, `rsrc` | Bitwise XOR: `rdst ^= rsrc`.                                                                     |
| `18`   | `VNOT`   | `rdst`         | Bitwise NOT: `rdst = ~rdst`.                                                                     |
| `19`   | `VSHL`   | `rdst`, `rsrc` | Shift `rdst` left by the number of bits stored in `rsrc`.                                        |
| `1A`   | `VSHR`   | `rdst`, `rsrc` | Shift `rdst` right by the number of bits stored in `rsrc`.                                       |

The original design treats register values as natural numbers with zero.
In practice, the bitwise operations still make two's-complement style manipulations possible.

### Examples

`VADD`

```asm
VSET R1, 5
VSET R2, 8
VADD R2, R1

; machine code:
; 01 01 05 00 00 00
; 01 02 08 00 00 00
; 10 02 01
```

`VSUB`

```asm
VSUB R1, R1

; machine code:
; 11 01 01
```

`VMUL`

```asm
VMUL R1, R1

; machine code:
; 12 01 01
```

`VDIV`

```asm
VSET R2, 10
VDIV R1, R2

; machine code:
; 01 02 0A 00 00 00
; 13 01 02
```

`VMOD`

```asm
VSET R2, 10
VMOD R1, R2

; machine code:
; 01 02 0A 00 00 00
; 14 01 02
```

`VOR`

```asm
VOR R1, R2

; machine code:
; 15 01 02
```

`VAND`

```asm
VSET R2, 0x0F
VAND R1, R2

; machine code:
; 01 02 0F 00 00 00
; 16 01 02
```

`VXOR`

```asm
VXOR R1, R1

; machine code:
; 17 01 01
```

`VNOT`

```asm
VNOT R5

; machine code:
; 18 05
```

`VSHL`

```asm
VSET R2, 3
VSHL R1, R2

; machine code:
; 01 02 03 00 00 00
; 19 01 02
```

`VSHR`

```asm
VSET R2, 4
VSHR R1, R2

; machine code:
; 01 02 04 00 00 00
; 1A 01 02
```

## 3. Comparison and conditional jump instructions

| Opcode | Mnemonic        | Operands       | Summary                                                                                      |
| :----- | :-------------- | :------------- | :------------------------------------------------------------------------------------------- |
| `20`   | `VCMP`          | `rdst`, `rsrc` | Compare `rdst` and `rsrc`. Internally, this acts like a subtraction that only updates flags. |
| `21`   | `VJZ` / `VJE`   | `imm16`        | Jump if `ZF` is set.                                                                         |
| `22`   | `VJNZ` / `VJNE` | `imm16`        | Jump if `ZF` is clear.                                                                       |
| `23`   | `VJC` / `VJB`   | `imm16`        | Jump if `CF` is set.                                                                         |
| `24`   | `VJNC` / `VJAE` | `imm16`        | Jump if `CF` is clear.                                                                       |
| `25`   | `VJBE`          | `imm16`        | Jump if `CF` or `ZF` is set.                                                                 |
| `26`   | `VJA`           | `imm16`        | Jump if both `CF` and `ZF` are clear.                                                        |

`VCMP` updates flags as follows:

- `ZF` is set when `rdst == rsrc`.
- `CF` is set when `rdst < rsrc`.

All conditional jumps use a relative 16-bit offset.
If the jump instruction starts at address `A`, then the address of the next instruction is `A + 3`.
The target address is computed like this:

```text
target = (A + 3 + imm16) mod 2^16
imm16  = (target - (A + 3)) mod 2^16
```

### Example

If `VJZ` is located at address `0x13` and should jump to `0x30`,
then `imm16 = 0x001A`:

```asm
VCMP R1, R2
VJZ  0x30

; machine code:
; 20 01 02
; 21 1A 00
```

## 4. Stack instructions

| Opcode | Mnemonic | Operands | Summary                                                                     |
| :----- | :------- | :------- | :-------------------------------------------------------------------------- |
| `30`   | `VPUSH`  | `rsrc`   | Decrement `SP` by 4, then push the 32-bit value from `rsrc` onto the stack. |
| `31`   | `VPOP`   | `rdst`   | Pop a 32-bit value from the stack into `rdst`, then increment `SP` by 4.    |

The stack grows toward lower addresses.
`SP` points to the top of the stack.

### Examples

`VPUSH`

```asm
VXOR  R1, R1
VPUSH R1
VPUSH R1

; machine code:
; 17 01 01
; 30 01
; 30 01
```

`VPOP`

```asm
VPOP R5

; machine code:
; 31 05
```

## 5. Unconditional jumps and calls

| Opcode | Mnemonic | Operands | Summary                                                                        |
| :----- | :------- | :------- | :----------------------------------------------------------------------------- |
| `40`   | `VJMP`   | `imm16`  | Unconditional relative jump.                                                   |
| `41`   | `VJMPR`  | `rsrc`   | Unconditional absolute jump to the address stored in `rsrc`.                   |
| `42`   | `VCALL`  | `imm16`  | Push the return address (`PC + 3`), then perform a relative jump.              |
| `43`   | `VCALLR` | `rsrc`   | Push the return address (`PC + 2`), then jump to the address stored in `rsrc`. |
| `44`   | `VRET`   | none     | Pop a return address from the stack and jump to it.                            |

`VRET` is equivalent to:

```asm
VPOP  RTMP
VJMPR RTMP
```

where `RTMP` is any temporary register.

### Examples

`VJMPR`

```asm
VSET  R1, 0x1234
VJMPR R1

; machine code:
; 01 01 34 12 00 00
; 41 01
```

`VCALL` and `VRET`

If `VCALL` is located at address `0x10` and `func` starts at `0x40`,
then the relative offset is `0x002D`:

```asm
VCALL func
...
func:
VSET R1, 0x1234
VRET

; machine code:
; 0x10: 42 2D 00
; ...
; 0x40: 01 01 34 12 00 00
; 0x46: 44
```

## 6. Control, I/O, and interrupt instructions

| Opcode | Mnemonic | Operands        | Summary                                                                                                                            |
| :----- | :------- | :-------------- | :--------------------------------------------------------------------------------------------------------------------------------- |
| `F0`   | `VCRL`   | `imm16`, `rsrc` | Write the value of `rsrc` to control register `imm16`. Invalid control-register numbers raise interrupt `2` (`INT_GENERAL_ERROR`). |
| `F1`   | `VCRS`   | `imm16`, `rdst` | Read control register `imm16` into `rdst`. Invalid control-register numbers raise interrupt `2` (`INT_GENERAL_ERROR`).             |
| `F2`   | `VOUTB`  | `imm8`, `rsrc`  | Send the low byte of `rsrc` to I/O port `imm8`.                                                                                    |
| `F3`   | `VINB`   | `imm8`, `rdst`  | Read one byte from I/O port `imm8` into `rdst`.                                                                                    |
| `F4`   | `VIRET`  | none            | Return from an interrupt by restoring the saved machine state, including `PC`.                                                     |
| `FF`   | `VOFF`   | none            | Power off the VM.                                                                                                                  |

For the instructions in this group,
the assembly syntax is written in immediate-then-register form,
but the machine-code examples below preserve the original encoded byte order from the spec.

### Examples

`VCRL`

```asm
VSET R0, 1
VCRL 0x110, R0

; machine code:
; 01 00 01 00 00 00
; F0 00 10 01
```

`VOUTB`

```asm
VSET  R0, 0x41
VOUTB 0x20, R0

; machine code:
; 01 00 41 00 00 00
; F2 00 20
```

`VINB`

```asm
VINB 0x20, R0

; machine code:
; F3 00 20
```

## Implementation-specific extension

The current Go implementation also defines one non-standard opcode that is not part of the main ISA above.

| Opcode | Mnemonic | Operands | Summary                                                          |
| :----- | :------- | :------- | :--------------------------------------------------------------- |
| `FE`   | `VCRSH`  | none     | Immediately crash the VM. This is a debugging/testing extension. |

[book]: https://github.com/gynvael/zrozumiec-programowanie

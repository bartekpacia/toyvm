%include "vm.inc"

vset r0, 0x48 ; letter 'H'
vset r1, 0x65 ; letter 'e'
vset r2, 0x6c ; letter 'l'
vset r3, 0x6c ; letter 'l'
vset r4, 0x6f ; letter 'o'
vset r5, 0x0a ; a newline (LF)

voutb 0x20, r0
voutb 0x20, r1
voutb 0x20, r2
voutb 0x20, r3
voutb 0x20, r4
voutb 0x20, r5
voff

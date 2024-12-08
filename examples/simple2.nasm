%include "vm.inc"

vset r0, 0x41 ; 0x41 is letter 'A'
voutb 0x20, r0
voff

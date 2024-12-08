%include "vm.inc"

vset r0, 1
vset r1, 0
vset r2, 0x54 ; letter 'T' (because "True")

vcmp r0, r1

vjz if_true
vjmp end_if

if_true:
  voutb 0x20, r2 ; print letter 'T'
  ; vjmp end_if
end_if:

voff

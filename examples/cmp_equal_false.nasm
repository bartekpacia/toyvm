%include "vm.inc"

vset r0, 1
vset r1, 0
vset r2, 0x54 ; letter 'T' (because "True")
vset r3, 0x46 ; letter 'F' (because "False")

vcmp r0, r1

vje if_true
vjmp if_false

if_true:
  voutb 0x20, r2 ; print letter 'T'
  vjmp end_if
if_false:
  voutb 0x20, r3 ; print letter 'F'
  vjmp end_if
end_if:

voff

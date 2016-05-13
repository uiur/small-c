.data
.text
.globl main

sum:
addi $sp, $sp, -20
sw $ra, 4($sp)
sw $fp, 0($sp)
addi $fp, $sp, 16
sw $a1, 0($fp)
sw $a0, -4($fp)
lw $t1, -4($fp)
addi $sp, $sp, -4
sw $t1, 0($sp)
lw $t2, 0($fp)
lw $t1, 0($sp)
addi $sp, $sp, 4
add $t0, $t1, $t2
sw $t0, -8($fp)
lw $v0, -8($fp)
j sum_exit
sum_exit:
lw $fp, 0($sp)
lw $ra, 4($sp)
addi $sp, $sp, 20
jr $ra

main:
addi $sp, $sp, -24
sw $ra, 4($sp)
sw $fp, 0($sp)
addi $fp, $sp, 20
li $t0, 100
sw $t0, 0($fp)
li $t0, 20
sw $t0, -4($fp)
lw $a1, -4($fp)
lw $a0, 0($fp)
jal sum
sw $v0, -8($fp)
lw $t0, -8($fp)
sw $t0, -12($fp)
li $v0, 1
lw $a0, -12($fp)
syscall
main_exit:
lw $fp, 0($sp)
lw $ra, 4($sp)
addi $sp, $sp, 24
jr $ra


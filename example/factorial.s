.data
.text
.globl main

fact:
addi $sp, $sp, -28
sw $ra, 4($sp)
sw $fp, 0($sp)
addi $fp, $sp, 24
sw $a0, 0($fp)
lw $t1, 0($fp)
addi $sp, $sp, -4
sw $t1, 0($sp)
li $t2, 1
lw $t1, 0($sp)
addi $sp, $sp, 4
beq $t1, $t2, beq_true_1
li $t0, 0
j beq_end_1
beq_true_1:
li $t0, 1
beq_end_1:
sw $t0, -4($fp)
lw $t0, -4($fp)
beq $t0, $zero, ir_if_false_1
j true_3
ir_if_false_1:
j false_3
ir_if_end_1:
true_3:
li $t0, 1
sw $t0, -8($fp)
lw $v0, -8($fp)
j fact_exit
j end_3
false_3:
lw $t1, 0($fp)
addi $sp, $sp, -4
sw $t1, 0($sp)
li $t2, 1
lw $t1, 0($sp)
addi $sp, $sp, 4
sub $t0, $t1, $t2
sw $t0, -12($fp)
lw $a0, -12($fp)
jal fact
sw $v0, -16($fp)
lw $t1, 0($fp)
addi $sp, $sp, -4
sw $t1, 0($sp)
lw $t2, -16($fp)
lw $t1, 0($sp)
addi $sp, $sp, 4
mul $t0, $t1, $t2
sw $t0, -8($fp)
lw $v0, -8($fp)
j fact_exit
end_3:
fact_exit:
lw $fp, 0($sp)
lw $ra, 4($sp)
addi $sp, $sp, 28
jr $ra

main:
addi $sp, $sp, -20
sw $ra, 4($sp)
sw $fp, 0($sp)
addi $fp, $sp, 16
li $t0, 4
sw $t0, 0($fp)
lw $a0, 0($fp)
jal fact
sw $v0, -4($fp)
lw $t0, -4($fp)
sw $t0, -8($fp)
li $v0, 1
lw $a0, -8($fp)
syscall
main_exit:
lw $fp, 0($sp)
lw $ra, 4($sp)
addi $sp, $sp, 20
jr $ra

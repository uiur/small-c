  .text
  .globl    main
sum:
  addi $sp, $sp, -8
  sw $ra, 4($sp)
  sw $fp, 0($sp)
  addi $fp, $sp, 8

  add $t0, $a0, $a1
  move $v0, $t0

  lw $fp, 0($sp)
  lw $ra, 4($sp)
  addi $sp, $sp, 8

  jr $ra

fib:
  subu $sp, $sp, 20
  sw $ra, 4($sp)
  sw $fp, 0($sp)

  addiu $fp, $sp, 20
  sw $a0, -4($fp)
  sw $s0, -8($fp)
  sw $s1, -12($fp)

  slti $t0, $a0, 1
  beq $t0, $zero, ifelse

  li $v0, 1
  j fibexit

ifelse:
  lw $t0, -4($fp)
  addi $a0, $t0, -1
  jal fib
  move $s0, $v0

  lw $t0, -4($fp)
  addi $a0, $t0, -2
  jal fib
  move $s1, $v0

  add $v0, $s0, $s1
  j fibexit

ifend:

fibexit:
  lw $s1, -12($fp)
  lw $s0, -8($fp)
  lw $a0, -4($fp)

  lw $fp, 0($sp)
  lw $ra, 4($sp)
  addiu $sp, $sp, 20

  jr $ra

main:
  subu $sp, $sp, 20
  sw $ra, 4($sp)
  sw $fp, 0($sp)
  addiu $fp, $sp, 20

  li $a0, 10
  jal fib

  move $a0, $v0
  li $v0, 1
  syscall

  lw $fp, 0($sp)
  lw $ra, 4($sp)
  addiu $sp, $sp, 20
  jr $ra

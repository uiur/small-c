#! /usr/bin/env racket
#lang racket

;; =========== mips-util.rkt ===========================
(require parser-tools/lex
         (prefix-in : parser-tools/lex-sre)
         parser-tools/yacc)

(define-tokens tokens-with-value
  (NUM CHAR STR ID DIR))

(define-empty-tokens tokens-without-value
  (COLON COMMA NEWLINE
   DOLLAR LPAR RPAR
   EOF))

(define-lex-abbrevs
  (digit            (char-range "0" "9"))
  (digit-non-zero   (char-range "1" "9"))
  (number  (:or "0"
                (:: digit-non-zero
                    (:* digit))))
  (identifier-char (:or (char-range "a" "z")
                        (char-range "A" "Z")
                        "_"))
  (identifier (:: identifier-char
                  (:* (:or identifier-char digit)))))

(define mips-lexer
  (lexer
   ("$"        (token-DOLLAR))
   (":"        (token-COLON))
   ("("        (token-LPAR))
   (")"        (token-RPAR))
   (","        (token-COMMA))
   ((:: (:or "" "+" "-") number) (token-NUM (string->number lexeme)))
   (identifier (token-ID (string->symbol lexeme)))
   ((:: "." identifier) (token-DIR (string->symbol lexeme)))
   ((:: "'" (:or any-char (:: "\\" any-char)) "'")
    (token-CHAR
     (string-ref 
      (read
       (open-input-string
        (string-append "\""
                       (substring lexeme 1 (- (string-length lexeme) 1))
                       "\"")))
      0)))
   ((:: "\"" (:* (:or any-char (:: "\\" any-char))) "\"")
    (token-STR (read (open-input-string lexeme))))
   ("\n" (token-NEWLINE))
   ((:or " " "\t") (mips-lexer input-port))
   ((:: "#" (:* (:~ "\n"))) (mips-lexer input-port))
   ((eof)      (token-EOF))))

(define mem-instrs
  (apply set '(lb lbu ld lh lhu ll lw lwc1 lwl lwr ulh ulhu
               ulw sb sc sd sh sw swc1 sdc1 swl swr ush usw )))

(define instrs
  (set-union
   mem-instrs
   (apply set '(abs add addi addiu addu and andi b bclf bclt beq
                    beqz bge bgeu bgez bgezal bgt bgtu bgtz ble bleu
                    blez blt bltu bltz bltzal bne bnez clo clz div
                    divu j jal jalr jr li lui la move movf movn movt
                    movz mfc0 mfc1 mfhi mflo mthi mtlo mtc0 mtc1 madd
                    maddu msub msubu mul mulo mulou mult multu neg negu
                    nop nor not or ori rem remu rol ror seq sge sgeu
                    sgt sgtu sle sleu slt slti sltiu sltu sne sll sllv
                    sra srav srl srlv sub subu syscall xor xori))))

(define directives
  (apply set '(.align .ascii .asciiz .byte .data .double .extern
               .float .globl .half .kdata .ktext .set .space .text .word
               .rdata .sdata)))

(define mips-parser
  (parser
   (start program)
   (end EOF)
   ;(debug "mips-parser.tbl")
   (suppress)
   (error (lambda (tok-ok? tok-name tok-value)
            (error "parse error:" tok-name tok-value)))
   (tokens tokens-with-value tokens-without-value)
   (grammar
    (program (() '())
             ((line program) (if $1 (cons $1 $2) $2)))
    (line ((instruction NEWLINE) $1)
          ((directive NEWLINE) $1)
          ((label) $1)
          ((NEWLINE) #f))
    (label ((ID COLON) `(#:label ,$1))
           ((ID COLON NEWLINE) `(#:label ,$1)))
    (instruction ((ID  operands-opt)
                  (if (set-member? instrs $1)
                      (cons $1 $2)
                      (error (format "illegal opcode: ~a" $1)))))
    (directive ((DIR operands-opt)
                (if (set-member? directives $1)
                    (cons $1 $2)
                    (error (format "illegal directive: ~a" $1)))))
    (operands-opt (() '())
                 ((operands) $1))
    (operands ((operand) (list $1))
              ((operand COMMA operands) (cons $1 $3)))
    (operand ((NUM)  $1)
             ((ID)   $1)
             ((DOLLAR ID) `($ ,$2))
             ((NUM LPAR operand RPAR) `(,$1 ,$3))
             ((STR) $1)
             ((CHAR) $1)))))

(define (mips-parse-port port)
  (mips-parser (lambda () (mips-lexer port))))

(define (mips-parse-string str)
  (mips-parse-port (open-input-string str)))

(define (mips-parse-file fname)
  (mips-parse-port (open-input-file fname)))

(define (mips-count-file fname)
  (let* ((code (mips-parse-file fname))
         (mems (filter (lambda (line) (set-member? mem-instrs (first line)))
                       code)))
    (display (format "total: ~a~%" (length code)))
    (display (format "memory access: ~a~%" (length mems)))))

;; =========== end of mips-util.rkt ====================

(define spim-command "/home/lab4/umatani/local/bin/spim -file ~a")

(define racket-compiler-module "compiler")
(define racket-compiler-function 'compile)

(define (usage)
  (display "Usage: scc (<option> | <input>)*
Options:
  -e, execute (i.e., compile and then invoke spim)
  -n, count the number of MIPS instructions
  -r <module>:<function>, specify compiler function (default=compiler:compile)
  -c <command>, specify compiler command
  -s <spim>, specify spim command
  -h, print this message
")
  (exit 1))

(define r (foldl
           (lambda (arg r)
             (let ((mode (first r))
                   (ins (second r))
                   (opts (third r)))
               (if (char=? (string-ref arg 0) #\-)
                   (let ((opt (string->symbol
                               (format "~a" (string-ref arg 1)))))
                     (when mode (usage))
                     (case opt
                       ((e n h)
                        `(#f ,ins ,(hash-set opts opt #t)))
                       ((r c s)
                        `(,opt ,ins ,opts))
                       (else (usage))))
                   (if mode
                       `(#f ,ins ,(hash-set opts mode arg))
                       `(#f ,(cons arg ins) ,opts)))))
           `(#f () ,(hasheq))
           (vector->list (current-command-line-arguments))))

(when (first r) (usage))
(define inputs (reverse (second r)))
(define options (third r))

(when (hash-ref options 'h #f) (usage))

(when (= (length inputs) 0)
  (eprintf "no input file specified~%")
  (usage))

(let ((s-option (hash-ref options 's #f)))
  (when s-option
    (set! spim-command (string-append s-option " -f ~a"))))

(define assems
  (map (lambda (input)
         (let ((r (regexp-match #rx"^(.+)\\.sc$" input)))
           (if r
               (format "~a.s" (second r))
               (begin
                 (eprintf "input file must be *.sc~%")
                 (exit 1)))))
       inputs))

(define external-compiler-command (void))

(define (string-split str sep-str)
  (let ((sep (string-ref sep-str 0)))
    (define (aux str-list line lines)
      (cond ((null? str-list) (cons line lines))
            ((char=? (first str-list) sep)
             (aux (rest str-list) '() (cons line lines)))
            (else (aux (rest str-list) (cons (first str-list) line) lines))))
    (let ((str-list (string->list str)))
      (let ((lines (aux str-list '() '())))
        (filter (lambda (s) (> (string-length s) 0))
                (map (lambda (line) (list->string (reverse line)))
                     (reverse lines)))))))

(let ((r-option (hash-ref options 'r #f))
      (c-option (hash-ref options 'c #f)))
  (cond
   ((and r-option c-option)
    (eprintf "-r and -c must be exclusive~%")
    (exit 1))
   (r-option
    (let ((strs (string-split r-option ":")))
      (set! racket-compiler-module (first strs))
      (set! racket-compiler-function (string->symbol (second strs)))))
   (c-option
    (set! external-compiler-command c-option))
   (else
    (set! r-option #t)))

  (define racket-compiler-file (string-append racket-compiler-module ".rkt"))

  (define (with-redirect-output port thunk)
    (let ((stdout (current-output-port)))
      (current-output-port port)
      (with-handlers ((exn:fail?
                       (lambda (e)
                         (current-output-port stdout)
                         (raise e))))
        (thunk)
        (current-output-port stdout))))

  (define (trim-spim-out str)
    (let ((lines (string-split str "\n")))
      (string-join (filter (lambda (line)
                             (not (regexp-match #rx"^Loaded:" line)))
                           lines)
                   "\n")))

  (let ((compile (cond
                  (r-option
                   (unless (file-exists? racket-compiler-file)
                     (error (format "module ~a not found~%"
                                    racket-compiler-file)))
                   (let ((compiler (dynamic-require racket-compiler-file
                                                    racket-compiler-function)))
                     (lambda (p in)
                       (with-redirect-output p (lambda () (compiler in))))))
                  (c-option
                   (unless (file-exists? external-compiler-command)
                     (error (format "command ~a not found~%"
                                    external-compiler-command)))
                   (lambda (p in)
                     (let* ((subp (process/ports
                                   p
                                   (current-input-port)
                                   (current-error-port)
                                   (format "~a ~a"
                                           external-compiler-command
                                           in)))
                            (comm-fun (fifth subp)))
                       (comm-fun 'wait)
                       (let ((exit-status (comm-fun 'exit-code)))
                         (unless (zero? exit-status)
                           (error
                            (format "external compiler error: ~a"
                                    exit-status)))))))
                  (else (error "scc: no such case")))))
    (for-each
     (lambda (input assem)
       (let ((p (open-output-file assem
                                  #:mode 'text
                                  #:exists 'replace)))
         (with-handlers ((exn:fail?
                          (lambda (e)
                            (close-output-port p)
                            (delete-file assem)
                            (eprintf "...compile error in ~a~%" input)
                            ;(eprintf "~a~%" e)
                            (newline (current-error-port)))))
           (compile p input)
           (close-output-port p)
           (when (hash-ref options 'n #f)
             (display (format "[# of instructions in ~a]~%" assem))
             (mips-count-file assem))
           (when (hash-ref options 'e #f)
             (display (format "[output of spim -f ~a]~%" assem))
             (let* ((spim-out (open-output-string))
                    (subp (process/ports spim-out
                                         (current-input-port)
                                         'stdout
                                         (format spim-command assem)))
                    (comm-fun (fifth subp)))
               (comm-fun 'wait)
               (displayln (trim-spim-out (get-output-string spim-out)))
               (close-output-port spim-out))))))
     inputs
     assems)))

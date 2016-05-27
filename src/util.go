package main

func ast(src string) []Statement {
	statements, _ := Parse(src)

	env := &Env{}
	Analyze(statements, env)

	return statements
}

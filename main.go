package main

func main() {
	env := ReadEnv()
	api := NewApiServer(&env)
	api.Start()
}

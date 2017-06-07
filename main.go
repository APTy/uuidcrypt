package main

func main() {
	cli := NewCLI(NewFlagConfig())
	cli.Run()
}

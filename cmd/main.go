package main

func main() {
	app := newApp(loadConfig())
	app.Run()
}

package main

import forum "forum/function"

func main() {
	forum.InitDB()
	forum.Server()
}

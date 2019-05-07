package main

func main() {
	id := SecretId
	id.content = "Hello World!"

	SendWeiXinMessage(id)
}

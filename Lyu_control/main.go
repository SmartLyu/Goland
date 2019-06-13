package main

func main(){
	go StartApi("9999")
	<-ListenSig
}

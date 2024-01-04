package main

func main() {
	address := "0x0000000000000000000000007c4fb3e5ba0a5d80658889715b307e66916f29b2"
	// 去掉前面的零，只保留后面40个字符
	standardAddress := "0x" + address[len(address)-40:]
	println(standardAddress)
}

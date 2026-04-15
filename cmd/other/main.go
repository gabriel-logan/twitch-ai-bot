package main

import (
	"log"
	"time"
	"unsafe"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
)

func main() {
	// Initialize logger - uses log package internally
	config.InitLogger()

	log.Println("=== Basic Types ===")

	log.Println("VAR SIZE ALIGN")

	log.Println("=== Void ===")

	log.Println("Void", unsafe.Sizeof(struct{}{}), unsafe.Alignof(struct{}{}))

	log.Println("=== Byte, Rune ===")
	log.Println("byte", unsafe.Sizeof(byte(0)), unsafe.Alignof(byte(0)))
	log.Println("rune", unsafe.Sizeof(rune(0)), unsafe.Alignof(rune(0)))

	log.Println("=== Boolean ===")

	log.Println("bool", unsafe.Sizeof(bool(false)), unsafe.Alignof(bool(false)))

	log.Println("=== Integers ===")

	log.Println("int8", unsafe.Sizeof(int8(0)), unsafe.Alignof(int8(0)))
	log.Println("int16", unsafe.Sizeof(int16(0)), unsafe.Alignof(int16(0)))
	log.Println("int32", unsafe.Sizeof(int32(0)), unsafe.Alignof(int32(0)))
	log.Println("int64", unsafe.Sizeof(int64(0)), unsafe.Alignof(int64(0)))

	log.Println("=== Unsigned Integers ===")

	log.Println("uint8", unsafe.Sizeof(uint8(0)), unsafe.Alignof(uint8(0)))
	log.Println("uint16", unsafe.Sizeof(uint16(0)), unsafe.Alignof(uint16(0)))
	log.Println("uint32", unsafe.Sizeof(uint32(0)), unsafe.Alignof(uint32(0)))
	log.Println("uint64", unsafe.Sizeof(uint64(0)), unsafe.Alignof(uint64(0)))

	log.Println("=== Platform-Dependent Integers ===")

	log.Println("int", unsafe.Sizeof(int(0)), unsafe.Alignof(int(0)))
	log.Println("uint", unsafe.Sizeof(uint(0)), unsafe.Alignof(uint(0)))
	log.Println("uintptr", unsafe.Sizeof(uintptr(0)), unsafe.Alignof(uintptr(0)))

	log.Println("=== Floats ===")

	// floats
	log.Println("float32", unsafe.Sizeof(float32(0)), unsafe.Alignof(float32(0)))
	log.Println("float64", unsafe.Sizeof(float64(0)), unsafe.Alignof(float64(0)))

	log.Println("=== Complex Numbers ===")

	// complex
	log.Println("complex64", unsafe.Sizeof(complex64(0)), unsafe.Alignof(complex64(0)))
	log.Println("complex128", unsafe.Sizeof(complex128(0)), unsafe.Alignof(complex128(0)))

	log.Println("=== Strings, Slices ===")

	// string / slice
	log.Println("string", unsafe.Sizeof(""), unsafe.Alignof(""))
	log.Println("[]string", unsafe.Sizeof([]string{}), unsafe.Alignof([]string{}))

	log.Println("=== Arrays ===")

	// array
	log.Println("[3]int", unsafe.Sizeof([3]int{}), unsafe.Alignof([3]int{}))

	log.Println("=== Maps ===")

	// map
	log.Println("map[string]int", unsafe.Sizeof(map[string]int{}), unsafe.Alignof(map[string]int{}))

	log.Println("=== Pointers ===")

	// pointer
	var p *int
	log.Println("*int", unsafe.Sizeof(p), unsafe.Alignof(p))

	log.Println("=== Structs ===")

	// struct
	type S struct {
		A int8
		B int64
	}
	log.Println("struct", unsafe.Sizeof(S{}), unsafe.Alignof(S{}))

	log.Println("=== Interfaces ===")

	// interface
	var i interface{}
	log.Println("interface{}", unsafe.Sizeof(i), unsafe.Alignof(i))

	log.Println("=== Functions ===")

	// func
	var fn func()
	log.Println("func()", unsafe.Sizeof(fn), unsafe.Alignof(fn))

	log.Println("=== Channels ===")

	// channel
	var ch chan int
	log.Println("chan int", unsafe.Sizeof(ch), unsafe.Alignof(ch))

	log.Println("=== Duration ===")

	// duration
	log.Println("time.Duration", unsafe.Sizeof(time.Duration(0)), unsafe.Alignof(time.Duration(0)))
}

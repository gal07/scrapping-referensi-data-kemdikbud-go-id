package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	defer TimeTrack(start, "main")
	defer fmt.Println("Selesai")

	var nav string
	fmt.Println("=== Pilih Menu ===")
	fmt.Println("1. Download Data Paud")
	fmt.Println("2. Download Data DIKDAS(SD & SMP)")
	fmt.Println("3. Download Data DIKMEN(SMA & SMK)")
	fmt.Scanf("%s", &nav)

	switch nav {
	case "1":
		fmt.Println("Tunggu Sebentar")
		GetProvince("https://referensi.data.kemdikbud.go.id/pendidikan/paud", "paud")
	case "2":
		fmt.Println("Tunggu Sebentar")
		GetProvince("https://referensi.data.kemdikbud.go.id/pendidikan/dikdas", "dikdas")
	case "3":
		fmt.Println("Tunggu Sebentar")
		GetProvince("https://referensi.data.kemdikbud.go.id/pendidikan/dikmen", "dikmen")
	}

}

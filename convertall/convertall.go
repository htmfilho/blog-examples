package main

import "fmt"

type Converter func(float32) float32

func convert (converter Converter, value float32) float32 {
	return converter(value)
}

func BarToPSI(bar float32) float32 {
	return bar * 14.503773773
}

func PSIToBar(PSI float32) float32 {
	return PSI * 0.0689475729
}

func KPaToPSI(kPa float32) float32 {
	return kPa * 0.1450377377
}

func PSIToKPa(PSI float32) float32 {
	return PSI * 6.8947572932
}

func KPaToBar(kPa float32) float32 {
	PSI := KPaToPSI(kPa)
	return PSIToBar(PSI)
}

func BarToKPa(bar float32) float32 {
	PSI := BarToPSI(bar)
	return PSIToKPa(PSI)
}

func main() {
	fmt.Printf("Bar to PSI: %v -> %v \n", 2.9,       convert(BarToPSI, 2.9))
	fmt.Printf("PSI to Bar: %v -> %v \n", 42.060944, convert(PSIToBar, 42.060944))
	fmt.Printf("kPa to PSI: %v -> %v \n", 240,       convert(KPaToPSI, 240))
	fmt.Printf("PSI to kPa: %v -> %v \n", 35,        convert(PSIToKPa, 35))
	fmt.Printf("kPa to Bar: %v -> %v \n", 240,       convert(KPaToBar, 240))
	fmt.Printf("Bar to kPa: %v -> %v ", 2.9,         convert(BarToKPa, 2.9))
}

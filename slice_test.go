package util

import (
	"log"
	"testing"
)

func BenchmarkDeleteByIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := []string{"中国", "and", "美国", "and", "法国"}
		DeleteByIndex(&slice, 1)
	}
}

func TestDeleteByIndex(t *testing.T) {
	slice := []string{"中国", "and", "美国", "and", "法国"}
	if err := DeleteByIndex(&slice, 1); err != nil {
		log.Println("TestDeleteByIndex error ->", err.Error())
	} else {
		log.Println(slice)
	}
}

func BenchmarkDeleteByValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := []string{"中国", "and", "美国", "and", "法国"}
		DeleteByValue(&slice, "and")
	}
}

func TestDeleteByValue(t *testing.T) {
	slice := []string{"中国", "and", "美国", "and", "法国"}
	if err := DeleteByValue(&slice, "and"); err != nil {
		t.Fatal(err)
	} else {
		log.Println(slice)
	}
}

func BenchmarkInSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := []string{"中国", "and", "美国", "and", "法国"}
		InSlice(slice, "and")
	}
}

func TestInSlice(t *testing.T) {
	slice := []string{"中国", "and", "美国", "and", "法国"}
	if exists, err := InSlice(slice, "and"); err != nil {
		t.Fatal(err)
	} else {
		if exists {
			log.Println("value exists in slice")
		} else {
			log.Println("value not exists in slice")
		}
	}
}

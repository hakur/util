package util

import (
	"log"
	"testing"
)

func TestDeleteByIndex(t *testing.T) {
	slice := []string{"中国", "and", "美国", "and", "法国"}
	if err := DeleteByIndex(&slice, 1); err != nil {
		log.Println("TestDeleteByIndex error ->", err.Error())
	} else {
		log.Println(slice)
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

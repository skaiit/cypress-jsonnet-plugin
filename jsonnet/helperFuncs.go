package main

import (
	"github.com/brianvoe/gofakeit/v7"
)

func CallGoFakeIt(pattern string) (string, error) {
	return gofakeit.Generate(pattern)
}

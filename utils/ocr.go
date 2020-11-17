package main

import (
	"fmt"
	"image"
	"log"
	"os"

	_ "image/png"
)

type red struct {
	r        uint32
	position struct {
		x int
		y int
	}
}

func main() {
	f, err := os.Open("/home/yang/Downloads/screenshot-192-168-1-44-8000-go-src-xueXi-sc-png-1605600429410.png")
	//	f, err := os.Open("../sc.png")
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r == 65535 && g == 0 && b == 0 {
				fmt.Print("x")
			} else {
				fmt.Print(".")
			}
		}

		fmt.Print("\n")
	}
}

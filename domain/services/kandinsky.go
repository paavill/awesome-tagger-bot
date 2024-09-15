package services

import "image"

type Kandinsky interface {
	GenerateImage(query string) (*image.Image, error)
}

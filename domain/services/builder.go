package services

type Builder interface {
	Kandinsky(kandinsky Kandinsky) Builder
	Build() (Services, error)
}

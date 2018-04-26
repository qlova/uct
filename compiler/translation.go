package compiler

type Language int

const (
	English Language = iota
	Maori
	
	Chinese
	All
)

type Translatable [All+1]string

func NoTranslation(name string) Translatable {
	var t Translatable
	for i := range t {
		t[i] = name
	}
	return t
}

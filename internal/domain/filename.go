package domain

type Filename string

func RandomFilename() (Filename, error) {
	f, err := RandomStringURLSafe(42)
	return Filename(f), err
}

package entity

type Phone struct {
	Id        int
	Brand     string
	Model     string
	Year      int
	OS        string
	Processor string
}

type PhoneInputDto struct {
	Brand     string
	Model     string
	Year      int
	OS        string
	Processor string
}

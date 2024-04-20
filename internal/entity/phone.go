package entity

type PhoneBuilder interface {
	Brand(brand string) PhoneBuilder
	Model(model string) PhoneBuilder
	Year(year int) PhoneBuilder
	OS(os string) PhoneBuilder
	Processor(processor string) PhoneBuilder
	Build() Phone
}

type Phone struct {
	Brand     string
	Model     string
	Year      int
	OS        string
	Processor string
}

type ConcretePhoneBuilder struct {
	phone Phone
}

func (pb *ConcretePhoneBuilder) Brand(brand string) PhoneBuilder {
	pb.phone.Brand = brand
	return pb
}

func (pb *ConcretePhoneBuilder) Model(model string) PhoneBuilder {
	pb.phone.Model = model
	return pb
}

func (pb *ConcretePhoneBuilder) Year(year int) PhoneBuilder {
	pb.phone.Year = year
	return pb
}

func (pb *ConcretePhoneBuilder) OS(os string) PhoneBuilder {
	pb.phone.OS = os
	return pb
}

func (pb *ConcretePhoneBuilder) Processor(processor string) PhoneBuilder {
	pb.phone.Processor = processor
	return pb
}

func (pb *ConcretePhoneBuilder) Build() Phone {
	return pb.phone
}

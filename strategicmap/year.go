package strategicmap

type YearInterface interface {
	Get() int
	Increment()
}

type YearSimpleProvider struct {
	Year int
}

func (simpleProvider YearSimpleProvider) Get() int {
	return simpleProvider.Year
}

func (simpleProvider *YearSimpleProvider) Increment() {
	simpleProvider.Year++
}

var simpleYearInterfaceChecker YearInterface = &YearSimpleProvider{}

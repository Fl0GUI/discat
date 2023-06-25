package catapi

type Response struct {
	Id        string
	Url       string
	Width     uint
	Height    uint
	Breeds    []Breed
	Favourite string
}

type Breed struct {
	Weight        Weight
	Id            string
	Name          string
	Temperament   string
	Origin        string
	Country_codes string
	Country_code  string
	Life_span     string
	Wikipedia_url string
}

type Weight struct {
	Imperial string
	Metric   string
}

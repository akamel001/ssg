package Results

type Result struct {
	service string,
	fvalue float64,
	ivalue int64
}

type Results struct {
	service string,
	results []Result
}

package domain

type Report struct {
	TotalItems       int
	TotalsByCategory map[string]int
	Top100ByPrice    []Item
}

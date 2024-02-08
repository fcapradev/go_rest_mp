package domain

type Item struct {
	ID       string
	Name     string
	Price    float64
	Category string
}
type PriceChangeRequest struct {
	TransactionID string
	ItemID        string
	Price         float64
	Status        string
}
type ByPriceDesc []Item

package view

const (
	help                 = "help"
	acceptOrder          = "accept"
	returnOrderToCourier = "return_courier"
	issueOrders          = "issue"
	acceptReturn         = "accept_return"
	listReturns          = "list_returns"
	listOrders           = "list_orders"
	setMaxGoroutines     = "set_mg"
	exit                 = "exit"
)

type command struct {
	name        string
	description string
}

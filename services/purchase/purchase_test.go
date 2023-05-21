package purchase

import (
	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/domains/vouchers"
)

var (
	ps  products.Manager
	pls plans.Manager
	vs  vouchers.Manager
)

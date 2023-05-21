package purchase

type Purchase struct {
	ProductID string `json:"product_id"`
	PlanID    string `json:"plan_id"`
	Voucher   string `json:"voucher"`
}

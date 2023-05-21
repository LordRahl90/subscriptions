package purchase

// Purchase DTO for managing purchase request values
type Purchase struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	PlanID    string `json:"plan_id"`
	Voucher   string `json:"voucher"`
}

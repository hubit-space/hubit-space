package model

type AuthResponse struct {
	UserID                 string `json:"user_id"`
	UserCode               string `json:"user_code"`
	RoleCode               string `json:"role_code"`
	Name                   string `json:"name"`
	SalesCode              string `json:"sales_code"`
	WarehouseCode          string `json:"warehouse_code"`
	WarehouseFreeTrialCode string `json:"warehouse_free_trial_code"`
	Email                  string `json:"email"`
	Menus                  []Menu `json:"menus"`
	MobileMenus            []Menu `json:"mobile_menus"`
	AccessToken            string `json:"access_token"`
	RefreshToken           string `json:"refresh_token"`
}

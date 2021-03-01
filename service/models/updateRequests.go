package models

type UpdateAccountOptions struct {
	AccountsAdd []string `json:"accounts_add,omitempty"`
	AccountsDel []string `json:"accounts_del,omitempty"`
}

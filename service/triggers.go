package service

import "net/url"

type updateOptions = struct {
	AccountsAdd []string `json:"accounts_add,omitempty"`
	AccountsDel []string `json:"accounts_del,omitempty"`
}

func AddAccountTrigger(r registrar, w *webhookService, account string, recipient url.URL) {
	w.addToCache(account, recipient)

	r.register("confirmation", updateOptions{
		AccountsAdd: []string{account},
	})
}

func AddAccountTriggers(r registrar, accounts []string) {
	r.register("confirmation", updateOptions{
		AccountsAdd: accounts,
	})
}

func RemoveAccountTrigger(r registrar, w *webhookService, account string, recipient url.URL) {
	_ = w.removeFromCache(account, recipient)

	r.register("confirmation", updateOptions{
		AccountsDel: []string{account},
	})
}

func RemoveAccountTriggers(r registrar, accounts []string) {
	r.register("confirmation", updateOptions{
		AccountsDel: accounts,
	})
}

package service

import (
	"github.com/tvandinther/nanohooks/service/models"
	"log"
)

func AddAccountTrigger(r registrar, w *webhookService, job webhookJob) {
	w.cache.add(job)

	r.register("confirmation", models.UpdateAccountOptions{
		AccountsAdd: job.accounts,
	})
}

func RemoveAccountTrigger(r registrar, w *webhookService, job webhookJob) {
	err := w.cache.remove(job)
	if err != nil {
		log.Println(err)
	}

	r.register("confirmation", models.UpdateAccountOptions{
		AccountsDel: job.accounts,
	})
}

package ambush

import (
	"bot/webhooks"
	"fmt"
	"time"
)

func (t *Task) FailedWebhook() {
	w := webhooks.Webhook{}
	e := webhooks.Embed{}

	e.SetTitle("Checkout failed :disappointed:")
	e.SetColor(0xed1c24)
	e.SetThumbnail(t.ProductImage)
	e.SetFooter(fmt.Sprintf("TraianBot - [%v]", time.Now().Format("2006-02-01 15:04:05.999999")), "https://i.pinimg.com/564x/17/6b/90/176b90a7cf4ca43e88de32369d053dab.jpg")

	e.AddField("Store", "Ambush", false)
	e.AddField("PID", t.SKU, true)
	e.AddField("Size", t.ProductSize, true)
	e.AddField("Email", fmt.Sprintf("||%v||", t.Email), true)

	w.AddEmbed(e)

	err := w.Send(t.Webhook)

	if err != nil {
		t.Error("Error sending webhook - %v", err.Error())
	} else {
		t.Info("Webhook sent")
	}
}

func (t *Task) PaypalWebhook() {

	w := webhooks.Webhook{}
	e := webhooks.Embed{}

	e.SetTitle("Succesful checkout :disappointed:")
	e.SetColor(0x00ff00)
	e.SetThumbnail(t.ProductImage)
	e.SetFooter(fmt.Sprintf("TraianBot - [%v]", time.Now().Format("2006-02-01 15:04:05.999999")), "https://i.pinimg.com/564x/17/6b/90/176b90a7cf4ca43e88de32369d053dab.jpg")
	e.SetDescription(fmt.Sprintf("[Click here](%v)", t.PaypalURL))
	e.AddField("Store", "Ambush", false)
	e.AddField("PID", t.SKU, true)
	e.AddField("Size", t.ProductSize, true)
	e.AddField("Email", fmt.Sprintf("||%v||", t.Email), true)

	w.AddEmbed(e)

	err := w.Send(t.Webhook)

	if err != nil {
		t.Error("Error sending webhook - %v", err.Error())
	} else {
		t.Info("Webhook sent")
	}

}

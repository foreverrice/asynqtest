package tpl

const EMAIL_TPL = "schedule:email"
const DELAY_TPL = "schedule:delay:email"

type EmailPayload struct {
	Email   string
	Content string
}

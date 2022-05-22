package auth

func (i EmailsMessage) Send(e Email) bool {
	return true
}

type EmailsMessage struct {
	Uid     string  `json:"uid"`
	From    string  `json:"from"`
	Subject string  `json:"subject"`
	Anchor  string  `json:"anchor"`
	HTML    string  `json:"html"`
	Emails  []Email `json:"emails"`
}

type Email struct {
	ID     string `json:"_id"`
	Email  string `json:"email"`
	Qrcode string `json:"qrcode"`
}

package email

const (
	BodyContentTypePlain string = "text/plain"
	BodyContentTypeHTML  string = "text/html"
)

type Recipient struct {
	ToEmails  []string
	CCEmails  []string
	BCCEmails []string
}
type SendEmailParams struct {
	Body        string
	BodyType    string
	Subject     string
	SenderName  string
	SenderEmail string
	Recipients  Recipient
	Attachments []string
	Headers     map[string]string
}

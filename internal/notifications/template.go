package notifications

type Email interface {
	GetRecipient() string
	GetTemplate() (string, error)
}

type ActivationEmail struct {
	Recipient string
	Username  string
}

func NewActivationEmail(recipient string, username string) *ActivationEmail {
	return &ActivationEmail{Recipient: recipient, Username: username}
}

func (a *ActivationEmail) GetRecipient() string {
	return a.Recipient
}

func (a *ActivationEmail) GetTemplate() (string, error) {
	return "", nil
}

type PasswordWasChanged struct {
	Recipient string
	Username  string
}

func NewPasswordWasChanged(recipient string, username string) *PasswordWasChanged {
	return &PasswordWasChanged{Recipient: recipient, Username: username}
}

func (a *PasswordWasChanged) GetRecipient() string {
	return a.Recipient
}

func (a *PasswordWasChanged) GetTemplate() (string, error) {
	return "", nil
}

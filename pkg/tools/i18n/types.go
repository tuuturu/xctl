package i18n

type HumanReadableError struct {
	Content string
	Key     string
}

func (receiver HumanReadableError) Error() string {
	return receiver.Content
}

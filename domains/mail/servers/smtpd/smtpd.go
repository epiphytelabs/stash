package smtpd

func Run() error {
	s, err := NewServer("api:4000")
	if err != nil {
		return err
	}

	return s.Listen(":2525")
}

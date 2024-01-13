package client

type Spec struct {
	ClientId     string
	ClientSecret string
	ApiHost      string
	RedirectUri  string
}

func (s *Spec) Validate() bool {
	return len(s.ClientId) > 0 && len(s.ClientSecret) > 0
}

package api

type Space struct {
	Identifier         string `json:"identifier"`
	Name               string `json:"name"`
	CategoryIdentifier string `json:"categoryIdentifier"`
}

func GetSpace(identifier string) (*Space, error) {
	return Exec[*Space](R(), "GET", "/workspace/space/"+identifier)
}

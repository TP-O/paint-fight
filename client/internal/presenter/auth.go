package presenter

type Login struct {
	Token  string  `json:"token"`
	Player *Player `json:"player"`
}

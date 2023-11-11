package presenter

type Login struct {
	Player       Player `json:"player"`
	AccessToken  Token  `json:"accessToken"`
	RefreshToken Token  `json:"refreshToken"`
}

type Token struct {
	Value     string `json:"value"`
	ExpiredAt int64  `json:"expiredAt"`
}

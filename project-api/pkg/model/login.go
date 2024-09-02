package model

type LoginReq struct {
	Account  string `json:"account" form:"account"`
	Password string `json:"password" form:"password"`
}

type LoginRsp struct {
	Member           Member             `json:"member"`
	TokenList        TokenList          `json:"tokenList"`
	OrganizationList []OrganizationList `json:"organizationList"`
}
type Member struct {
	Id               int64  `json:"id"`
	Code             string `json:"code"`
	Name             string `json:"name"`
	Mobile           string `json:"mobile"`
	Status           int    `json:"status"`
	CreateTime       string `json:"create_time"`
	LastLoginTime    string `json:"last_login_time"`
	OrganizationCode string `json:"organization_code"`
}

type TokenList struct {
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	TokenType      string `json:"tokenType"`
	AccessTokenExp int64  `json:"accessTokenExp"`
}

type OrganizationList struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	OwnerCode   string `json:"owner_code"`
	CreateTime  string `json:"create_time"`
	//MemberId    int64  `json:"memberId"`
	Personal int32  `json:"personal"`
	Address  string `json:"address"`
	Province int32  `json:"province"`
	City     int32  `json:"city"`
	Area     int32  `json:"area"`
	Code     string `json:"code"`
}

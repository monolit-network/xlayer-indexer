package thirdweb

import (
	"time"
)

type UserInfoProfile struct {
	Email         string  `json:"email"`
	EmailVerified bool    `json:"emailVerified"`
	HD            string  `json:"hd"`
	ID            string  `json:"id"`
	Locale        string  `json:"locale"`
	Picture       string  `json:"picture"`
	Type          string  `json:"type"`
	FamilyName    *string `json:"familyName"`
	GivenName     *string `json:"givenName"`
	Name          *string `json:"name"`
}

type UserInfo struct {
	UserID             string    `json:"userId"`
	Address            string    `json:"address"`
	CreatedAt          time.Time `json:"createdAt"`
	PublicKey          string    `json:"publicKey"`
	SmartWalletAddress string    `json:"smartWalletAddress"`

	Profiles []UserInfoProfile `json:"profiles"`
}

type WalletMeResponse struct {
	Result UserInfo `json:"result"`
}

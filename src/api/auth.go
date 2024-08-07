package api

type UserInfo struct {
	// Email: The user's email address.
	Email string `json:"email,omitempty"`
	// Gender: The user's gender.
	Gender string `json:"gender,omitempty"`
	// GivenName: The user's first name.
	GivenName string `json:"given_name,omitempty"`
	// Id: The obfuscated ID of the user.
	Id string `json:"id,omitempty"`
	// Name: The user's full name.
	Name string `json:"name,omitempty"`
	// Picture: URL of the user's picture image.
	Picture string `json:"picture,omitempty"`
	// VerifiedEmail: Boolean flag which is true if the email address is verified.
	// Always verified because we only return the user's primary email address.
	//
	// Default: true
	VerifiedEmail *bool `json:"verified_email,omitempty"`
}

package http

type BecomeArtistRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

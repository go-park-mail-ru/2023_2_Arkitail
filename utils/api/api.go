package api

const (
	Auth         = "/auth"
	Login        = "/login"
	Signup       = "/signup"
	Logout       = "/logout"
	User         = "/user"
	Places       = "/places"
	UserById     = "/users/{user_id}"
	ReviewById   = "/reviews/{reviewId}"
	Review       = "/review"
	PlaceReviews = "/places/{placeId}/reviews"
	UserReviews  = "/users/{userId}/reviews"
	Trip         = "/trip"
	Trips        = "/trips"
	TripById     = "/trips/{tripId}"
)

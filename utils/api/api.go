package api

const (
	Auth            = "/auth"
	Login           = "/login"
	Signup          = "/signup"
	Logout          = "/logout"
	User            = "/user"
	Places          = "/places"
	PlaceById       = "/places/{placeId}"
	UserById        = "/users/{user_id}"
	ReviewById      = "/reviews/{reviewId}"
	Review          = "/review"
	PlaceReviews    = "/places/{placeId}/reviews"
	UserReviews     = "/users/{userId}/reviews"
	UserAvatar      = "/user/avatar"
	Trip            = "/trip"
	Trips           = "/trips"
	TripById        = "/trips/{tripId}"
	PlaceInTrip     = "/trips/{tripId}/place"
	PlaceInTripById = "/trips/places/{placeInTripId}"
)

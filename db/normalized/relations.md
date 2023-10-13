Relation Users:
{id} -> name, username, email, location, webSite, about, avatarUrl
Таблица с пользователями.

Relation Reviews:
{id} -> user_id, placeId, text, rating

Relation Places:
{id} -> name, description, rating, cost, imageUrl

Relation Trips:
{id} -> user_id, name, isPublic

Relation Trips_to_places:
{id} -> placeId, tripId

Relation Users:
{id, name, username, email} -> location, webSite, about, avatarUrl, creationDate, lastUpdated
Таблица с пользователями.

Relation Reviews:
{id} -> userId, placeId, text, rating, creationDate, lastUpdated
Таблица с отзывами.

Relation Places:
{id} -> name, description, cost, imageUrl, creationDate, lastUpdated
Таблица с достопримечательностями.

Relation Trips:
{id} -> userId, name, isPublic, creationDate, lastUpdated
Таблица с поездками.

Relation Trips_to_places:
{id} -> placeId, tripId
Таблица для отношения многие ко многим между поездками и достопримечательностями.

1ая нормальная форма:
Каждая таблица содержит только 1 значение для каждого атрибута из кортежа.

2ая нормальная форма:
Нет ни одной функциональной зависимости от части ключа.

3ья нормальная форма:
Нет ни одной функциональной зависимости от неключевых атрибутов.

Нормальная форма Бойса — Кодда:
Все функциональные зависимости имеют в качестве детерминанты потенциальные ключи.

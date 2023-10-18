CREATE TABLE
    "user" (
        id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        name VARCHAR(30) NOT NULL,
        username VARCHAR(30) UNIQUE NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        location VARCHAR(100) DEFAULT '' NOT NULL,
        web_site TEXT DEFAULT '' NOT NULL,
        about TEXT DEFAULT '' NOT NULL,
        avatar_url VARCHAR(255) NOT NULL,
        creation_date TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
        last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
    );

CREATE TABLE
    place (
        id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        description TEXT NOT NULL,
        cost DECIMAL(19, 4) NOT NULL,
        image_url VARCHAR(255) NOT NULL,
        creation_date TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
        last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
    );

CREATE TABLE
    trip (
        id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES "user" ON DELETE CASCADE,
        name VARCHAR(30) NOT NULL,
        is_public BOOLEAN DEFAULT FALSE NOT NULL,
        creation_date TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
        last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
    );

CREATE TABLE
    review (
        id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES "user" ON DELETE CASCADE,
        place_id INTEGER NOT NULL REFERENCES place(id) ON DELETE CASCADE,
        content TEXT NOT NULL,
        rating INTEGER CHECK (rating BETWEEN 1 AND 5) NOT NULL,
        creation_date TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
        last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
    );


CREATE TABLE
    trip_to_place (
        id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        place_id INTEGER NOT NULL REFERENCES place(id) ON DELETE CASCADE,
        trip_id INTEGER NOT NULL REFERENCES trip(id) ON DELETE CASCADE
    );

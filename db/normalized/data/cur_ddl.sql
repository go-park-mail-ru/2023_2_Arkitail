CREATE TABLE
    "user" (
        id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        NAME VARCHAR(30) NOT NULL,
        PASSWORD VARCHAR(100) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        birth_date date not null,
        about TEXT,
        avatar_url VARCHAR(255),
        creation_date TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
        last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
    );

CREATE TABLE
    place (
        id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        NAME VARCHAR(100) NOT NULL,
        description TEXT NOT NULL,
        COST VARCHAR(10) NOT NULL,
        image_url VARCHAR(255),
        open_time time,
        close_time time CHECK(open_time is null or close_time > open_time),
        adress text,
        web_site varchar(40),
        email varchar(40),
        phone_number varchar(30),
        creation_date TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
        last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
    );

CREATE TYPE publicity AS ENUM ('public', 'private');
CREATE TABLE
    trip (
        id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES "user" ON DELETE CASCADE,
        description text,
        NAME VARCHAR(30) NOT NULL,
        publicity publicity DEFAULT 'private' NOT NULL,
        creation_date TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
        last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
    );

CREATE TABLE
    review (
        id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES "user" ON DELETE CASCADE,
        place_id INTEGER NOT NULL REFERENCES place (id) ON DELETE CASCADE,
        CONTENT TEXT,
        rating INTEGER CHECK (rating BETWEEN 1 AND 5) NOT NULL,
        creation_date TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
        last_updated TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
    );

CREATE TABLE
    trip_to_place (
        id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        place_id INTEGER NOT NULL REFERENCES place (id) ON DELETE CASCADE,
        trip_id INTEGER NOT NULL REFERENCES trip (id) ON DELETE CASCADE,
        first_date date,
        last_date date CHECK(last_date is null and first_date is null or last_date is null or first_date is not null and last_date > first_date)
    );

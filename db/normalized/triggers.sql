create or replace function update_time()
returns trigger
AS
$$
    declare
    begin
        new.last_updated := now();
        return new;
    end;
$$
language plpgsql;

create trigger "before_update_user_trigger"
    BEFORE update on "user"
    for each row
    execute procedure "update_time"();

create trigger "before_update_place_trigger"
    BEFORE update on place
    for each row
    execute procedure "update_time"();

create trigger "before_update_trip_trigger"
    BEFORE update on trip
    for each row
    execute procedure "update_time"();

create trigger "before_update_review_trigger"
    BEFORE update on review
    for each row
    execute procedure "update_time"();

create trigger "before_update_place_photo"
    BEFORE update on place_photo
    for each row
    execute procedure "update_time"();

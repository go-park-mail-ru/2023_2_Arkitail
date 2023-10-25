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

create trigger "update_after_user_trigger"
    AFTER update on "user"
    for each row
    execute procedure "update_time"();

create trigger "update_after_place_trigger"
    AFTER update on place
    for each row
    execute procedure "update_time"();

create trigger "update_after_trip_trigger"
    AFTER update on trip
    for each row
    execute procedure "update_time"();

create trigger "update_after_review_trigger"
    AFTER update on review
    for each row
    execute procedure "update_time"();

create trigger "update_after_place_photo"
    AFTER update on place_photo
    for each row
    execute procedure "update_time"();

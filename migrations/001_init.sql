CREATE TABLE IF NOT EXISTS segmentation (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    phone_number TEXT UNIQUE NOT NULL,
    canary_version INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- updated_at function
create function update_timestamp()
    returns trigger
as $$
begin
    new.updated_at = current_timestamp;
    return new;
end;
$$ language plpgsql;

create trigger update_segmentation_timestamp
    before update on segmentation
for each row
execute procedure update_timestamp();
CREATE TYPE days_of_month AS ENUM (
  'monday',
  'tuesday',
  'wednesday'
  'thursday',
  'friday',
  'saturday',
  'sunday'
);

CREATE TYPE dietary_restrictions AS ENUM (
  'gluten_free',
  'vegan',
  'vegetarian'
  'spicy'
);

CREATE TABLE IF NOT EXISTS business
(
  business_id SERIAL,
  name        TEXT NOT NULL,
  phone       TEXT NOT NULL,
  website     TEXT,
  PRIMARY KEY (business_id)
);

CREATE TABLE IF NOT EXISTS business_hours
(
  business_id INT           NULL,
  day         DAYS_OF_MONTH NOT NULL,
  open_time   TIME          NOT NULL,
  close_time  TIME          NOT NULL,
  FOREIGN KEY (business_id) REFERENCES business (business_id)
);

CREATE TABLE IF NOT EXISTS business_cuisine
(
  business_id INT  NULL,
  cuisine     TEXT NOT NULL,
  FOREIGN KEY (business_id) REFERENCES business (business_id)
);

CREATE TABLE IF NOT EXISTS business_user
(
  user_id  SERIAL,
  name     TEXT NOT NULL,
  email    TEXT NOT NULL,
  password TEXT NOT NULL,
  phone    TEXT,
  PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS user_to_business
(
  user_id     INT NULL,
  business_id INT NULL,
  PRIMARY KEY (user_id, business_id),
  FOREIGN KEY (user_id) REFERENCES user (user_id),
  FOREIGN KEY (business_id) REFERENCES business (business_id)
);

CREATE TABLE IF NOT EXISTS country
(
  country_id SERIAL,
  NAME       TEXT NOT NULL,
  PRIMARY KEY (country_id)
);

CREATE TABLE IF NOT EXISTS address
(
  address_id    SERIAL UNIQUE,
  business_id   INT  NULL,
  street        TEXT NOT NULL,
  city          TEXT NOT NULL,
  postal_code   INT  NOT NULL,
  state         TEXT NOT NULL,
  country_id    INT  NOT NULL,
  other_details TEXT,
  PRIMARY KEY (address_id, business_id),
  FOREIGN KEY (country_id) REFERENCES country (country_id)
);

CREATE TABLE IF NOT EXISTS address_geo
(
  geo_id      SERIAL,
  address_id  INT NULL,
  business_id INT NULL,
  latitude    NUMERIC,
  longitude   NUMERIC,
  PRIMARY KEY (geo_id, address_id),
  FOREIGN KEY (address_id) REFERENCES address (address_id),
  FOREIGN KEY (business_id) REFERENCES business (business_id)
);

CREATE TABLE IF NOT EXISTS business_image
(
  image_id    SERIAL,
  business_id INT  NOT NULL,
  path        TEXT NOT NULL,
  PRIMARY KEY (image_id, business_id),
  FOREIGN KEY (business_id) REFERENCES business (business_id)
);

CREATE TABLE IF NOT EXISTS listing
(
  listing_id          SERIAL UNIQUE,
  business_id         INT       NOT NULL,
  title               TEXT      NOT NULL,
  description         TEXT,
  old_price           DECIMAL   NOT NULL,
  new_price           DECIMAL   NOT NULL,
  listing_date        DATE      NOT NULL,
  start_time          TIME,
  end_time            TIME,
  recurring           BOOLEAN   NOT NULL,
  listing_create_date TIMESTAMP NOT NULL,
  PRIMARY KEY (listing_id),
  FOREIGN KEY (business_id) REFERENCES business (business_id)
);

CREATE TABLE IF NOT EXISTS recurring_listing
(
  listing_id INT           NULL,
  day        days_of_month NOT NULL,
  start_time TIME,
  end_time   TIME,
  PRIMARY KEY (listing_id, day),
  FOREIGN KEY (listing_id) REFERENCES listing (listing_id)
);

CREATE TABLE IF NOT EXISTS listing_image
(
  image_id   SERIAL,
  listing_id INT  NOT NULL,
  path       TEXT NOT NULL,
  PRIMARY KEY (image_id, listing_id),
  FOREIGN KEY (listing_id) REFERENCES listing (listing_id)
);

CREATE TABLE IF NOT EXISTS listing_dietary_restrictions
(
  listing_id  INT                  NULL,
  restriction DIETARY_RESTRICTIONS NOT NULL,
  PRIMARY KEY (listing_id, restriction),
  FOREIGN KEY (listing_id) REFERENCES listing (listing_id)
);

CREATE TABLE IF NOT EXISTS favourites
(
  favourite_id SERIAL,
  phone_id     TEXT NOT NULL,
  listing_id   INT  NULL,
  business_id  INT  NOT NULL,
  PRIMARY KEY (favourite_id, phone_id),
  FOREIGN KEY (business_id) REFERENCES business (business_id),
  FOREIGN KEY (listing_id) REFERENCES listing (listing_id)
);

CREATE TABLE IF NOT EXISTS happy_hour
(
  happy_hour_id       SERIAL,
  business_id         INT       NOT NULL,
  title               TEXT      NOT NULL,
  description         TEXT,
  percentage_discount DECIMAL   NOT NULL,
  listing_date        DATE      NOT NULL,
  start_time          TIME,
  end_time            TIME,
  recurring           BOOLEAN   NOT NULL,
  listing_create_date TIMESTAMP NOT NULL,
  PRIMARY KEY (happy_hour_id),
  FOREIGN KEY (business_id) REFERENCES business (business_id)
);

CREATE TABLE IF NOT EXISTS happy_hour_listing
(
  happy_hour_id INT           NULL,
  day           days_of_month NOT NULL,
  PRIMARY KEY (happy_hour_id, day),
  FOREIGN KEY (happy_hour_id) REFERENCES happy_hour (happy_hour_id)
);

CREATE TABLE IF NOT EXISTS happy_hour_image
(
  image_id      SERIAL,
  happy_hour_id INT  NOT NULL,
  path          TEXT NOT NULL,
  PRIMARY KEY (image_id, happy_hour_id),
  FOREIGN KEY (happy_hour_id) REFERENCES happy_hour (happy_hour_id)
);


INSERT INTO country (name) VALUES ('USA');
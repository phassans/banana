ALTER DATABASE banana SET timezone TO 'US/Pacific';
CREATE TYPE days_of_month AS ENUM (
  'monday',
  'tuesday',
  'wednesday',
  'thursday',
  'friday',
  'saturday',
  'sunday'
);

CREATE TYPE dietary_restrictions AS ENUM (
  'gluten free',
  'vegan',
  'vegetarian',
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
  business_id INT           NOT NULL,
  day         DAYS_OF_MONTH NOT NULL,
  open_time   TIME          NOT NULL,
  close_time  TIME          NOT NULL,
  FOREIGN KEY (business_id) REFERENCES business (business_id)
);

CREATE TABLE IF NOT EXISTS business_cuisine
(
  business_id INT  NOT NULL,
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
  user_id     INT NOT NULL,
  business_id INT NOT NULL,
  PRIMARY KEY (user_id, business_id),
  FOREIGN KEY (user_id) REFERENCES business_user (user_id),
  FOREIGN KEY (business_id) REFERENCES business (business_id)
);

CREATE TABLE IF NOT EXISTS business_country
(
  country_id SERIAL,
  NAME       TEXT NOT NULL,
  PRIMARY KEY (country_id)
);

CREATE TABLE IF NOT EXISTS business_address
(
  address_id  SERIAL UNIQUE,
  business_id INT     NOT NULL,
  street      TEXT    NOT NULL,
  city        TEXT    NOT NULL,
  postal_code INT     NOT NULL,
  state       TEXT    NOT NULL,
  country_id  INT     NOT NULL,
  latitude    NUMERIC NOT NULL,
  longitude   NUMERIC NOT NULL,
  PRIMARY KEY (address_id, business_id),
  FOREIGN KEY (country_id) REFERENCES business_country (country_id)
);

CREATE TABLE IF NOT EXISTS listing
(
  listing_id           SERIAL UNIQUE,
  business_id          INT       NOT NULL,
  title                TEXT      NOT NULL,
  old_price            DECIMAL,
  new_price            DECIMAL,
  discount             DECIMAL,
  discount_description TEXT,
  description          TEXT,
  start_date           DATE      NOT NULL,
  start_time           TIME      NOT NULL,
  end_time             TIME      NOT NULL,
  multiple_days        BOOLEAN   NOT NULL,
  end_date             DATE,
  recurring            BOOLEAN   NOT NULL,
  recurring_end_date   DATE,
  listing_type         TEXT      NOT NULL,
  listing_create_date  TIMESTAMP NOT NULL,
  PRIMARY KEY (listing_id),
  FOREIGN KEY (business_id) REFERENCES business (business_id)
);

CREATE TABLE IF NOT EXISTS listing_date
(
  listing_date_id SERIAL UNIQUE,
  listing_id      INT  NOT NULL,
  listing_date    DATE NOT NULL,
  start_time      TIME,
  end_time        TIME
);

CREATE TABLE IF NOT EXISTS listing_recurring
(
  listing_id INT           NOT NULL,
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
  listing_id  INT                  NOT NULL,
  restriction DIETARY_RESTRICTIONS NOT NULL,
  PRIMARY KEY (listing_id, restriction),
  FOREIGN KEY (listing_id) REFERENCES listing (listing_id)
);

CREATE TABLE IF NOT EXISTS favorites
(
  favorite_id SERIAL,
  phone_id    TEXT NOT NULL,
  listing_id  INT  NOT NULL,
  listing_date_id INT,
  favorite_add_date TIMESTAMP NOT NULL,
  PRIMARY KEY (favorite_id, phone_id),
  FOREIGN KEY (listing_id) REFERENCES listing (listing_id)
);

CREATE TABLE IF NOT EXISTS notifications
(
  notification_id SERIAL UNIQUE,
  notification_name TEXT,
  phone_id        TEXT NOT NULL,
  latitude        NUMERIC,
  longitude       NUMERIC,
  location        TEXT,
  price_filter    TEXT,
  distance_filter TEXT,
  keywords        TEXT,
  PRIMARY KEY (notification_id, phone_id)
);

CREATE TABLE IF NOT EXISTS notifications_dietary_restrictions
(
  notification_id INT                  NOT NULL,
  restriction     DIETARY_RESTRICTIONS NOT NULL,
  PRIMARY KEY (notification_id, restriction),
  FOREIGN KEY (notification_id) REFERENCES notifications (notification_id)
);

CREATE TABLE IF NOT EXISTS register_phone
(
  registration_id SERIAL UNIQUE,
  registration_token TEXT UNIQUE NOT NULL,
  phone_id           TEXT UNIQUE NOT NULL,
  phone_model        TEXT,
  register_date      TIMESTAMP NULL,
  update_date        TIMESTAMP NULL,
  PRIMARY KEY (registration_id)
);

CREATE TABLE IF NOT EXISTS preferences
(
  prefernce_id   SERIAL UNIQUE,
  phone_id       TEXT UNIQUE NOT NULL,
  cuisine        TEXT,
  PRIMARY KEY (phone_id,cuisine)
);

CREATE TABLE IF NOT EXISTS contact_us
(
  contact_us_id   SERIAL UNIQUE,
  phone_id        TEXT  NOT NULL,
  uname           TEXT,
  email           TEXT,
  subject         TEXT,
  comments        TEXT,
  contact_date    TIMESTAMP NOT NULL,
  PRIMARY KEY (contact_us_id)
);

CREATE TABLE IF NOT EXISTS category_to_keyword
(
  category        TEXT  NOT NULL,
  keyword         TEXT,
  PRIMARY KEY (category,keyword)
);

CREATE TABLE IF NOT EXISTS address_to_geo
(
  address TEXT,
  latitude        NUMERIC,
  longitude       NUMERIC,
  PRIMARY KEY (address)
);

CREATE TABLE IF NOT EXISTS search
(
  search_id       SERIAL UNIQUE,
  search_request  TEXT  NOT NULL,
  search_date     TIMESTAMP NOT NULL,
  PRIMARY KEY (search_id)
);

CREATE TABLE IF NOT EXISTS upvotes
(
  upvote_id SERIAL,
  phone_id    TEXT NOT NULL,
  listing_id  INT  NOT NULL,
  upvote_date TIMESTAMP NOT NULL,
  PRIMARY KEY (upvote_id, phone_id),
  FOREIGN KEY (listing_id) REFERENCES listing (listing_id)
);

CREATE TABLE IF NOT EXISTS happyhour
(
  hh_id   SERIAL UNIQUE,
  phone_id        TEXT  NOT NULL,
  name            TEXT,
  email           TEXT,
  business_owner  TEXT,
  restaurant      TEXT,
  city            TEXT,
  description     TEXT,
  submission_date TIMESTAMP NOT NULL,
  PRIMARY KEY (hh_id)
);

CREATE TABLE IF NOT EXISTS happyhour_images
(
  hh_image_id     SERIAL UNIQUE,
  hh_id           INT  NOT NULL,
  image_name      TEXT,
  submission_date TIMESTAMP NOT NULL,
  PRIMARY KEY (hh_image_id),
  FOREIGN KEY (hh_id) REFERENCES happyhour (hh_id)
);

CREATE TABLE IF NOT EXISTS report_inaccurate
(
  report_inaccurate_id SERIAL,
  phone_id    TEXT NOT NULL,
  listing_id  INT  NOT NULL,
  report_inaccurate_add_date TIMESTAMP NOT NULL,
  PRIMARY KEY (report_inaccurate_id),
  FOREIGN KEY (listing_id) REFERENCES listing (listing_id)
);

INSERT INTO business_country (name)
VALUES ('USA');
INSERT INTO business_country (name)
VALUES ('INDIA');
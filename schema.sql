CREATE TABLE IF NOT EXISTS business
  (
     business_id SERIAL,
     name        TEXT NOT NULL,
     phone       TEXT NOT NULL,
     website     TEXT,
     PRIMARY KEY (business_id)
  );

CREATE TABLE IF NOT EXISTS owner
  ( 
     owner_id    SERIAL, 
     business_id INT NULL, 
     first_name  TEXT NOT NULL, 
     last_name   TEXT NOT NULL, 
     phone       TEXT NOT NULL,
     email       TEXT NOT NULL,
     PRIMARY KEY (owner_id, business_id), 
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
     business_id   INT NULL,
     line1         TEXT NOT NULL,
     line2         TEXT, 
     city          TEXT NOT NULL, 
     postal_code   INT NOT NULL, 
     state         TEXT NOT NULL, 
     country_id    INT NOT NULL, 
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
     path        TEXT NOT NULL, 
     business_id INT NOT NULL, 
     PRIMARY KEY(image_id, business_id), 
     FOREIGN KEY (business_id) REFERENCES business (business_id)
  ); 

CREATE TABLE IF NOT EXISTS listing
  ( 
     listing_id  SERIAL UNIQUE,
     title        TEXT NOT NULL,
     description TEXT, 
     price       DECIMAL NOT NULL,
     start_time  TIMESTAMP, 
     end_time    TIMESTAMP,
     business_id INT NOT NULL,
     PRIMARY KEY (listing_id),
     FOREIGN KEY (business_id) REFERENCES business (business_id)
  );

CREATE TABLE IF NOT EXISTS listing_image
  ( 
     image_id   SERIAL, 
     path       TEXT NOT NULL, 
     listing_id INT NOT NULL, 
     PRIMARY KEY(image_id, listing_id), 
     FOREIGN KEY (listing_id) REFERENCES listing (listing_id)
  );

INSERT INTO country (name) VALUES ('USA');
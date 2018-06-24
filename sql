CREATE TABLE owner 
  ( 
     owner_id    SERIAL, 
     business_id INT NULL, 
     first_name  TEXT NOT NULL, 
     last_name   TEXT NOT NULL, 
     phone       TEXT NOT NULL, 
     PRIMARY KEY (owner_id, business_id), 
     FOREIGN KEY (business_id) REFERENCES business (business_id)
  ); 

CREATE TABLE business 
  ( 
     business_id SERIAL, 
     NAME        TEXT NOT NULL, 
     PRIMARY KEY (business_id)
  ); 

CREATE TABLE address 
  ( 
     address_id    SERIAL, 
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

CREATE TABLE geo 
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

CREATE TABLE country 
  ( 
     country_id SERIAL, 
     NAME       TEXT NOT NULL, 
     PRIMARY KEY (country_id)
  ); 

CREATE TABLE business_image 
  ( 
     image_id    SERIAL, 
     path        TEXT NOT NULL, 
     business_id INT NOT NULL, 
     PRIMARY KEY(image_id, business_id), 
     FOREIGN KEY (business_id) REFERENCES business (business_id)
  ); 

CREATE TABLE listing 
  ( 
     listing_id  SERIAL,
     name        TEXT NOT NULL,
     description TEXT, 
     price       NUMERIC NOT NULL, 
     start_time  TIMESTAMP, 
     end_time    TIMESTAMP, 
     PRIMARY KEY (listing_id)
  );

CREATE TABLE listing_to_business
  (
     listing_id    INT NOT NULL,
     business_id    INT NOT NULL,
     PRIMARY KEY(listing_id, business_id),
     FOREIGN KEY (business_id) REFERENCES business (business_id),
      FOREIGN KEY (listing_id) REFERENCES listing (listing_id)
  );

CREATE TABLE listing_image 
  ( 
     image_id   SERIAL, 
     path       TEXT NOT NULL, 
     listing_id INT NOT NULL, 
     PRIMARY KEY(image_id, listing_id), 
     FOREIGN KEY (listing_id) REFERENCES listing (listing_id)
  ); 
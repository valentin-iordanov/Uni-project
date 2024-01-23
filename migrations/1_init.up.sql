-- Table for Location
CREATE TABLE IF NOT EXISTS `location` (
    id INT PRIMARY KEY AUTO_INCREMENT,
    street VARCHAR(255) NOT NULL,
    number VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL
);

-- Table for Holiday
CREATE TABLE IF NOT EXISTS `holiday` (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    duration INT NOT NULL,
    startDate DATE NOT NULL,
    price FLOAT NOT NULL,
    freeSlots INT NOT NULL,
    locationID INT NOT NUll,
    FOREIGN KEY (locationID) REFERENCES `location`(id)
);

-- Table for Reservation
CREATE TABLE IF NOT EXISTS `reservation` (
    id INT PRIMARY KEY AUTO_INCREMENT,
    contactName VARCHAR(255) NOT NULL,
    phoneNumber VARCHAR(20) NOT NULL,
    holidayID INT NOT NULL,
    FOREIGN KEY (holidayID) REFERENCES `holiday`(id)
);

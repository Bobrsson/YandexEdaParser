-- +migrate Up
CREATE SCHEMA IF NOT EXISTS YandexEdaParser;
CREATE TABLE YandexEdaParser.restaurants (
                                             id                   SERIAL                              NOT NULL ,
                                             yandexId             INTEGER NOT NULL,
                                             created_at           TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
                                             updated_at           TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
                                             name                TEXT NOT NULL ,
                                             slag                TEXT                                   NOT NULL UNIQUE,
                                             rating              FLOAT,
                                             minimalOrderPrice FLOAT,
                                             InternalRating FLOAT,
                                             PRIMARY KEY (id)
);
CREATE TABLE YandexEdaParser.menuitem
(
    id         SERIAL                                 NOT NULL,
    yandexId   INTEGER                                NOT NULL,
    restaurantId INTEGER     NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    name       TEXT                                   NOT NULL,
    description       TEXT ,
    price     FLOAT,
    value     INTEGER,
    internalRating FLOAT,
    PRIMARY KEY (id)
);


-- +migrate Down
DROP TABLE IF EXISTS YandexEdaParser.restaurants;
DROP TABLE IF EXISTS YandexEdaParser.menuitem;
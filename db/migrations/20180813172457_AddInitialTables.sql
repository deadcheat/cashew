
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- tickets
CREATE TABLE tickets
(
    id varchar(256) PRIMARY KEY NOT NULL COMMENT 'ticket id',
    client_hostname varchar(255) NOT NULL COMMENT 'hostname of client',
    created_at datetime DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT 'created datetime'
)
DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin
COMMENT 'all tickets are not classified by their types';
CREATE UNIQUE INDEX tickets_id_uindex ON tickets (id);

-- tickets_service
CREATE TABLE ticket_service
(
    ticket_id varchar(256) NOT NULL COMMENT 'ticket id',
    service text(2083) NOT NULL COMMENT 'service url that will be or has been granted',
    created_at datetime DEFAULT current_timestamp NOT NULL COMMENT 'created datetime',
    CONSTRAINT fk_ticket_service_tickets FOREIGN KEY (ticket_id) REFERENCES tickets (id)
)
DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin
COMMENT 'service url for ticket that issued as service-ticket';
-- service is too long to create index

-- ticket_type
CREATE TABLE ticket_type
(
    ticket_id varchar(256) NOT NULL COMMENT 'ticket id',
    type varchar(50) NOT NULL COMMENT 'ticket type login/service/proxy/proxy granting/ticket granting',
    created_at datetime DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT 'created datetime',
    CONSTRAINT fk_ticket_type_tickets FOREIGN KEY (ticket_id) REFERENCES tickets (id)
)
DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin
COMMENT 'type for ticket';
CREATE INDEX idx_ticket_type_type ON ticket_type (type);

-- ticket_username
CREATE TABLE ticket_username
(
    ticket_id varchar(256) NOT NULL COMMENT 'ticket id',
    username varchar(255) NOT NULL COMMENT '',
    created_at datetime DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT 'created datetime',
    CONSTRAINT fk_ticket_username_tickets FOREIGN KEY (ticket_id) REFERENCES tickets (id)
)
DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin
COMMENT 'username for ticket such as proxy grainting ticket and service ticket';
CREATE INDEX idx_ticket_username_username ON ticket_username (username);

-- ticket_last_referenced
CREATE TABLE ticket_last_referenced
(
    ticket_id varchar(256) NOT NULL COMMENT 'ticket id',
    last_referenced_at datetime NOT NULL COMMENT '',
    created_at datetime DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT 'created datetime',
    updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL COMMENT 'updated datetime',
    CONSTRAINT fk_ticket_last_referenced_tickets FOREIGN KEY (ticket_id) REFERENCES tickets (id)
)
DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin
COMMENT 'expire datetime for ticket';
CREATE INDEX idx_ticket_last_referenced_last_referenced_at ON ticket_last_referenced (last_referenced_at);

-- ticket_iou
CREATE TABLE ticket_iou
(
    ticket_id varchar(256) NOT NULL COMMENT 'ticket id',
    iou varchar(256) NOT NULL COMMENT '',
    created_at datetime DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT 'created datetime',
    CONSTRAINT fk_ticket_iou_tickets FOREIGN KEY (ticket_id) REFERENCES tickets (id)
)
DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin
COMMENT 'iou for ticket. this will be used when ticket is proxy granting ticket';
CREATE INDEX idx_ticket_iou_iou ON ticket_iou (iou);

-- ticket_extra_attribute
CREATE TABLE ticket_extra_attribute
(
    ticket_id varchar(256) NOT NULL COMMENT 'ticket id',
    extra_attribute json NOT NULL COMMENT 'extra attributes serialized to json',
    created_at datetime DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT 'created datetime',
    CONSTRAINT fk_ticket_extra_attribute_tickets FOREIGN KEY (ticket_id) REFERENCES tickets (id)
)
DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin
COMMENT '';

-- ticket_grant_ticket
CREATE TABLE ticket_grant_ticket
(
    source_ticket_id varchar(256) PRIMARY KEY NOT NULL COMMENT 'ticket id has granted',
    destination_ticket_id varchar(256) NOT NULL COMMENT 'ticket id has been granted',
    created_at datetime DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT 'created datetime',
    CONSTRAINT fk_ticket_grant_ticket_tickets_as_source FOREIGN KEY (source_ticket_id) REFERENCES tickets (id),
    CONSTRAINT fk_ticket_grant_ticket_tickets_as_destination FOREIGN KEY (destination_ticket_id) REFERENCES tickets (id)
)
DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin
COMMENT 'ticket id that grant other ticket';
CREATE UNIQUE INDEX ticket_grant_ticket_id_uindex ON ticket_grant_ticket (id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE ticket_grant_ticket;
DROP TABLE ticket_extra_attribute;
DROP TABLE ticket_iou;
DROP TABLE ticket_expires;
DROP TABLE ticket_username;
DROP TABLE ticket_type;
DROP TABLE ticket_service;
DROP TABLE tickets;

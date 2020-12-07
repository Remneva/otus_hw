-- +goose Up
CREATE TABLE events (
    id serial primary key,
    owner bigint,
    title text,
    descr text,
    start_date date not null,
    start_time time,
    end_date date not null,
    end_time time


);
--create index owner_idx on events (owner);
--create index start_idx on events using btree (start_date, start_time);

INSERT INTO events (owner, title, descr, start_date, start_time, end_date, end_time)
VALUES
(0001, 'Atlcahualo', 'Ceasing of Water, Rising Trees', '2020-03-01', now(), '2020-03-20', now()),
(0002, 'Tlacaxipehualiztli', 'Rites of Fertility; Xipe-Totec', '2020-03-21', now(), '2020-04-09', now()),
(0002, 'Tozoztontli', 'Lesser Perforation', '2020-04-10', now(), '2020-04-29', now()),
(0002, 'Huey Tozoztli', 'Greater Perforation', '2020-04-30', now(), '2020-05-19', now()),
(0003, 'Tōxcatl', 'Dryness', '2020-05-20', now(), '2020-05-08', now());


-- +goose Down
drop table events;
create USER calendar with PASSWORD 'calendar';
create SCHEMA calendar;
grant all privileges on schema calendar TO calendar;
alter user calendar set SEARCH_PATH to calendar;
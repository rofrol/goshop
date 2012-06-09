BEGIN TRANSACTION;
CREATE TABLE products (
  id integer primary key autoincrement,
  title string not null,
  text string not null,
  price float not null
);
CREATE TABLE users (
  id integer primary key autoincrement,
  name1 string not null,
  name2 string,
  surname string not null,
);
CREATE TABLE orders (
  id integer primary key autoincrement,
  client_id integer not null,
  product_id integer not null,
  price float not null,
  sell_date date not null
);
COMMIT;

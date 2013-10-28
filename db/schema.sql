PRAGMA foreign_keys = ON;
BEGIN TRANSACTION;
drop table if exists products;
CREATE TABLE products (
  id integer primary key autoincrement,
  title string not null,
  text string not null,
  price float not null
);
drop table if exists users;
CREATE TABLE users (
  id integer primary key autoincrement,
  login string not null,
  password string not null,
  name1 string not null,
  name2 string,
  surname string not null
);
drop table if exists orders;
CREATE TABLE orders (
  id integer primary key autoincrement,
  user_id integer not null,
  product_id integer not null,
  price float not null,
  sell_date date not null,
  FOREIGN KEY(user_id) REFERENCES users(id),
  FOREIGN KEY(product_id) REFERENCES products(id)
);

create unique index users_u on users(login);

insert into users (login,password,name1,name2,surname) values ('admin','adminp','admin1','admin2','admin3');
COMMIT;

CREATE DATABASE casdb CHARACTER SET utf8mb4 collate utf8mb4_bin;
CREATE USER casuser IDENTIFIED WITH mysql_native_password BY 'password';
GRANT ALL PRIVILEGES on casdb.* to 'casuser'@'%';

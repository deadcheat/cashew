CREATE database casdb;
CREATE USER casuser IDENTIFIED WITH mysql_native_password BY 'password';
GRANT ALL PRIVILEGES on casdb.* to 'casuser'@'%';

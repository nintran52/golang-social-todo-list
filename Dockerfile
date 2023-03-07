FROM alpine

WORKDIR /app/
ADD ./app /app/

ENTRYPOINT ["./app"]

#
#JWT_SECRET=200Lab.io;
#MYSQL_GORM_DB_TYPE=mysql;
#MYSQL_GORM_DB_URI=root:my-secret-pw@tcp(g09-mysql:3306)/social-todo-list?charset=utf8mb4&parseTime=True&loc=Local
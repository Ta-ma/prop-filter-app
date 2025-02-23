docker run --name filter-prop-db \
	-e POSTGRES_USER=dbadmin -e POSTGRES_PASSWORD=filterpr0p \
	-e POSTGRES_DB=filter-prop \
	-p 8080:5432 -it postgres
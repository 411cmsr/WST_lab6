WST_lab4_server

/api/v1

http Method              Route
GET                /api/v1/persons                  Search persons 
GET                /api/v1/persons/list             Fetch all persons
POST               /api/v1/persons                  Add person
GET                /api/v1/person/:personId         Retrieve person     
PUT                /api/v1/person/:personId         Update person
DELETE             /api/v1/person/:personId         Delete person    




Invoke-WebRequest -Uri "http://localhost:8095/api/v1/persons?query=34" -Method GET
Invoke-WebRequest -Uri "http://localhost:8095/api/v1/persons?query=Ольга" -Method GET

Invoke-WebRequest -Uri "http://localhost:8095/api/v1/persons/list" -Method GET

Invoke-RestMethod -Uri "http://localhost:8095/api/v1/persons" -Method POST -ContentType "application/json; charset=utf-8" -Body '{"name":"OLga","surname":"Ditr","age":34,"email":"olga@mail.com","telephone":"+70011234576"}'

Invoke-WebRequest -Uri "http://localhost:8095/api/v1/person/1817" -Method GET

Invoke-RestMethod -Uri "http://localhost:8095/api/v1/person/1629" -Method PUT -ContentType "application/json; charset=utf-8" -Body '{"name":"Olga","surname":"Berig","age":34,"email":"olga@mail.com","telephone":"+70011234576"}'

Invoke-WebRequest -Uri "http://localhost:8095/api/v1/person/1817" -Method DELETE

# factory

## run
``` 
docker compose -d up
```

## create database structure
```
$ migrate create -ext sql -dir database/migration -seq create_tables
$ make migration_up
$ migrate create -ext sql -dir database/migration -seq create_insert_question_answer_tags_stored_procedure
$ make migration_up
```

## ask a question to the chat
```
curl "http://localhost:8080/get?question=Quelle%20est%20la%20capital%20de%20la%20France%3F"
```
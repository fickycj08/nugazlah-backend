## Endpoint

- login

- create class
- join class
- get my class

- create task (only class creator)
- get my task
- mark task done
- get detail task

new : ```migrate create -ext sql -dir migration -seq {migration-name}```

up : ```migrate -database postgres://postgres:root@localhost:5432/nugazlah?sslmode=disable -path ./migration up```

down : ```migrate -database postgres://postgres:root@localhost:5432/nugazlah?sslmode=disable -path ./migration down```
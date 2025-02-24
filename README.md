# connect-task-manger

## 開発時に使用するコマンド群

### Dockerコンテナ起動

```zsh
# GoコンテナとMySQLコンテナを構築
make up
```

### sqlcによるコード生成

```zsh
sqlc generate
```

### bufによるコード生成

```zsh
buf generate
```

### grpcurlコマンドツールをインストール

```zsh
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### Goバックエンドサービスの起動

```zsh
go run cmd/server/main.go
```

## user関連のエンドポイント一覧

```zsh
grpcurl -plaintext -d '{"name": "Test User", "email": "test@example.com", "password": "password123"}' localhost:8080 user.v1.UserService/CreateUser

grpcurl -plaintext -d '{"email": "test@example.com", "password": "password123"}' localhost:8080 user.v1.UserService/Login

grpcurl -plaintext -H "Authorization: Bearer <取得したaccess_token>" localhost:8080 user.v1.UserService/GetMe

grpcurl -plaintext -H "Authorization: Bearer <取得したaccess_token>" -d '{"name": "New Name", "email": "new_email@example.com", "password": "new_password"}' localhost:8080 user.v1.UserService/UpdateUser
```

## task関連のエンドポイント一覧

```zsh
grpcurl -plaintext -H "Authorization: Bearer <取得したaccess_token>" -d '{"title": "My First Task", "description": "This is a description of my first task."}' localhost:8080 task.v1.TaskService/CreateTask

grpcurl -plaintext -H "Authorization: Bearer <取得したaccess_token>" localhost:8080 task.v1.TaskService/ListTasks

grpcurl -plaintext -H "Authorization: Bearer <取得したaccess_token>" -d '{"id": "<タスクのID>", "title": "Updated Task Title", "description": "Updated task description.", "isCompleted": true}' localhost:8080 task.v1.TaskService/UpdateTask

grpcurl -plaintext -H "Authorization: Bearer <取得したaccess_token>" -d '{ "id": "<タスクのID>"}' localhost:8080 task.v1.TaskService/DeleteTask
```

## grpcurl 実行例

### user.v1.UserService/CreateUser

```zsh
/opt/task_manage # grpcurl -plaintext -d '{"name": "Test User", "email": "test@example.com", "password": "password123"}' localhost:8080 user.v1.UserService/CreateUser
{}
```

### user.v1.UserService/Login

```zsh
/opt/task_manage # grpcurl -plaintext -d '{"email": "test@example.com", "password": "password123"}' localhost:8080 user.v1.UserService/Login
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiM2UxZWFmNzMtZDlkNi00NGY5LWEyZDEtMGJiZDkwNzJlYTIzIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDE1OTI4OCwibmJmIjoxNzQwMTU4Mzg4LCJpYXQiOjE3NDAxNTgzODh9.NBTfXvWEWlnm_csCVsGqEh72Ql0ISpBhIZ4LqD7kJXc"
}
```

### user.v1.UserService/GetMe

```zsh
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZWRmZDY1MzUtZjRlOS00OTBlLWEyNTAtNmEyMGMzNDc0M2NmIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDQwNDY0NywibmJmIjoxNzQwNDAzNzQ3LCJpYXQiOjE3NDA0MDM3NDd9.GiuDEoPE1vlB_xEQnwjWE71fmQeSb-3UFIKmnqoR-0c" localhost:8080 user.v1.UserService/GetMe
{
  "user": {
    "id": "3e1eaf73-d9d6-44f9-a2d1-0bbd9072ea23",
    "name": "Test User",
    "email": "test@example.com",
    "createdAt": "2025-02-22T02:19:24Z",
    "updatedAt": "2025-02-22T02:19:24Z"
  }
}
```

### user.v1.UserService/UpdateUser

```zsh
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiM2UxZWFmNzMtZDlkNi00NGY5LWEyZDEtMGJiZDkwNzJlYTIzIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDE1OTI4OCwibmJmIjoxNzQwMTU4Mzg4LCJpYXQiOjE3NDAxNTgzODh9.NBTfXvWEWlnm_csCVsGqEh72Ql0ISpBhIZ4LqD7kJXc" -d '{"name": "New Name", "email": "new_email@example.com", "password": "new_password"}' localhost:8080 user.v1.UserService/UpdateUser
{
  "user": {
    "id": "3e1eaf73-d9d6-44f9-a2d1-0bbd9072ea23",
    "name": "New Name",
    "email": "new_email@example.com",
    "createdAt": "2025-02-22T02:19:24Z",
    "updatedAt": "2025-02-22T02:20:50Z"
  }
}
```

### task.v1.TaskService/CreateTask

```zsh
# パターン①
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZWRmZDY1MzUtZjRlOS00OTBlLWEyNTAtNmEyMGMzNDc0M2NmIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDQwNzc0MCwibmJmIjoxNzQwNDA2ODQwLCJpYXQiOjE3NDA0MDY4NDB9.25wrI8edTyn4ucPeLTrz6xRbeuXxyCcUnf7T3SJYrJo" -d '{"title": "Buy groceries","description": "Milk, eggs, bread, cheese","priority": "high","due_date": null}' localhost:8080 task.v1.TaskService/CreateTask
{}
```

```zsh
# パターン②
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZWRmZDY1MzUtZjRlOS00OTBlLWEyNTAtNmEyMGMzNDc0M2NmIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDQwNzc0MCwibmJmIjoxNzQwNDA2ODQwLCJpYXQiOjE3NDA0MDY4NDB9.25wrI8edTyn4ucPeLTrz6xRbeuXxyCcUnf7T3SJYrJo" -d '{"title": "Buy groceries","description": "Milk, eggs, bread, cheese","priority": "high","due_date": "2025-03-29T12:00:00Z"}' localhost:8080 task.v1.TaskService/CreateTask
{}
```

### task.v1.TaskService/ListTasks

```zsh
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZWRmZDY1MzUtZjRlOS00OTBlLWEyNTAtNmEyMGMzNDc0M2NmIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDQwNzc0MCwibmJmIjoxNzQwNDA2ODQwLCJpYXQiOjE3NDA0MDY4NDB9.25wrI8edTyn4ucPeLTrz6xRbeuXxyCcUnf7T3SJYrJo" localhost:8080 task.v1.TaskService/ListTasks
{
  "tasks": [
    {
      "id": "98c29d67-564e-41f2-9a89-e231381be0e8",
      "title": "Buy groceries",
      "description": "Milk, eggs, bread, cheese",
      "userId": "edfd6535-f4e9-490e-a250-6a20c34743cf",
      "createdAt": "2025-02-24T23:23:33Z",
      "updatedAt": "2025-02-24T23:23:33Z",
      "priority": "high",
      "dueDate": "2025-03-29T00:00:00Z"
    },
    {
      "id": "a466aff0-48c6-42a0-9bf1-bf7d61062cc5",
      "title": "Buy groceries",
      "description": "Milk, eggs, bread, cheese",
      "userId": "edfd6535-f4e9-490e-a250-6a20c34743cf",
      "createdAt": "2025-02-24T23:23:05Z",
      "updatedAt": "2025-02-24T23:23:05Z",
      "priority": "high"
    },
    {
      "id": "1c27e56e-9ca2-452f-9c7f-7f22258f0b2d",
      "title": "Buy groceries",
      "description": "Milk, eggs, bread, cheese",
      "userId": "edfd6535-f4e9-490e-a250-6a20c34743cf",
      "createdAt": "2025-02-24T23:21:03Z",
      "updatedAt": "2025-02-24T23:21:03Z",
      "priority": "high"
    }
  ]
}
```

### task.v1.TaskService/UpdateTask

```zsh
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZWRmZDY1MzUtZjRlOS00OTBlLWEyNTAtNmEyMGMzNDc0M2NmIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDQwODY3NiwibmJmIjoxNzQwNDA3Nzc2LCJpYXQiOjE3NDA0MDc3NzZ9.JUD9BjcIj76TcgnSLMY4yw7lVbVvuMnY0Cj0MSWuc-M" -d '{ "id": "1c27e56e-9ca2-452f-9c7f-7f22258f0b2d", "title": "Buy groceries and snacks", "description": "Milk, eggs, bread, cheese, chips, soda", "isCompleted": true, "assignee_id": "edfd6535-f4e9-490e-a250-6a20c34743cf", "priority": "medium", "due_date": "2025-04-10T00:00:00Z"}' localhost:8080 task.v1.TaskService/UpdateTask
{
  "task": {
    "id": "1c27e56e-9ca2-452f-9c7f-7f22258f0b2d",
    "title": "Buy groceries and snacks",
    "description": "Milk, eggs, bread, cheese, chips, soda",
    "isCompleted": true,
    "userId": "edfd6535-f4e9-490e-a250-6a20c34743cf",
    "createdAt": "2025-02-24T23:21:03Z",
    "updatedAt": "2025-02-24T23:36:36Z",
    "assigneeId": "edfd6535-f4e9-490e-a250-6a20c34743cf",
    "priority": "medium",
    "dueDate": "2025-04-10T00:00:00Z"
  }
}
```

### task.v1.TaskService/DeleteTask

```zsh
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjhiODQ0NzEtY2U5Mi00ZTNmLWI3OGMtOGQyYTcxMzEwYzZlIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDIzMDE4MCwibmJmIjoxNzQwMjI5MjgwLCJpYXQiOjE3NDAyMjkyODB9.PgxBXd9X1zN5REHsyMkxO3t9ZBHayJ30UzEVaYjcXL0" -d '{ "id": "d1970f2c-345c-45a3-9141-887ac761f74e"}' localhost:80
80 task.v1.TaskService/DeleteTask
{}
```

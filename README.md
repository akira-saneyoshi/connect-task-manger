# connect-task-manger

## grpcurlコマンドツールをインストール

```zsh
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

## バックエンドサービスの起動

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
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiM2UxZWFmNzMtZDlkNi00NGY5LWEyZDEtMGJiZDkwNzJlYTIzIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDE1OTI4OCwibmJmIjoxNzQwMTU4Mzg4LCJpYXQiOjE3NDAxNTgzODh9.NBTfXvWEWlnm_csCVsGqEh72Ql0ISpBhIZ4LqD7kJXc" localhost:8080 user.v1.UserService/GetMe
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
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjhiODQ0NzEtY2U5Mi00ZTNmLWI3OGMtOGQyYTcxMzEwYzZlIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDIzMDE4MCwibmJmIjoxNzQwMjI5MjgwLCJpYXQiOjE3NDAyMjkyODB9.PgxBXd9X1zN5REHsyMkxO3t9ZBHayJ30UzEVaYjcXL0" -d '{"title": "My First Task", "description": "This is a description of my first task."}' localhost:8080 task.v1.TaskService/CreateTask
{}
```

### task.v1.TaskService/ListTasks

```zsh
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjhiODQ0NzEtY2U5Mi00ZTNmLWI3OGMtOGQyYTcxMzEwYzZlIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDIzMDE4MCwibmJmIjoxNzQwMjI5MjgwLCJpYXQiOjE3NDAyMjkyODB9.PgxBXd9X1zN5REHsyMkxO3t9ZBHayJ30UzEVaYjcXL0" localhost:8080 task.v1.TaskService/ListTasks
{
  "tasks": [
    {
      "id": "d1970f2c-345c-45a3-9141-887ac761f74e",
      "title": "My First Task",
      "description": "This is a description of my first task.",
      "userId": "f8b84471-ce92-4e3f-b78c-8d2a71310c6e",
      "createdAt": "2025-02-22T22:04:02Z",
      "updatedAt": "2025-02-22T22:04:02Z"
    }
  ]
}
```

### task.v1.TaskService/UpdateTask

```zsh
opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjhiODQ0NzEtY2U5Mi00ZTNmLWI3OGMtOGQyYTcxMzEwYzZlIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDIzMDE4MCwibmJmIjoxNzQwMjI5MjgwLCJpYXQiOjE3NDAyMjkyODB9.PgxBXd9X1zN5REHsyMkxO3t9ZBHayJ30UzEVaYjcXL0" -d '{"id": "d1970f2c-345c-45a3-9141-887ac761f74e", "title": "Updated Task Title", "description": "Updated task description.", "isCompleted": true}' localhost:8080 task.v1.TaskService/UpdateTask
{
  "task": {
    "id": "d1970f2c-345c-45a3-9141-887ac761f74e",
    "title": "Updated Task Title",
    "description": "Updated task description.",
    "isCompleted": true,
    "userId": "f8b84471-ce92-4e3f-b78c-8d2a71310c6e",
    "createdAt": "2025-02-22T22:04:02Z",
    "updatedAt": "2025-02-22T22:05:45Z"
  }
}
```

### task.v1.TaskService/DeleteTask

```zsh
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjhiODQ0NzEtY2U5Mi00ZTNmLWI3OGMtOGQyYTcxMzEwYzZlIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDIzMDE4MCwibmJmIjoxNzQwMjI5MjgwLCJpYXQiOjE3NDAyMjkyODB9.PgxBXd9X1zN5REHsyMkxO3t9ZBHayJ30UzEVaYjcXL0" -d '{ "id": "d1970f2c-345c-45a3-9141-887ac761f74e"}' localhost:80
80 task.v1.TaskService/DeleteTask
{}
```
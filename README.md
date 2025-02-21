# connect-task-manger

## grpcurlコマンドツールをインストール

```
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

## バックエンドサービスの起動

```
go run cmd/server/main.go
```

## user関連のエンドポイント一覧

```
grpcurl -plaintext -d '{"name": "Test User", "email": "test@example.com", "password": "password123"}' localhost:8080 user.v1.UserService/CreateUser

grpcurl -plaintext -d '{"email": "test@example.com", "password": "password123"}' localhost:8080 user.v1.UserService/Login

grpcurl -plaintext -H "Authorization: Bearer <取得したaccess_token>" localhost:8080 user.v1.UserService/GetMe

grpcurl -plaintext -H "Authorization: Bearer <取得したaccess_token>" -d '{"name": "New Name", "email": "new_email@example.com", "password": "new_password"}' localhost:8080 user.v1.UserService/UpdateUser
```

## grpcurl 実行例

### user.v1.UserService/CreateUser

```
/opt/task_manage # grpcurl -plaintext -d '{"name": "Test User", "email": "test@example.com", "password": "password123"}' localhost:8080 
user.v1.UserService/CreateUser
{}
```

### user.v1.UserService/Login

```
/opt/task_manage # grpcurl -plaintext -d '{"email": "test@example.com", "password": "password123"}' localhost:8080 user.v1.UserService/L
ogin
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiM2UxZWFmNzMtZDlkNi00NGY5LWEyZDEtMGJiZDkwNzJlYTIzIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDE1OTI4OCwibmJmIjoxNzQwMTU4Mzg4LCJpYXQiOjE3NDAxNTgzODh9.NBTfXvWEWlnm_csCVsGqEh72Ql0ISpBhIZ4LqD7kJXc"
}
```

### user.v1.UserService/GetMe

```
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiM2UxZWFmNzMtZDlkNi0
0NGY5LWEyZDEtMGJiZDkwNzJlYTIzIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDE1OTI4OCwibmJmIjoxNzQwMTU4Mzg
4LCJpYXQiOjE3NDAxNTgzODh9.NBTfXvWEWlnm_csCVsGqEh72Ql0ISpBhIZ4LqD7kJXc" localhost:8080 user.v1.UserService/GetMe
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

```
/opt/task_manage # grpcurl -plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiM2UxZWFmNzMtZDlkNi0
0NGY5LWEyZDEtMGJiZDkwNzJlYTIzIiwiaXNzIjoiZGRkLWF1dGgtYXBwIiwic3ViIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTc0MDE1OTI4OCwibmJmIjoxNzQwMTU4Mzg
4LCJpYXQiOjE3NDAxNTgzODh9.NBTfXvWEWlnm_csCVsGqEh72Ql0ISpBhIZ4LqD7kJXc" -d '{"name": "New Name", "email": "new_email@example.com", "passw
ord": "new_password"}' localhost:8080 user.v1.UserService/UpdateUser
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

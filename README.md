# atoolforparsingandrunninghttpfiles

`atfparhf` for short

This is a command-line tool that can parse `.http` files and run the resulting collection
of any requests found inside them. 

## Why?

I don't like Postman and other tools like it, and I often forget the semantics of the `curl` command. The `.http` file
just makes sense to me.

<details>

<summary>example.http</summary>

```http request
### GET request to /health
GET http://localhost:8080/health

### POST request to /todos
POST http://localhost:8080/todos
Content-Type: application/json

{
  "task": "learn go"
}
```

</details>

# atoolforparsingandrunninghttpfiles

`atfparhf` for short

This is a command-line tool that can parse `.http` files and run the resulting collection
of any requests found inside of them. 

## Why?

I don't like Postman, and I often forget the semantics of the `curl` command. The `.http` file
just makes sense to me.

<details>

<summary>get_health_check.http</summary>

```http
### GET request to /api/health
GET http:${ROOT_URL}/api/health

###
```

</details>


### Run the Command

```bash
atoolforparsingandrunninghttpfiles --file get_health_check.http
```

### Result

```text
OK
```

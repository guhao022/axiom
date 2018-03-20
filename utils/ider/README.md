# ID生成器

## 使用

```go
# idgener
id := ID()
for i := 1; i <= 1000000; i++ {
	fmt.Printf("ID: %s \n", <-id)
}

# snowflake
id := NewID(1)
for i := 1; i <= 10; i++ {
	fmt.Printf("%d: %d \n", wid, id.Next())
}

```
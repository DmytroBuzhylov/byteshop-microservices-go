module AdminService

go 1.23.6

require (
	ByteShop/generated/auth v0.0.0-00010101000000-000000000000
	ByteShop/generated/product v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/labstack/echo/v4 v4.13.3
	google.golang.org/grpc v1.70.0
)

replace ByteShop/generated/product => ../generated/product

replace ByteShop/generated/auth => ../generated/auth

require (
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

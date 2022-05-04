# receipt-uploader

## Available Commands

In the project directory, you can run:

### `go mod tidy`

This command will download all the dependencies that are required \
in your source files and update go.mod file with that dependency.

### `go run github.com/99designs/gqlgen generate`
After make some modification in file `schema.graphqls`, run this command to generate models and resolver

### `go run server`
Start web server


## Available Endpoints
1. signup(graphql mutation) -- register new user
2. login(graphql mutation) -- login, if succeed a valid JWT will be returned
3. uploadReceipt(graphql mutation) -- create receipt with uploaded image file(a valid token needed)
4. fetchReceipts(graphql query) -- fetch all receipts created by user(a valid token needed)
5. fetchReceipt(graphql query) -- fetch a specific receipt, receipt image will be scaled if paramter `scaleRatio` is a valid number (0, 100)


### Note
I do not have golong working experience before,there might be some stupid errors in my implementation.
Test is missing, error handlering need improve as well. Need some time to learn this part of golang.
# golang-assignment

## Setup Instructions

1. **Install Dependencies**:
    ```sh
    go mod tidy
    ```

2. **MySQL Setup**:
    - Create a database named `golang_assignement`.

3. **Redis Setup**:
    - Ensure Redis is running on `localhost:6379`.
	- To check  cache use `KEYS employee*`

4. **Run the Server**:
    ```sh
    go run main.go
    ```

5. **Upload Excel File**:
    - Use Postman to upload an Excel file to `http://localhost:8080/upload`.

6. **View, Edit, and Delete Data**:
    - Use the following endpoints:
        - `GET /employees`: View data.
        - `PUT /employee/:id`: Edit a record.
        - `DELETE /employee/:id`: Delete a record.

## Project Structure

- **main.go**: Entry point of the application.
- **handlers/**: Contains API handlers.
- **config/**: Configuration for database and Redis.
- **utils/**: Utility functions for processing Excel files.

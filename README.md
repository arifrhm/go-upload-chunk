# Go Upload Chunk

A Go application for handling chunked file uploads with support for resumable uploads, file type validation, and file size checks.

## Features

- **Chunked File Uploads**: Efficiently upload large files in chunks.
- **File Type Validation**: Only allows PDF files.
- **File Size Validation**: Ensures each chunk is no larger than 1 MB.
- **Resumable Uploads**: Supports resuming file uploads if interrupted.

## Prerequisites

- Go 1.18 or higher
- `github.com/labstack/echo/v4` v4.12.0 for HTTP server functionality
- `github.com/joho/godotenv` for managing environment variables

## Setup

1. **Clone the Repository**

    ```bash
    git clone https://github.com/arifrhm/go-upload-chunk.git
    cd go-upload-chunk
    ```

2. **Install Dependencies**

    ```bash
    go mod tidy
    ```

3. **Create a `.env` File**

    Create a `.env` file in the root directory with the following content:

    ```dotenv
    APP_PORT=8001
    UPLOAD_PATH=./uploads
    ```

4. **Run the Application**

    ```bash
    go run main.go
    ```

    The application will start on the port specified in the `.env` file (default: `8001`).

## Directory Structure

- `README.md`: Project overview and setup instructions.
- `main.go`: Entry point of the application.
- `main_test.go`: Contains tests for the main application.
- `config/`: Configuration settings.
  - `config.go`: Loads environment variables and provides configuration settings.
- `handlers/`: HTTP route handlers.
  - `handlers.go`: Defines the routes and handlers for file uploads.
- `services/`: Business logic and file handling.
  - `services.go`: Functions for handling file uploads and assembling chunks.
- `utils/`: Utility functions.
  - `fileutils.go`: Functions for initializing directories and managing files.
  - `validator.go`: Functions for validating file types and sizes.
- `upload.html`: HTML form for uploading files.
- `uploads/`: Directory where uploaded files are stored.
  - `upload_here.txt`: Example file for testing.

## Endpoints

### Upload File

- **Endpoint**: `/upload`
- **Method**: `GET`
- **Description**: Serves the HTML form for file uploads.

### Upload Chunk

- **Endpoint**: `/upload-chunk`
- **Method**: `POST`
- **Description**: Uploads a file chunk. Requires the following form-data fields:
  - `file`: The file chunk
  - `chunk_index`: The index of the chunk
  - `total_chunks`: Total number of chunks
  - `file_name`: The name of the file

### Resume Upload

- **Endpoint**: `/resume-upload`
- **Method**: `GET`
- **Description**: Retrieves the status of the file upload. Requires a query parameter:
  - `file_name`: The name of the file

## Contributing

Contributions are welcome! Please open issues or submit pull requests for suggestions or improvements.
# Beatify Core

Welcome to Beatify Core! This project is a backend service for a music streaming application that handles music management, user authentication, and streaming functionalities.

## Features

- **Music Hosting**: Store and manage music files.
- **User Authentication**: Secure user registration and login.
- **Music Streaming**: Stream music with support for HTTP range requests.
- **Music Uploading**: Allow users to upload their music.
- **Music Fetching**: Retrieve music from alternative resources if not available on the platform.

## Project Structure

- **Music Management**: Handles operations related to storing and managing music.
- **User Authentication**: Manages user accounts and authentication using JWTs.
- **Streaming Service**: Supports streaming of music files with range requests.

## Installation

To run this project locally, you need to have Go installed. Follow these steps to get started:

1. Clone the repository:
    ```bash
    git clone https://github.com/YuanziX/beatify-core.git
    ```
2. Navigate to the project directory:
    ```bash
    cd beatify-core
    ```
3. Install dependencies:
    ```bash
    go mod tidy
    ```

## Configuration

Create a `.env` file in the root directory of the project and add your environment variables. Example configuration is included in .env.sample

## Running the Project

To start the server, use the following command:

```bash
make run
```

## API Endpoints

- **GET /users**: Get a list of users (To be protected with RBA).
- **POST /user**: Create a new user.
- **GET /user/{email}/verify**: Verify a user's email.
- **GET /user/{email}/isVerified**: Check if a user's email is verified.
- **GET /user/{email}/resendVerificationMail**: Resend verification email.
- **GET /user/{email}**: Get user details by email (protected).
- **DELETE /user/{email}**: Delete a user by email (protected).
- **POST /login**: Log in a user.
- **GET /logout**: Log out a user (protected).
- **GET /music**: Get a list of available music.
- **GET /music/stream?id=n**: Stream music file with id n.

## Contributing

If you would like to contribute to Beatify Core, please fork the repository and create a pull request with your changes. Make sure to follow the code style guidelines and add appropriate tests.

## License

This project is licensed under the MIT License.

## Contact

For any questions or feedback, please contact [YuanziX](mailto:achubadyal4@gmail.com).
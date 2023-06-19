Car Rentals API
This is a Go API for car rentals, built using the Gin framework. It uses MongoDB as the database, Cloudinary for storing images, JWT for authentication, and role-based authorization.

Getting Started
To get started with the Car Rentals API, follow the steps below:

Prerequisites
Go (version 1.16 or later) installed on your machine
MongoDB installed and running
Cloudinary account (for storing images)
Postman or any API testing tool


Clone the repository:

PORT=5000
MONGODB_URI=your-mongodb-uri
CLOUDINARY_API_KEY=your-cloudinary-api-key
CLOUDINARY_API_SECRET=your-cloudinary-api-secret
JWT_SECRET=your-jwt-secret



Start the server by running the following command:

go run main.go




Car Rentals API
This is a Go API for car rentals, built using the Gin framework. It uses MongoDB as the database, Cloudinary for storing images, JWT for authentication, and role-based authorization.

Getting Started
To get started with the Car Rentals API, follow the steps below:

Prerequisites
Go (version 1.16 or later) installed on your machine
MongoDB installed and running
Cloudinary account (for storing images)
Postman or any API testing tool
Installation
Clone the repository:

shell
Copy code
git clone https://github.com/your-username/car-rentals-api.git
Change into the project directory:

shell
Copy code
cd car-rentals-api
Run go mod tidy to download and install the project dependencies.

Set up the required environment variables. Create a .env file in the project root directory and add the following variables:

plaintext
Copy code
PORT=5000
MONGODB_URI=your-mongodb-uri
CLOUDINARY_API_KEY=your-cloudinary-api-key
CLOUDINARY_API_SECRET=your-cloudinary-api-secret
JWT_SECRET=your-jwt-secret
Replace the values with your MongoDB URI, Cloudinary API credentials, and a secret key for JWT.

Start the server by running the following command:

shell
Copy code
go run main.go
The server will start running on http://localhost:5000.

API Endpoints
The following API endpoints are available:

POST /api/v1/auth/signup - User registration
POST /api/v1/auth/login - User login
GET /api/v1/cars - Get all cars
GET /api/v1/cars/:id - Get a specific car by ID
POST /api/v1/cars - Create a new car (requires authentication)
PUT /api/v1/cars/:id - Update a car by ID (requires authentication)
DELETE /api/v1/cars/:id - Delete a car by ID (requires authentication)
Make sure to replace your-mongodb-uri, your-cloudinary-api-key, your-cloudinary-api-secret, and your-jwt-secret with the actual values in your .env file.

Authentication and Authorization
User registration (/api/v1/auth/signup) returns a JWT token upon successful registration.
User login (/api/v1/auth/login) returns a JWT token upon successful authentication.
To access protected routes, include the Authorization header with the JWT token prefixed by Bearer. Example: Authorization: Bearer <token>.
Role-based authorization is implemented, with roles ADMIN and USER. Only users with the ADMIN role can create, update, and delete cars.
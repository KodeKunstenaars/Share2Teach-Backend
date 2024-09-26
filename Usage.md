# Project Overview

Welcome to the **Share2Teach API** usage guide. This document provides a comprehensive overview of the API, detailing how you can interact with it to manage educational resources and authenticate users. This guide will help you understand and utilise the Share2Teach API effectively.

# Introduction

**Share2Teach** is an Open Educational Resource (OER) platform aimed at fostering a community of learners and educators. The platform encourages the sharing of educational resources to promote self-directed and collaborative learning. The Share2Teach API is the backbone of this platform, providing endpoints for:
  - **User Authentication**: Registering and logging in users, with role-based access control.
  - **Document Management**: Uploading, downloading, and managing educational documents and their metadata.
  - **User Interactions**: Rating and reporting documents, as well as searching for resources.

Built with **Go (Golang)** and utilising the **Chi router** for handling HTTP requests, the API leverages **MongoDB Atlas** for database management and **AWS S3** for secure file storage.


# Table of contents
- [Project Overview](#project-overview)
- [Introduction](#introduction)
- [Table of contents](#table-of-contents)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Endpoint Summaries](#api-endpoint-summaries)
  - [Authentication Endpoints](#authentication-endpoints)
  - [Document Management Endpoints](#document-management-endpoints)
  - [User Interaction Endpoints](#user-interaction-endpoints)
- [Usage](#usage)
- [Postman Collection](#postman-collection)
- [Database Migration](#database-migration)
- [Troubleshooting](#troubleshooting)
- [Testing](#testing)
- [License](#license)
- [Additional Resources](#additional-resources)
    


# Features
- **Document Management**
  - Upload educational documents to AWS S3
  - Manage document metadata with MongoDB
    - The MongoDB cluster contains 7 collections:
      - **faqs**: Contains all FAQ questions and answers.
      - **metadata**: Contains all metadata associated with the document, including:
 
          - Document ID
          - Title
          - Date and Time created
          - Moderation status
          - Document subject
          - Document grade
          - Moderation ID
          - Report status
  
      - **moderate**: Contains all data associated with the moderation process, including:
          - Moderator Information (who the moderator is)
          - Date and Time Moderated
          - Approval Status
          - Comments
          - Moderation Date and Time
      - **password_reset**: Contains data associated with password reset requests, including:
          - User ID
          - Reset Token
          - Token Expiry Date
          - Token Usage Status (a boolean value indicating if the token was used)
      - **ratings**: Contains data associated with document reporting, including:
          - Number of Times Rated
          - Total Rating
          - Average Rating
      - **reports**: Contains data associated with document rating, including:
          - Date and Time Reported
          - Reporter Information (who reported the document)
          - Reason for Report
      - **user_info**: Contains all user data, including:
          - First Name
          - Last Name
          - Email
          - Password Hash
          - User Role
          - User Qualification
    
    **Note:** The collections in the database are automatically updated based on the requests executed using Postman. Each API request interacts with specific collections, ensuring that the database reflects the most recent data corresponding to user actions.

- **User Authentication**
  - Register and authenticate Educators, Moderators and Admins
  - Token-based authentication (JWT)
- **User interactions**
  - Rate, report, moderate, and search for documents
  - Educators can upload documents
  - Download documents
- **RESTful API**
  - Built with Go and the Chi router
  - Follows RESTful principles
- **Educational Focus**
  - Supports self-directed and collaborative learning
  - Emphasizes cooperative learning and project-based teaching strategies

# Prerequisites

* Go: `Version 1.23.0` or higher
* MongoDB: Access to a MongoDB Atlas database cluster (no local installation required)
* AWS Account: With S3 access and credentials
* Git: Git `version 2.46.8` or higher installed on your machine
* Postman: For API testing (Optional but recommended)
* An IDE or Text Editor: Such as VS Code or GoLand
  
# Installation
1. **Clone the Repository**
   
   ```bash 
   git clone https://github.com/KodeKunstenaars/Share2Teach.git
   ```

2. **Navigate to the Project Directory**
   
   ```bash 
   cd Share2Teach
   ```

3. **Install Dependencies**
   ```bash 
   go mod download
   ```
# Configuration

Create a `.env` file in the root directory and add the following environment variables:

```
MONGODB_URI=your_mongodb_connection_string
AWS_REGION=your_aws_region
AWS_ACCESS_KEY_ID=your_aws_access_key_id
AWS_SECRET_ACCESS_KEY=your_aws_secret_access_key
```
**Notes:**
- Replace the placeholders with your actual credentials and settings.
  - `MONGODB_URI`: This is your remote MongoDB connection string, e.g.,
  `mongodb+srv://username:password@cluster0.mongodb.net/mydatabase?retryWrites=true&w=majority`.
  - `AWS_REGION`: Your AWS region, e.g., `us-east-1`.
  - `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`: Your AWS credentials for accessing S3.
- Keep your `.env` file **secure** and do not commit it to version control.

# Running the Application
1. **Start the Server**
   ```bash 
   go run ./cmd/api
   ```
   **Note:** Make sure you are in the root directory

   In the terminal, there will be an indication saying you are connected to mongoDB and connected to port 8080.

# API Endpoint Summaries

## Authentication Endpoints

   - **Register** (`POST /register`): Create a new user account.
   - **Login** (`POST /authenticate`): Authenticate a user and obtain JWT tokens.

## Document Management Endpoints
   - **Presign Upload** (`GET /presigned-url`): Get a presigned URL for uploading a document to AWS S3.
   - **Confirm Upload** (`POST /confirm`): Submit document metadata after uploading.
   - **Download Document** (`GET /download-document/{id}`): Retrieve a document from AWS S3.
## User Interaction Endpoints
   - **Rate Document** (`POST /rate-document/{id}`): Rate a document.
   - **Report Document** (`Route /docuements/{id}/report`) Report a document.
   - **Search Documents** (`GET /search`): Search for documents.
  
  **Note:** For detailed request and response formats, refer to the API documentation [here](./doc/)

# Usage
It is recommended to use the provided Postman collection to interact with the API. The collection includes requests for:
* **User Authentication**
  * Registering a new user
  * Logging in to obtain tokens
  * Refreshing tokens
* **Document management**
  * Uploading documents (presign, upload, confirm)
  * Downloading documents
  * Searching for documents
* **User interactions**
  * Rating documents
  * Reporting documents
  * Moderating documents (for admin and moderator users)

**User Roles and Permissions**
 
 The Share2Teach API defines four user roles:
  - **Open Access User**
    - Can search, view, and rate documents.
    - Can use the FAQ.
  - **Educator**
    - Can search, view, contribute and rate documents.
    - Can use the FAQ.
  - **Moderator**
    - Can search, view, contribute, rate and moderate documents.
    - Can use the FAQ.
  - **Admin**
    - Has unrestricted access to all system components.
  
  **Example Workflow**
  1. **Register a user**
       - Run the **Register** request with your user details.
  2. **Login**
       - Run the **Login** request to authenticate and retrieve tokens.
  3. **Upload a Document**
        - Execute the **Upload** sequence:
          - **Presign:** Get the presigned URL and `document_id`
          - **Upload to AWS:** Upload your file using the presigned URL.
          - **Confirm:** Submit the document metadata.
  4. **Search for Documents**
        - Use the **Search** request to find documents based on the parameters.
          - `key:` **title**, `parameter:` The document title, e.g., `vimcheatsheet`.
          - `key:` **subject**, `parameter:` The document subject, e.g., `IT`.
          - `key:` **grade**, `parameter:` The document grade, e.g., `cmpg323`.
  5. **Download a Document**
        - Use the **Download** request to retrieve a document.

**Note:** Ensure your API server is running and properly configured before making any requests.

# Postman Collection

We have provided a postman collection to help you test and interact with the Share2Teach API efficiently.

Download the Postman collection
- [Share2Teach.postman_collection.json](./doc/Share2Teach.postman_collection.json)

**Importing the Collection into Postman**

1. **Download** the collection file from the link above.
2. **Open Postman** on your computer.
3. **Import the collection:**
   - Click on the `Import` button in the  top-left corner.
   - Choose `Upload Files` and select the downloaded JSON file.
   - The collection and associated environment variables will be imported automatically.

**Using the Collection**

**Note:** Make sure that the backend is running on your machine before making any requests. In the terminal make sure you are in the `Root Directory` and type: `go run ./cmd/api`

- **Authentication:**
  - Start by running the **Register** request to create a new user.
  - Then, execute the **Login** request to obtain the `access_token` and the `refresh_token`. These tokens are stored within the collection variables automatically.
- **Document Operations**
  - Use the **Upload** folder to upload documents in three steps:
    1.  **Presign:** Get a presigned URL and `document_id`.
    2.  **Upload to AWS:** Use the presigned URL to upload your file.
    3.  **Confirm:** Confirm the upload by providing the metadata (found under Body-Raw)
  - Use the **Download** folder to download documents.
- **Search and Interact:**
  - Perform searches using the **Search** requests. (only `moderators` / `admins` can use the **Admin Search**)
  - Rate, report, or moderate documents using the **Interact** request.
    - Only `educators, admins,` and `moderators` can utilise the **Report** request.
    - Only `admins` and `moderators` can utilise the **Moderate** request.


# Database Migration

If you plan to migrate to a different database system (e.g., PostgreSQL, MySQL), you'll need to:

  1. **Export the Data from MongoDB**

      Use MongoDB's `mongodump` utility to export your data:


     ```bash
      mongodump --uri="your_mongodb_connection_string" --out=./mongo_backup 
      ```
    
      Here is the documentation for the `mongodump` utility: [MONGODUMP](https://www.mongodb.com/docs/database-tools/mongodump/)

  2. **Convert MongoDB Data to a Compatible Format**
   
      Since other databases use different data formats, you'll need to convert the MongoDB `BSON` data to a format like `JSON` or `CSV`.

      You can use `mongoexport` to export collections to JSON:

      ```bash
      mongoexport --uri="your_mongodb_connection_string" --collection=your_collection_name --out=your_collection_name.json
      ```

      Here is the documentation for the `mongoexport` utility: [MONGOEXPORT](https://www.mongodb.com/docs/database-tools/mongoexport/#:~:text=0%20of%20mongoexport%20.-,Synopsis,tool%20for%20backing%20up%20deployments.)

  3. **Transform the Data for the Target Database**

      Depending on the target database, you may need to transform the JSON data to match the schema and data types of the new system.

      - You can write custom scripts using Python, Node.js or Go to transform the data.

  4. **Import the Data into the new Database**

      Use the target database's import tools to load the transformed data.
        - **For PostgreSQL**:
          - You can use `psql` to import CSV files.
          - The documentation can be found here: [PSQL-DOCUMENTATION](https://www.postgresql.org/docs/current/app-psql.html)
        - **For MySQL**:
          - You can Use the `LOAD DATA INFILE` command.
          - The documentation can be found here: [MYSQL-LOAD_DATA](https://dev.mysql.com/doc/refman/8.4/en/load-data.html)
  
  5. **Update the Application Configuration**
   
      Modify the application's configuration to connect to the new database:

     - **Database Driver**: Ensure that you're using the correct databae driver.
         - Replace the MongoDB driver with the appropriate one (e.g., `lib/pq` for PostgreSQL)
     - **Connection String**: Update the connection string in your `.env` file:

        ``DATABASE_URL=your_new_database_connection_string``

  6. **Modify Data Access Code**
      - **Refactor Code**: Update the data access layer to use SQL queries instead of MongoDB queries.

  7. **Test the Application**
      - Run the test suite to ensure that all functionalities work correctly with the new database
      - Perform integration testing to check interactions with the database
  
  **Note:** Ensure to backup all data before starting the migration process.

# Troubleshooting
  - **AWS S3 Access Denied Errors**
    - **Cause**: Incorrect AWS credentials or insufficient permissions.
    - **Solution**: Verify your AWS access key ID and secret access key. Ensure the IAM user has the necessary permissions for S3 operations.
  - **Applicaton Crashes on Startup**:
    - **Cause**: Missing environment variables or misconfigured settings.
    - **Solution**: Ensure all required environment variables are set and correctly configured in your `.env` file.
  - **Invalid Tokens:** 
    - If you receive authentication errors, re-run the **Login** request to refresh your tokens.
    - Make sure you are logging in with an account that has the correct role.
  - **Server Errors:** 
    - Verify that your server is running at the URL specified in the `URL` variable in the collection.
  - **Database Connectivity Issues:** 
    - Ensure your application can connect to the remote MongoDB instance. Check your `MONGODB_URI` and network access settings.
  - **.env log error**: 
    - Ensure that the `.env` file is placed in the root directory.

Detailed API documentation is available in the `docs` directory


# Testing
Run the test suite using:

cd into the correct file directory. (`Share2Teach\internal\repository\dbrepo`)

```bash
go test ./internal/repository/dbrepo -run TestMongoDBRepo
```
This will execute the unit test.

**Note:** Ensure your test environment has access to test instances of your remote MongoDB database and AWS S3 bucket to avoid affecting production data.


# License

The product license is still pending.

# Additional Resources

- **Project Repository:** [GitHub-KodeKunstenaars/Share2Teach](https://github.com/KodeKunstenaars/Share2Teach)
- **MongoDB Atlas:** [MongoDB Atlas Documentation](https://www.mongodb.com/docs/atlas/)
- **Go Language Documentation:** [Go Official Website](https://go.dev/doc/)
- **AWS S3 Documentation:** [Amazon S3 Documentation](https://docs.aws.amazon.com/s3/)
- **Postman Documentation:** [Postman Learning Center](https://learning.postman.com/docs/introduction/overview/)

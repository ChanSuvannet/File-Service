# File Service API

A simple **File Service** built with **Go** and **Gin** for managing file uploads, retrieval, and metadata management. This service supports file uploads, fetching files by filename, and saving metadata to a MySQL database using **GORM**.

---

## Table of Contents

- [Description](#description)
- [Installation](#installation)
- [Usage](#usage)
- [File Endpoints](#file-endpoints)
- [Contributing](#contributing)
- [License](#license)

---

## Description

This **File Service API** allows users to upload files and store related metadata, such as filenames, original names, MIME types, and file paths in a database. The service is built using **Gin** for the web framework and **GORM** for interacting with a MySQL database.

### Key Features:
- Upload files (single/multiple).
- Fetch files by filename.
- Store file metadata in a MySQL database.
- Support for **file validation** and **error handling**.

---

## Installation

To get started with the File Service API, follow these steps:

### Prerequisites
- **Go** (version 1.18+)
- **MySQL** (or compatible database)
- **Git**

### Steps
1. Install dependencies
```
$ go mod tidy
```
2. Install dependencies:
```
$ go mod tidy
```
3. Run the application
```
go run main.go
```

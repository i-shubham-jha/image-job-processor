This project has been divided into several Go modules to enhance modularity and code maintainability. They implement a separation of concerns. Below is a brief description of each package. (All of them have been marked as internal for use by this project only).

# api/handlers.go
- Contains handler functions associated with API routes.
- Also includes basic data validation functions.

# cmd/retail_pulse/main.go
- Reads command line arguments (if any).
- Initializes reading of the CSV file for store IDs.
- Initializes the connection to the database.

# db
- Establishes a connection to the database.
- Follows a singleton pattern to maintain a single instance of the client in memory.

# files
- Contains data structures and functions required for downloading and saving images from URLs.

# job
- Contains the image processing function responsible for perimeter and other calculations.

# logger
- Contains functions and variables for the logger object.
- Implements an asynchronous logger using goroutines and channels.
- Follows a singleton pattern to maintain a single instance of the logger in memory.

# model
- Defines the required data models.

# store
- Contains functions to read the CSV and query for the presence of a store ID.
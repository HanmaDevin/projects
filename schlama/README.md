# Schlama Chat Application

Schlama is a modern chat application built with Go (Golang) that provides a sleek and user-friendly interface. It leverages DaisyUI for styling and htmx for dynamic updates, ensuring a responsive and visually appealing user experience.

## Features

### CLI

The CLI makes it easy to chat with local models, install new ones and add files to them.

### Web APP

The web app is just a chat interface for those who like GUIs.
It can select local models and add files.

## Prerequisites

- Go 1.20 or later

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/HanmaDevin/projects.git
   cd projects/schlama
   ```

2. Build the application:

   ```bash
   make build
   ```

3. Run the application:

   ```bash
   make run
   ```

Or make sure you have the '$GOPATH' variable set and in '$PATH' then just:

```bash
make install
```

After that you should be able to use it by just typing 'schlama' in the command line

## Usage

### Web Application

- First start the application with:

    ```bash
    ./bin/schlama chat
    ```

- Access the application in your browser at `http://localhost:8080`.
- Use the dropdown menu to select a model.
- Enter your message in the text input and click "Send".
- Upload files using the file input above the text box.

### Command-Line Tools

- **Get Help**:

    ```bash
    ./bin/schlama -h
    ```

- **List Models**:

  ```bash
  ./bin/schlama list
  ```

- **Show Local Models**:

  ```bash
  ./bin/schlama local
  ```

- **Select Model**:

  ```bash
  ./bin/schlama select <model-name>
  ```

- **Send Prompt**:

  ```bash
  ./bin/schlama prompt "Your message here"
  ```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

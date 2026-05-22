# envcrypt

A Python utility to securely encrypt and sync `.env` files across team members using age encryption.

## Installation

```bash
pip install envcrypt
```

## Usage

Encrypt your `.env` file before sharing it with your team:

```bash
# Encrypt a .env file
envcrypt encrypt .env --output .env.age

# Decrypt a received .env file
envcrypt decrypt .env.age --output .env

# Sync encrypted env with a remote store
envcrypt sync --push .env.age
envcrypt sync --pull .env.age
```

You can also use it directly in Python:

```python
from envcrypt import encrypt, decrypt

encrypt(".env", output=".env.age", recipients=["age1..."])
decrypt(".env.age", output=".env", identity="~/.age/key.txt")
```

## Configuration

Add a `envcrypt.toml` to your project root to define recipients and settings:

```toml
[envcrypt]
recipients = ["age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p"]
identity   = "~/.age/key.txt"
```

## Requirements

- Python 3.8+
- [age](https://github.com/FiloSottile/age) installed on your system

## License

This project is licensed under the [MIT License](LICENSE).
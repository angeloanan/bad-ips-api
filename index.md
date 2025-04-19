# Bad IPs API

A simple API to check if an IP address is bad or not.

Currently only supports IPv4 addresses.

## API Endpoints

### `POST /check` - Check IP

#### Request

```json
{
    "ip": "a.valid.ipv4.address"
}
```

- **`ip` (string)**: The IP address to check. It must be a valid IPv4 address.

#### Response

```json
{
    "ip": "a.valid.ipv4.address",
    "isBlocked": true
}
```

- **`ip` (string)**: The IP address that was checked.
- **`isBlocked` (boolean)**: Indicates whether the IP address is blocked or not.

## IP Source

- https://github.com/X4BNet/lists_vpn
- https://github.com/bitwire-it/ipblocklist

## Motivation

Learning Go.

That's it really.

"""JSON RPC 1.0 client."""
import http.client as http_client
import json
import logging

import requests

http_client.HTTPConnection.debuglevel = 1

logging.basicConfig()
logging.getLogger().setLevel(logging.DEBUG)
requests_log = logging.getLogger("requests.packages.urllib3")
requests_log.setLevel(logging.DEBUG)
requests_log.propagate = True


def rpc_call(url: str, method: str, args: dict):
    payload = {
        "method": method,
        "params": args,
        "id": 1,
    }
    headers = {"Content-Type": "application/json"}
    response = requests.post(url, data=json.dumps(payload), headers=headers)

    return response


def main():
    url = "http://localhost:8888/jrpc"
    result = rpc_call(url, "Arith.Multiply", [{"A": 3, "B": 5}])
    print(result.json())

    result = rpc_call(url, "Arith.Divide", [{"A": 8, "B": 2}])
    print(result.json())


if __name__ == "__main__":
    main()

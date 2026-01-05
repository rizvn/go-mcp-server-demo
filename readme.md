Based on blog post:
https://medium.com/@wadahiro/protecting-mcp-server-with-oauth-2-1-a-practical-guide-using-go-and-keycloak-7544eb5379d3


Update host file with entry:
```
127.0.0.1 host.docker.internal
```

Run Keycloak server and nginx reverse proxy to get around nginx cors with dcr
```
cd docker-compose 
docker-compose up -d
``` 


Run MCP Inspector
```
npx @modelcontextprotocol/inspector
```


In MCP Inspector connect to url:
- MCP Server URL: http://localhost:8000

When prompted to autentication use:
Username: demo
Password: 123

When prompted for consent approve the scopes.
Grant Access to MCP Inspector: Yes


# Usage
Clicking on List Tools should show available tools from MCP Server. This will list the echo tool.

Click echo tool to test it out. It will echo back the input string with currently logged in user

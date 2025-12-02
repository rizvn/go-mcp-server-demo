https://medium.com/@wadahiro/protecting-mcp-server-with-oauth-2-1-a-practical-guide-using-go-and-keycloak-7544eb5379d3

Run MCP Inspector
=================
npx @modelcontextprotocol/inspector


Keycloak Setup
=================
Keycloak Configuration

Create a new Realm: Demo

Go to Client scopes → Create client scope to add a client scope
- Name: mcp:tools
- Include in token scope: On


After creation, from the mcp:tools details screen, 
go to Mappers → Configure a new mapper → Select Audience
- Name: audience-config
- Included Custom Audience: http://localhost:8000


Go  Clients → Client registration → Client Policies → Policies tab
- Delete the default “Trusted Hosts” policy (this policy allows DCR only from trusted hosts; we delete it to allow connections from arbitrary clients like MCP Inspector)
- Open the default “Allowed Client Scopes” policy settings and set mcp:tools in "Allowed Client Scopes"


Create user in realm to login with user


# Another tutorial
https://www.simondrake.dev/blog/create-a-mcp-server-with-go
https://medium.com/@wadahiro/protecting-mcp-server-with-oauth-2-1-a-practical-guide-using-go-and-keycloak-7544eb5379d3

Run MCP Inspector

    npx @modelcontextprotocol/inspector

# Keycloak Setup

Create a new Realm: Demo

Create mcp:tools scope 

    Go to Client scopes → Create client scope to add a client scope
    - Name: mcp:tools
      - Include in token scope: On


Set custom aud value for mcp:tools scope tokens

    Client Scopes → mcp:tools → Mappers tab -> Add Mapper 
    → Configure a new mapper → Select Audience
    - Name: audience-config
      - Included Custom Audience: http://localhost:8000

    This will add aud: http://localhost:8000 to the token issued to clients with this scope.


Add email and email verified to mcp:tools scope

    Client Scopes → mcp:tools → Mappers tab -> Add Mapper > From Predefined Mapper

    Tick:
    [x] email
    [x] email verified
    
    [Save]


Allow Dynamic Client Registration (DCR) from Arbitrary Clients

    Clients → Client registration → Client Policies → Policies tab
    - Delete the default “Trusted Hosts” policy (this policy allows DCR only from trusted hosts; we delete it to allow connections from arbitrary clients like MCP Inspector)
    - Open the default “Allowed Client Scopes” policy settings and set mcp:tools in "Allowed Client Scopes"
    

Create user in realm to login with user

    Users → Add user
    - Username: demo
    - email: demo@example.com
    - email verified: On
    - Set Password: 123


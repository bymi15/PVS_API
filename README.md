# Project Virtual Showcase (PVS)

## Example API Endpoints

The route prefix is configured as `/api`, but you can change this in the `netlify.toml` config file under the `[[redirects]]` section.

| Route                        | Method | Description                                                                |
| ---------------------------- | ------ | -------------------------------------------------------------------------- |
| **/api/helloworld**          | GET    | Returns "Hello World"                                                      |
| **/api/showcase-rooms**      | GET    | Example tasks endpoint - returns all tasks                                 |
| **/api/showcase-rooms?id=1** | GET    | Example tasks endpoint - returns a task by id                              |
| **/api/showcase-rooms**      | POST   | Example tasks endpoint - creates a task (parameters parsed from body)      |
| **/api/showcase-rooms?id=1** | PUT    | Example tasks endpoint - updates a task by id(parameters parsed from body) |
| **/api/showcase-rooms?id=1** | DELETE | Example tasks endpoint - deletes a task by id                              |

## Running functions locally

```bash
go run functions/src/showcaseRooms/showcaseRooms.go -port 8000
```

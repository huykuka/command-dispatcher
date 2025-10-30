# LLM Service - Project Task Management

A Firebase Genkit-based LLM microservice for intelligent project breakdown with full task management capabilities.

## Features

✅ **AI-Powered Project Breakdown**: Break down projects into actionable tasks based on user roles (PM, Developer, Solo Dev, Architect)

✅ **Task Management**: Full CRUD operations for tasks with edit history tracking

✅ **Manual Verification**: Each task requires manual user check before project can be saved

✅ **Edit Capabilities**: Every task has an edit button to modify content

✅ **Persistent Storage**: File-based storage for projects and tasks

✅ **Validation System**: Ensures all tasks are checked before saving

## Architecture

```
backend/llm/
├── config/              # Configuration management
├── models/              # Data models and DTOs
├── prompts/             # AI prompt templates
├── services/            # Business logic layer
│   ├── genkit_llm_service.go          # LLM integration
│   ├── project_service.go             # Project breakdown
│   └── project_management_service.go  # Task management
├── storage/             # Persistence layer
├── flows/               # Genkit flow definitions
├── handlers/            # HTTP handlers
├── utils/               # Helper utilities
├── data/                # Storage directory (auto-created)
└── main.go              # Application entry point
```

## Quick Start

### 1. Prerequisites
- Go 1.21+
- Google AI API Key (for Gemini)
- Genkit CLI installed

### 2. Setup

```bash
# Navigate to service directory
cd backend/llm

# Copy environment template
cp .env.example .env

# Edit .env and add your Google API key
# GOOGLE_API_KEY=your_api_key_here

# Install dependencies
go mod tidy
```

### 3. Run the Service

```bash
# Run with environment variables
go run main_updated.go

# Or with explicit environment
GOOGLE_API_KEY=xxx go run main_updated.go
```

The service will start on:
- HTTP API: `http://localhost:8080`
- Genkit Flow UI: `http://localhost:3400`

## API Endpoints

### Health Check
```bash
GET /health
```

### LLM Operations (One-off)

#### Break Down Project
```bash
POST /api/llm/break-down-project
Content-Type: application/json

{
  "overview": "Build a user management system with authentication",
  "role": "DEVELOPER",
  "detailLevel": 3,
  "context": "Using NestJS and PostgreSQL",
  "constraints": ["2 weeks timeline", "Solo developer"]
}
```

#### Analyze Complexity
```bash
POST /api/llm/analyze-complexity
Content-Type: application/json

{
  "overview": "Build a microservices architecture with 5 services"
}
```

### Project Management (Persistent)

#### Create Project
```bash
POST /api/projects
Content-Type: application/json

{
  "name": "User Management System",
  "description": "Complete user auth and management",
  "input": {
    "overview": "Build user management with RBAC",
    "role": "SOLO_DEV",
    "detailLevel": 4
  },
  "createdBy": "user123"
}
```

**Response**: Returns project with all tasks marked as `isEditable: true` and `isChecked: false`

#### Get Project
```bash
GET /api/projects/get?id={projectId}
```

**Response**:
```json
{
  "project": {
    "id": "uuid",
    "name": "Project Name",
    "status": "draft",
    "breakdown": {...}
  },
  "managedTasks": {
    "task-1": {
      "id": "task-1",
      "title": "Setup Authentication",
      "description": "...",
      "isChecked": false,
      "isEditable": true,
      "version": 1,
      "editHistory": []
    }
  }
}
```

#### List All Projects
```bash
GET /api/projects
```

#### Update Task (Edit Button)
```bash
PUT /api/projects/tasks/update?projectId={projectId}
Content-Type: application/json

{
  "taskId": "task-1",
  "title": "Updated title",
  "description": "Updated description",
  "priority": "High",
  "estimatedHours": 8.5,
  "tags": ["backend", "security"],
  "editedBy": "user123"
}
```

**Features**:
- ✅ Tracks edit history with old/new values
- ✅ Records who edited and when
- ✅ Increments version number
- ✅ Only editable tasks can be modified

#### Check/Uncheck Task (Checkbox)
```bash
POST /api/projects/tasks/check?projectId={projectId}
Content-Type: application/json

{
  "taskId": "task-1",
  "isChecked": true,
  "checkedBy": "user123"
}
```

**Important**: User must manually check each task to verify they reviewed it.

#### Validate Project Before Save
```bash
GET /api/projects/validate?id={projectId}
```

**Response**:
```json
{
  "canSave": false,
  "allTasksChecked": false,
  "totalTasks": 10,
  "checkedTasks": 7,
  "uncheckedTaskIds": ["task-8", "task-9", "task-10"],
  "message": "3 of 10 tasks are not checked. Please review and check all tasks before saving."
}
```

#### Save Project (Only if All Checked)
```bash
POST /api/projects/save
Content-Type: application/json

{
  "projectId": "uuid",
  "userId": "user123"
}
```

**Validation**:
- ❌ Returns 400 error if any task is unchecked
- ✅ Saves and changes status to "in_progress" if all tasks checked

#### Delete Project
```bash
DELETE /api/projects/delete?id={projectId}
```

## Role-Based Task Breakdown

### PM (Project Manager)
- **Focus**: Milestones, deliverables, stakeholder management
- **Granularity**: Coarse (Epic → Feature level)
- **Output**: Gantt-compatible, dependencies, critical path

### Developer
- **Focus**: Technical implementation, code modules, testing
- **Granularity**: Medium (Feature → Story level)
- **Output**: Sprint-ready stories, acceptance criteria

### Solo Developer
- **Focus**: End-to-end ownership, MVP approach
- **Granularity**: Fine (Story → Task level)
- **Output**: Priority-sorted backlog, time-boxed tasks

### Architect
- **Focus**: System design, patterns, scalability
- **Granularity**: Architecture Decision Records
- **Output**: Component boundaries, interface definitions

## Task Management Workflow

```
1. Create Project
   ↓
2. LLM generates tasks (all isEditable=true, isChecked=false)
   ↓
3. User reviews each task
   ↓
4. User clicks "Edit" to modify task (optional)
   ↓
5. User clicks "Check" to mark task as reviewed
   ↓
6. Repeat for all tasks
   ↓
7. Validate project (check if all tasks are checked)
   ↓
8. Save project (only allowed if validation passes)
```

## Task Features

### ✅ Edit Button
- Every task can be edited
- Tracks complete edit history
- Records field changes (old value → new value)
- Stores editor and timestamp
- Increments version on each edit

### ✅ Check Button
- Manual user verification required
- Cannot save project until all tasks checked
- Visual indicator of review status
- Tracks who checked the task

### 🚫 Save Validation
```
Project can only be saved when:
- All tasks are manually checked by user
- No unchecked tasks remain
- User has reviewed the entire breakdown
```

## Storage

Projects are stored in `./data/projects/` directory as JSON files:
```
data/
└── projects/
    ├── project-uuid-1.json
    ├── project-uuid-2.json
    └── project-uuid-3.json
```

## Configuration

Environment variables in `.env`:
```env
# Required
GOOGLE_API_KEY=your_google_api_key

# Optional (with defaults)
GENKIT_FLOW_ADDR=:3400
MODEL_NAME=gemini-1.5-pro
SERVER_PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
```

## Example Usage Flow

### 1. Create a New Project
```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "E-commerce Platform",
    "description": "Build MVP for online store",
    "input": {
      "overview": "Build e-commerce with payment integration",
      "role": "SOLO_DEV",
      "detailLevel": 4,
      "constraints": ["3 weeks", "Use Stripe"]
    },
    "createdBy": "john@example.com"
  }'
```

### 2. Get Project with Tasks
```bash
curl http://localhost:8080/api/projects/get?id=abc-123
```

### 3. Edit a Task
```bash
curl -X PUT "http://localhost:8080/api/projects/tasks/update?projectId=abc-123" \
  -H "Content-Type: application/json" \
  -d '{
    "taskId": "task-1",
    "title": "Setup Stripe Payment Gateway",
    "estimatedHours": 6,
    "editedBy": "john@example.com"
  }'
```

### 4. Check Tasks
```bash
# Check task 1
curl -X POST "http://localhost:8080/api/projects/tasks/check?projectId=abc-123" \
  -H "Content-Type: application/json" \
  -d '{"taskId": "task-1", "isChecked": true, "checkedBy": "john@example.com"}'

# Check task 2
curl -X POST "http://localhost:8080/api/projects/tasks/check?projectId=abc-123" \
  -H "Content-Type: application/json" \
  -d '{"taskId": "task-2", "isChecked": true, "checkedBy": "john@example.com"}'

# ... repeat for all tasks
```

### 5. Validate Before Save
```bash
curl http://localhost:8080/api/projects/validate?id=abc-123
```

### 6. Save Project (if all checked)
```bash
curl -X POST http://localhost:8080/api/projects/save \
  -H "Content-Type: application/json" \
  -d '{
    "projectId": "abc-123",
    "userId": "john@example.com"
  }'
```

## Development

### Run Tests
```bash
go test ./...
```

### Run with Genkit Dev UI
```bash
genkit start
```

### Build for Production
```bash
go build -o llm-service main_updated.go
```

## Integration with SAF-Tools

This service integrates with the main SAF-Tools backend:
- Use JWT tokens from main backend for authentication
- Sync user IDs with SAF backend user system
- Store project references in SAF backend database
- Use consistent DTOs for data exchange

## Troubleshooting

### "Project cannot be saved" Error
- Check validation endpoint to see which tasks are unchecked
- Ensure all tasks have `isChecked: true`
- User must manually verify each task

### LLM Generation Fails
- Verify GOOGLE_API_KEY is set correctly
- Check API quota and rate limits
- Review prompt in `prompts/project_prompts.go`

### Storage Issues
- Ensure `./data/projects/` directory has write permissions
- Check disk space
- Review logs for file I/O errors

## License

Part of the SAF-Tools project.


# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`rtool` is a CLI tool written in Go that generates files for the Repeat Store (rpst) PHP application. It automates the creation of DDL files and PHP Request Model classes following specific naming conventions and templates.

## Build and Run

```bash
# Build the binary
go build

# Run directly with go
go run main.go <command>

# After building, run the binary
./rtool <command>
```

## Commands

### Create DDL Files

Creates SQL DDL files following rpst's naming convention: `YYYYMMDD_SORTNUM_TICKET_WHEN_EXPLANATION.sql`

```bash
./rtool create ddl \
  --ticket_number <redmine_ticket_number> \
  --when <execution_timing> \
  --explanation <table_operation_description> \
  [--sort_number <execution_order>] \
  [--file_path <output_directory>]
```

Required flags:
- `-t, --ticket_number`: Redmine ticket number
- `-w, --when`: When to execute the SQL (deployment timing)
- `-e, --explanation`: Brief description using lowercase alphanumeric and underscores (e.g., `create_users_table`, `update_items_status`)

Optional flags:
- `-s, --sort_number`: Execution order when multiple SQL files exist (default: "0000")
- `-f, --file_path`: Output directory path (default: current directory)

### Create Request Models

Generates PHP Request Model classes from templates, optionally with Factory classes.

```bash
./rtool create requestModel \
  --name <ModelName> \
  --namespace <namespace_path> \
  [--mode <request_mode>] \
  [--with-factory] \
  [--stdout]
```

Required flags:
- `-n, --name`: Request Model class name
- `-a, --namespace`: Namespace under `App\model\requests\rpst\` (e.g., `admin\system`)

Optional flags:
- `-m, --mode`: Request mode constant value (default: "default")
- `-f, --with-factory`: Also generate a Factory class
- `--stdout`: Output to stdout instead of creating files

## Architecture

### Command Structure

Built with [Cobra](https://github.com/spf13/cobra), the CLI follows this structure:

```
rtool
└── create
    ├── ddl        (creates DDL files)
    └── requestModel (creates PHP Request Models)
```

- `cmd/root.go`: Root command definition and initialization
- `cmd/create.go`: All `create` subcommands and their logic
- `main.go`: Entry point that calls `cmd.Execute()`

### Template System

Templates are located in `cmd/templates/`:
- `requestModel.tmpl`: PHP Request Model template with validation structure
- `requestModelFactory.tmpl`: Factory pattern implementation for Request Models

Templates use Go's `text/template` package with the following data structure:
- `{{.Name}}`: Class name
- `{{.Namespace}}`: PHP namespace path
- `{{.Mode}}`: Request mode constant

### Generated PHP Files

Request Models implement:
- `LegacyRpstWebRequestModelInterface`
- Traits: `AdminOnetimeCsrfTokenTrait`, `LegacyRpstWebRequestModelTrait`
- Mode constant for request routing
- Validation rules via `getValidationRuleForData()`

Request Model Factories:
- Implement `LegacyRpstWebRequestModelFactoryInterface`
- Use `LegacyRpstWebRequestModelFactoryTrait`
- Map modes to Request Model classes via `$CLASS_MAP`

## Code Patterns

### Adding New Commands

1. Define a new `cobra.Command` in `cmd/create.go` (or create a new file)
2. Implement the `RunE` function with error handling
3. Register the command in the `init()` function using `AddCommand`
4. Define flags with appropriate types and validation
5. Use `MarkFlagRequired` for mandatory flags

### Adding New Templates

1. Create `.tmpl` file in `cmd/templates/`
2. Use `{{.FieldName}}` syntax for template variables
3. Load with `template.New("filename.tmpl").ParseFiles("./cmd/templates/filename.tmpl")`
4. Execute with data map: `template.Execute(writer, data)`

### File Creation Pattern

The codebase uses a consistent pattern for file generation:
1. Parse command-line flags into option structs
2. Generate filename based on conventions and input
3. Create parent directories if needed with `os.MkdirAll`
4. Choose output destination (file or stdout)
5. Execute template or write content
6. Provide user feedback with file path

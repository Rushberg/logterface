# Logterface

A lightweight utility to analyze logs in real time.

## Features
- Processes piped input logs from another application.
- Supports various visualization types like graphs, progress bars, and numeric statistics.
- Configurable via JSON to adapt to different log formats and visualization needs.

## How It Works
1. The application reads a JSON configuration file that defines handlers and layouts for processing logs.
2. Logs are fed via piped input, and each line is parsed and processed based on the configured handlers.
3. Outputs are displayed in an organized format as defined by the layouts in the configuration.

## Getting Started

### Prerequisites
- Go SDK 1.23.6 or above.
- A log source that can pipe its output into this tool.

### Installation
1. Download windows or linux release.
2. Alternatively Clone the repository.

   2.1. Build the application:
   ```bash
   go build -o logterface
   ```

### Usage
Run the application by specifying the configuration file as an argument and piping input logs:
```bash
your-log-generator | ./logterface path/to/config.json
```
You can run example configuration to see that everything working as expected:
```bash
go run ./example/logsGen.go | go run ./main.go example/config.json
```
*(Make sure you have enough space in your terminal for the layout)*

### Configuration

#### Handlers Format
Handlers are responsible for extracting and processing specific log patterns. Each handler processes data matched by its `regex`, capturing relevant information for analysis. Below is the structure of a handler:

- **type**: The type of handler, which describes its function. Examples include:
    - **Graph**: Displays data as a graph.
    - **Numbers**: Performs arithmetic operations on numeric data (e.g., max, min, average).
    - **Progress**: Tracks progress based on the log data.
    - **Filter**: Captures and filters specific log patterns for display.

- **id**: A unique identifier for the handler, used to reference the handler in layouts.

- **regex**: A regular expression used to match specific patterns from the logs.  
  **The first capturing group in the regex defines what data will be extracted from the log for processing.** For example:
    - Regex: `.*number: (\\d+)` → Captures the number following the phrase "number: ".
    - Regex: `progress: (\\d+)` → Captures the number after "progress: ".

- **params**: Additional parameters specific to the handler type:
    - `name`: A descriptive name for the handler.
    - `method`: Describes how data is processed (e.g., `max`, `min`, `avg`, `latest`).
    - `thresholdMethod` and `threshold`: Define threshold logic for certain handlers.
    - Dimensions (e.g., `width`, `height`) used for graphical representation.

### Example Handlers
```json
{
    "type": "Numbers",
    "id": "max",
    "regex": ".*number: (\\d+)",
    "params": {
        "name": "Max Number",
        "method": "max"
    }
},
{
    "type": "Graph",
    "id": "graph",
    "regex": ".*number: (\\d+)",
    "params": {
        "name": "Numbers Graph",
        "width": 50,
        "height": 15
    }
}
```
- **Explanation**:
    - The first handler (**Numbers**) extracts numbers captured by the regex `.*number: (\\d+)` and calculates the maximum value.
    - The second handler (**Graph**) uses the same regex to generate a graphical representation of the captured numbers.

## Layouts Format
Layouts define how the processed data from handlers is displayed. Each layout organizes the output of one or more handlers in a visually meaningful way.

- **type**: Specifies the layout type. Options include:
    - **Pipe**: Prints the original logs.
    - **Line**: Displays multiple handlers side by side in a single row and refreshed every <refresh_mills>.

- **handlers**: A list of handler `id`s whose outputs will be included in the layout.

### Example Layouts
```json
{
    "type": "Line",
    "handlers": [
        {
            "id": "graph"
        },
        {
            "id": "max"
        },
        {
            "id": "min"
        },
        {
            "id": "avg"
        }
    ]
}
```
- **Explanation**:
    - This layout displays the outputs of the handlers with IDs `graph`, `max`, `min`, and `avg` in a single line.

## Capturing Data with Regex
The **regex** field in each handler plays a vital role in extracting the relevant data from log lines. The data captured by the first group in a regex (i.e., the part inside `(...)`) is what the handler processes. For example:
- Log Line: `number: 64`
    - Regex: `.*number: (\\d+)` → Captures `64` for processing.
- Log Line: `progress: 42`
    - Regex: `progress: (\\d+)` → Captures `42` for processing.

This captured data is then used by the handler, depending on its type and parameters.

## Combining Handlers and Layouts
Handlers process the log data, while layouts define how the processed data is displayed. Together, they provide flexibility for organizing, visualizing, and analyzing log information in real time.

### Sample Configuration
```json
{
  "refresh_mills": 1000,
  "handlers": [
    {
      "type": "Numbers",
      "id": "max",
      "regex": ".*number: (\\d+)",
      "params": {
        "name": "Max Number",
        "method": "max"
      }
    },
    {
      "type": "Graph",
      "id": "graph",
      "regex": ".*number: (\\d+)",
      "params": {
        "name": "Numbers Graph",
        "width": 50,
        "height": 15
      }
    }
  ],
  "layouts": [
    {
      "type": "Line",
      "handlers": [
        {
          "id": "graph"
        },
        {
          "id": "max"
        }
      ]
    }
  ]
}
```
In this sample:
- The `max` handler calculates the highest number from log entries matching `.*number: (\\d+)`.
- The `graph` handler creates a graphical visualization of the numbers captured by the same regex.
- The `Line` layout places both `graph` and `max` outputs in a single row.
### Sample Input Logs
The application works with structured logs. For example:
```
number: 64
progress: 42 
resources usage: 50/100 
error: Uh oh! Something bad happened. Code:123
```

## License
MIT License.

---

Developed for simplified and powerful real-time log analysis.
# Golang Kubernetes Operator

This project is a simple Kubernetes operator written in Go that watches for changes in ConfigMaps. It specifically looks for ConfigMaps that contain the annotation `secret-injected: true` and updates their data accordingly.

## Project Structure

```
golang-k8s-operator
├── cmd
│   └── manager
│       └── main.go          # Entry point of the operator
├── config
│   ├── crd
│   │   └── bases            # Directory for Custom Resource Definitions (CRDs)
│   ├── default              # Default configuration files
│   ├── manager              # Configuration files for the manager
│   ├── rbac                 # Role-Based Access Control (RBAC) configuration files
│   └── samples              # Sample ConfigMaps and resources
├── controllers
│   └── configmap_controller.go # Implementation of the ConfigMap controller
├── api
│   └── v1
│       └── groupversion_info.go # API group and version information
├── Dockerfile                # Instructions for building the Docker image
├── go.mod                    # Module dependencies
├── go.sum                    # Checksums for module dependencies
├── Makefile                  # Commands for building and managing the operator
└── README.md                 # Project documentation
```

## Getting Started

To get started with this operator, follow these steps:

1. **Clone the repository**:
   ```
   git clone <repository-url>
   cd golang-k8s-operator
   ```

2. **Build the operator**:
   ```
   make build
   ```

3. **Run the operator**:
   ```
   make run
   ```

## Usage

This operator will watch for ConfigMaps in the Kubernetes cluster. Ensure that your ConfigMaps have the annotation `secret-injected: true` for the operator to take action.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.